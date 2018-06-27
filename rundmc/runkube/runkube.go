package runkube

import (
	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gardener"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	"code.cloudfoundry.org/lager"
)

type RunKube struct{}

func New() *RunKube {
	return &RunKube{}
}

func (r *RunKube) Create(log lager.Logger, bundlePath, id string, io garden.ProcessIO) error {
	return errors.New("not implemented")
}

func (r *RunKube) Exec(log lager.Logger, bundlePath, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
	return nil, errors.New("not implemented")
}

func (r *RunKube) Attach(log lager.Logger, bundlePath, id, processId string, io garden.ProcessIO) (garden.Process, error) {
	return nil, errors.New("not implemented")
}

func (r *RunKube) Kill(log lager.Logger, bundlePath string) error {
	return errors.New("not implemented")
}

func (r *RunKube) Delete(log lager.Logger, force bool, id string) error {
	return errors.New("not implemented")
}

func (r *RunKube) State(log lager.Logger, id string) (runrunc.State, error) {
	return runrunc.State{}, errors.New("not implemented")
}

func (r *RunKube) Stats(log lager.Logger, id string) (gardener.ActualContainerMetrics, error) {
	return gardener.ActualContainerMetrics{}, errors.New("not implemented")
}

func (r *RunKube) WatchEvents(log lager.Logger, id string, eventsNotifier runrunc.EventsNotifier) error {
	return errors.New("not implemented")
}
