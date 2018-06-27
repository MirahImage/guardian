package kubener

import (
	"errors"
	"io"
	"time"

	"code.cloudfoundry.org/garden"
)

type KubeContainer struct {
	handle string
}

type KubeProcess struct{}

func NewKubeContainer(handle string) *KubeContainer {
	return &KubeContainer{handle: handle}
}

func (p *KubeProcess) ID() string                  { return "id-not-implemented" }
func (p *KubeProcess) Wait() (int, error)          { return 0, errors.New("wait not implemented") }
func (p *KubeProcess) SetTTY(garden.TTYSpec) error { return errors.New("settty not implemented") }
func (p *KubeProcess) Signal(garden.Signal) error  { return errors.New("signal not implemented") }

func (c *KubeContainer) Handle() string       { return c.handle }
func (c *KubeContainer) Stop(kill bool) error { return errors.New("stop not implemented") }
func (c *KubeContainer) Info() (garden.ContainerInfo, error) {
	return garden.ContainerInfo{}, errors.New("info not implemented")
}
func (c *KubeContainer) StreamIn(spec garden.StreamInSpec) error {
	return errors.New("streamin not implemented")
}
func (c *KubeContainer) StreamOut(spec garden.StreamOutSpec) (io.ReadCloser, error) {
	return nil, errors.New("streamout not implemented")
}
func (c *KubeContainer) CurrentBandwidthLimits() (garden.BandwidthLimits, error) {
	return garden.BandwidthLimits{}, errors.New("currentbandwidthlimits not implemented")
}
func (c *KubeContainer) CurrentCPULimits() (garden.CPULimits, error) {
	return garden.CPULimits{}, errors.New("currentspulimits not implemented")
}
func (c *KubeContainer) CurrentDiskLimits() (garden.DiskLimits, error) {
	return garden.DiskLimits{}, errors.New("currentdisklimits not implemented")
}
func (c *KubeContainer) CurrentMemoryLimits() (garden.MemoryLimits, error) {
	return garden.MemoryLimits{}, errors.New("currentmemorylimits not implemented")
}
func (c *KubeContainer) NetIn(hostPort, containerPort uint32) (uint32, uint32, error) {
	return 0, 0, errors.New("netin not implemented")
}
func (c *KubeContainer) NetOut(netOutRule garden.NetOutRule) error {
	return errors.New("netout not implemented")
}

func (c *KubeContainer) BulkNetOut(netOutRules []garden.NetOutRule) error {
	return errors.New("bulknetout not implemented")
}

func (c *KubeContainer) Run(garden.ProcessSpec, garden.ProcessIO) (garden.Process, error) {
	return &KubeProcess{}, errors.New("run not implemented")
}

func (c *KubeContainer) Attach(processID string, io garden.ProcessIO) (garden.Process, error) {
	return &KubeProcess{}, errors.New("attach not implemented")
}
func (c *KubeContainer) Metrics() (garden.Metrics, error) {
	return garden.Metrics{}, errors.New("metrics not implemented")
}
func (c *KubeContainer) SetGraceTime(graceTime time.Duration) error {
	return errors.New("setgracetime not implemented")
}
func (c *KubeContainer) Properties() (garden.Properties, error) {
	return garden.Properties{}, errors.New("properties not implemented")
}
func (c *KubeContainer) Property(name string) (string, error) {
	return "", errors.New("property not implemented")
}
func (c *KubeContainer) SetProperty(name string, value string) error {
	return errors.New("setproperty not implemented")
}
func (c *KubeContainer) RemoveProperty(name string) error {
	return errors.New("removeproperty not implemented")
}
