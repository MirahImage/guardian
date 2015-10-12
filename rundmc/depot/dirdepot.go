package depot

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pivotal-golang/lager"
)

var ErrDoesNotExist = errors.New("does not exist")

//go:generate counterfeiter . BundleSaver
type BundleSaver interface {
	Save(path string) error
}

// a depot which stores containers as subdirs of a depot directory
type DirectoryDepot struct {
	dir string
}

func New(dir string) *DirectoryDepot {
	return &DirectoryDepot{
		dir: dir,
	}
}

func (d *DirectoryDepot) Create(log lager.Logger, handle string, bundle BundleSaver) error {
	log = log.Session("create", lager.Data{"handle": handle})

	log.Info("started")
	defer log.Info("finished")

	path := d.toDir(handle)
	if err := os.MkdirAll(path, 0700); err != nil {
		log.Error("mkdir", err, lager.Data{"path": path})
		return err
	}

	if err := bundle.Save(path); err != nil {
		removeOrLog(log, path)
		log.Error("create", err, lager.Data{"path": path})
		return err
	}

	return nil
}

func (d *DirectoryDepot) Lookup(log lager.Logger, handle string) (string, error) {
	log = log.Session("lookup", lager.Data{"handle": handle})

	log.Info("started")
	defer log.Info("finished")

	if _, err := os.Stat(d.toDir(handle)); err != nil {
		return "", ErrDoesNotExist
	}

	return d.toDir(handle), nil
}

func (d *DirectoryDepot) Destroy(log lager.Logger, handle string) error {
	log = log.Session("destroy", lager.Data{"handle": handle})

	log.Info("started")
	defer log.Info("finished")

	return os.RemoveAll(d.toDir(handle))
}

func (d *DirectoryDepot) toDir(handle string) string {
	return filepath.Join(d.dir, handle)
}

func removeOrLog(log lager.Logger, path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Error("remove-failed", err, lager.Data{"path": path})
	}
}
