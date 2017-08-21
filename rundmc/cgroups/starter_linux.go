package cgroups

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	specs "github.com/opencontainers/runtime-spec/specs-go"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/guardian/logging"
	"code.cloudfoundry.org/lager"
)

func deviceWildcard() *int64 {
	var i int64 = -1
	return &i
}

var (
	worldReadWrite = os.FileMode(0666)
	FuseDevice     = specs.LinuxDevice{
		Path:     "/dev/fuse",
		Type:     "c",
		Major:    10,
		Minor:    229,
		FileMode: &worldReadWrite,
	}

	allowedDevices = []*specs.LinuxDeviceCgroup{{Access: "rwm", Type: "c", Major: intRef(FuseDevice.Major), Minor: intRef(FuseDevice.Minor), Allow: true}}

	ociDefaultAllowedDevices = []*specs.LinuxDeviceCgroup{
		{Access: "m", Type: "c", Major: deviceWildcard(), Minor: deviceWildcard(), Allow: true},
		{Access: "m", Type: "b", Major: deviceWildcard(), Minor: deviceWildcard(), Allow: true},
		{Access: "rwm", Type: "c", Major: intRef(1), Minor: intRef(3), Allow: true},          // /dev/null
		{Access: "rwm", Type: "c", Major: intRef(1), Minor: intRef(8), Allow: true},          // /dev/random
		{Access: "rwm", Type: "c", Major: intRef(1), Minor: intRef(7), Allow: true},          // /dev/full
		{Access: "rwm", Type: "c", Major: intRef(5), Minor: intRef(0), Allow: true},          // /dev/tty
		{Access: "rwm", Type: "c", Major: intRef(1), Minor: intRef(5), Allow: true},          // /dev/zero
		{Access: "rwm", Type: "c", Major: intRef(1), Minor: intRef(9), Allow: true},          // /dev/urandom
		{Access: "rwm", Type: "c", Major: intRef(5), Minor: intRef(1), Allow: true},          // /dev/console
		{Access: "rwm", Type: "c", Major: intRef(136), Minor: deviceWildcard(), Allow: true}, // /dev/pts/*
		{Access: "rwm", Type: "c", Major: intRef(5), Minor: intRef(2), Allow: true},          // /dev/ptmx
		{Access: "rwm", Type: "c", Major: intRef(10), Minor: intRef(200), Allow: true},       // /dev/net/tun
	}
)

type Starter struct {
	*CgroupStarter
}

const cgroupsHeader = "#subsys_name hierarchy num_cgroups enabled"

type CgroupsFormatError struct {
	Content string
}

func intRef(i int64) *int64 {
	return &i
}

func (err CgroupsFormatError) Error() string {
	return fmt.Sprintf("unknown /proc/cgroups format: %s", err.Content)
}

func NewStarter(logger lager.Logger, procCgroupReader io.ReadCloser, procSelfCgroupReader io.ReadCloser, cgroupMountpoint, gardenCgroup string, runner commandrunner.CommandRunner, chowner Chowner) *Starter {
	return &Starter{
		&CgroupStarter{
			CgroupPath:      cgroupMountpoint,
			GardenCgroup:    gardenCgroup,
			ProcCgroups:     procCgroupReader,
			ProcSelfCgroups: procSelfCgroupReader,
			CommandRunner:   runner,
			Logger:          logger,
			Chowner:         chowner,
			//TODO: add teh device list here
		},
	}
}

type CgroupStarter struct {
	CgroupPath    string
	GardenCgroup  string
	CommandRunner commandrunner.CommandRunner

	ProcCgroups     io.ReadCloser
	ProcSelfCgroups io.ReadCloser

	Logger  lager.Logger
	Chowner Chowner
}

func (s *CgroupStarter) Start() error {
	return s.mountCgroupsIfNeeded(s.Logger)
}

func (s *CgroupStarter) mountCgroupsIfNeeded(logger lager.Logger) error {
	defer s.ProcCgroups.Close()
	defer s.ProcSelfCgroups.Close()
	if err := os.MkdirAll(s.CgroupPath, 0755); err != nil {
		return err
	}

	if !s.isMountPoint(s.CgroupPath) {
		s.mountTmpfsOnCgroupPath(logger, s.CgroupPath)
	} else {
		logger.Info("cgroups-tmpfs-already-mounted", lager.Data{"path": s.CgroupPath})
	}

	subsystemGroupings, err := s.subsystemGroupings()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(s.ProcCgroups)

	if !scanner.Scan() {
		return CgroupsFormatError{Content: "(empty)"}
	}

	if _, err := fmt.Sscanf(scanner.Text(), cgroupsHeader); err != nil {
		return CgroupsFormatError{Content: scanner.Text()}
	}

	for scanner.Scan() {
		var subsystem string
		var skip, enabled int
		n, err := fmt.Sscanf(scanner.Text(), "%s %d %d %d ", &subsystem, &skip, &skip, &enabled)
		if err != nil || n != 4 {
			return CgroupsFormatError{Content: scanner.Text()}
		}

		if enabled == 0 {
			continue
		}

		subsystemToMount, dirToCreate := subsystem, s.GardenCgroup
		if v, ok := subsystemGroupings[subsystem]; ok {
			subsystemToMount = v.SubSystem
			dirToCreate = path.Join(v.Path, s.GardenCgroup)
		}

		subsystemMountPath := path.Join(s.CgroupPath, subsystem)
		// TODO rename
		if err := s.idempotentCgroupMount(logger, subsystemMountPath, subsystemToMount); err != nil {
			return err
		}

		gardenCgroupPath := filepath.Join(s.CgroupPath, subsystem, dirToCreate)
		if err := s.createGardenCgroup(logger, gardenCgroupPath); err != nil {
			return err
		}

		if subsystem == "devices" {
			if err := s.modifyAllowedDevices(gardenCgroupPath, append(allowedDevices, ociDefaultAllowedDevices...)); err != nil {
				return err
			}
		}

		if err := s.Chowner.RecursiveChown(gardenCgroupPath); err != nil {
			return err
		}
	}

	return nil
}

