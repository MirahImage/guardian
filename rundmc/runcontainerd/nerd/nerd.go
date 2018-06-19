package nerd

import (
	"context"
	"strconv"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/linux/runctypes"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type Nerd struct {
	client  *containerd.Client
	context context.Context
}

func New(client *containerd.Client, context context.Context) *Nerd {
	return &Nerd{
		client:  client,
		context: context,
	}
}

func WithNoNewKeyring(ctx context.Context, c *containerd.Client, ti *containerd.TaskInfo) error {
	ti.Options = &runctypes.CreateOptions{NoNewKeyring: true}
	return nil
}

func (n *Nerd) Create(log lager.Logger, containerID string, spec *specs.Spec) error {
	log.Debug("creating-container", lager.Data{"containerID": containerID})
	container, err := n.client.NewContainer(n.context, containerID, containerd.WithSpec(spec))
	if err != nil {
		return err
	}

	log.Debug("creating-task", lager.Data{"containerID": containerID})
	task, err := container.NewTask(n.context, cio.NullIO, WithNoNewKeyring)
	if err != nil {
		return err
	}

	log.Debug("starting-task", lager.Data{"containerID": containerID})
	return task.Start(n.context)
}

func (n *Nerd) Delete(log lager.Logger, containerID string) error {
	container, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return err
	}

	log.Debug("deleting-task", lager.Data{"containerID": containerID})
	_, err = task.Delete(n.context, containerd.WithProcessKill)
	if err != nil {
		return err
	}

	log.Debug("deleting-container", lager.Data{"containerID": containerID})
	return container.Delete(n.context)
}

func (n *Nerd) State(log lager.Logger, containerID string) (int, containerd.ProcessStatus, error) {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return 0, "", err
	}

	log.Debug("getting-task-status", lager.Data{"containerID": containerID})
	status, err := task.Status(n.context)
	if err != nil {
		return 0, "", err
	}

	log.Debug("task-result", lager.Data{"containerID": containerID, "pid": strconv.Itoa(int(task.Pid())), "status": string(status.Status)})
	return int(task.Pid()), status.Status, nil
}

func (n *Nerd) Exec(log lager.Logger, containerID, processID string, spec *specs.Process, io garden.ProcessIO) (*Process, error) {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return nil, err
	}

	log.Debug("execing-task", lager.Data{"containerID": containerID, "processID": processID})
	process, err := task.Exec(n.context, processID, spec, cio.NewCreator(cio.WithStdio, withGardenProcessIO(io)))
	if err != nil {
		return nil, err
	}

	if err := process.Start(n.context); err != nil {
		return nil, err
	}

	return &Process{context: n.context, containerdProcess: process}, nil
}

func (n *Nerd) GetContainerPID(log lager.Logger, containerID string) (uint32, error) {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return 0, err
	}

	return task.Pid(), nil
}

func withGardenProcessIO(io garden.ProcessIO) cio.Opt {
	return func(opt *cio.Streams) {
		if io.Stdin != nil {
			opt.Stdin = io.Stdin
		}

		if io.Stdout != nil {
			opt.Stdout = io.Stdout
		}

		if io.Stderr != nil {
			opt.Stderr = io.Stderr
		}
	}
}

func (n *Nerd) loadContainerAndTask(log lager.Logger, containerID string) (containerd.Container, containerd.Task, error) {
	log.Debug("loading-container", lager.Data{"containerID": containerID})
	container, err := n.client.LoadContainer(n.context, containerID)
	if err != nil {
		return nil, nil, err
	}

	log.Debug("loading-task", lager.Data{"containerID": containerID})
	task, err := container.Task(n.context, nil)
	if err != nil {
		return nil, nil, err
	}

	return container, task, nil
}

type Process struct {
	context           context.Context
	containerdProcess containerd.Process
}

func (p *Process) Wait() (int, error) {
	exitCh, err := p.containerdProcess.Wait(p.context)
	if err != nil {
		return 0, err
	}

	exitStatus := <-exitCh
	return int(exitStatus.ExitCode()), nil
}
