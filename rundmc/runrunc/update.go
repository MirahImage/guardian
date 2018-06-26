package runrunc

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
)

type Updater struct {
	runner RuncCmdRunner
	runc   RuncBinary
}

func NewUpdater(runner RuncCmdRunner, runc RuncBinary) *Updater {
	return &Updater{
		runner: runner,
		runc:   runc,
	}
}

func (u *Updater) UpdateLimits(log lager.Logger, handle string, limits garden.Limits) error {
	log = log.Session("update", lager.Data{"handle": handle})

	log.Info("started")
	defer log.Info("finished")

	limitsBytes, err := json.Marshal(limits)
	if err != nil {
		return err
	}

	return u.runner.RunAndLog(log, func(logFile string) *exec.Cmd {
		cmd := u.runc.UpdateCommand(handle, logFile)
		cmd.Stdin = bytes.NewReader(limitsBytes)
		return cmd
	})
}