func (s *CgroupStarter) modifyAllowedDevices(dir string, devices []*specs.LinuxDeviceCgroup) error {
	if has, err := hasSubdirectories(dir); err != nil {
		return err
	} else if has {
		return nil
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "devices.deny"), []byte("a"), 0770); err != nil {
		return err
	}
	for _, device := range devices {
		data := fmt.Sprintf("%s %s:%s %s", device.Type, s.deviceNumberString(device.Major), s.deviceNumberString(device.Minor), device.Access)
		if err := s.setDeviceCgroup(dir, "devices.allow", data); err != nil {
			return err
		}
	}

	return nil
}

func hasSubdirectories(dir string) (bool, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, fileInfo := range dirs {
		if fileInfo.Mode().IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func (d *CgroupStarter) setDeviceCgroup(dir, file, data string) error {
	if err := ioutil.WriteFile(filepath.Join(dir, file), []byte(data), 0); err != nil {
		return fmt.Errorf("failed to write %s to %s: %v", data, file, err)
	}

	return nil
}

func (s *CgroupStarter) deviceNumberString(number *int64) string {
	if *number == -1 {
		return "*"
	}
	return fmt.Sprint(*number)
}

func (s *CgroupStarter) createGardenCgroup(log lager.Logger, gardenCgroupPath string) error {
	log = log.Session("creating-garden-cgroup", lager.Data{"gardenCgroup": gardenCgroupPath})
	log.Info("started")
	defer log.Info("finished")

	if err := os.MkdirAll(gardenCgroupPath, 0700); err != nil {
		return err
	}

	return nil
}

func (s *CgroupStarter) mountTmpfsOnCgroupPath(log lager.Logger, path string) {
	log = log.Session("cgroups-tmpfs-mounting", lager.Data{"path": path})
	log.Info("started")

	if err := s.CommandRunner.Run(exec.Command("mount", "-t", "tmpfs", "-o", "uid=0,gid=0,mode=0755", "cgroup", path)); err != nil {
		log.Error("mount-failed-continuing-anyway", err)
	} else {
		log.Info("finished")
	}
}

type group struct {
	SubSystem string
	Path      string
}

func (s *CgroupStarter) subsystemGroupings() (map[string]group, error) {
	groupings := map[string]group{}

	scanner := bufio.NewScanner(s.ProcSelfCgroups)
	for scanner.Scan() {
		segs := strings.Split(scanner.Text(), ":")
		if len(segs) != 3 {
			continue
		}

		subsystems := strings.Split(segs[1], ",")
		for _, subsystem := range subsystems {
			groupings[subsystem] = group{segs[1], segs[2]}
		}
	}

	return groupings, scanner.Err()
}

func (s *CgroupStarter) idempotentCgroupMount(logger lager.Logger, cgroupPath, subsystems string) error {
	logger = logger.Session("mount-cgroup", lager.Data{
		"path":       cgroupPath,
		"subsystems": subsystems,
	})

	logger.Info("started")

	if !s.isMountPoint(cgroupPath) {
		if err := os.MkdirAll(cgroupPath, 0755); err != nil {
			return fmt.Errorf("mkdir '%s': %s", cgroupPath, err)
		}

		cmd := exec.Command("mount", "-n", "-t", "cgroup", "-o", subsystems, "cgroup", cgroupPath)
		cmd.Stderr = logging.Writer(logger.Session("mount-cgroup-cmd"))
		if err := s.CommandRunner.Run(cmd); err != nil {
			return fmt.Errorf("mounting subsystems '%s' in '%s': %s", subsystems, cgroupPath, err)
		}
	} else {
		logger.Info("subsystems-already-mounted")
	}

	logger.Info("finished")

	return nil
}

func (s *CgroupStarter) isMountPoint(path string) bool {
	// append trailing slash to force symlink traversal; symlinking e.g. 'cpu'
	// to 'cpu,cpuacct' is common
	return s.CommandRunner.Run(exec.Command("mountpoint", "-q", path+"/")) == nil
}