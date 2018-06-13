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
	socket  string
	client  *containerd.Client
	context context.Context
}

func New(containerdSocket string, context context.Context) *Nerd {
	return &Nerd{
		socket:  containerdSocket,
		context: context,
	}
}

func (n *Nerd) getClient() (*containerd.Client, error) {
	if n.client != nil {
		return n.client, nil
	}
	var err error
	n.client, err = containerd.New(n.socket)
	return n.client, err
}

func WithNoNewKeyring(ctx context.Context, c *containerd.Client, ti *containerd.TaskInfo) error {
	ti.Options = &runctypes.CreateOptions{NoNewKeyring: true}
	return nil
}

func (n *Nerd) Create(log lager.Logger, containerID string, spec *specs.Spec) error {
	log.Debug("creating-container", lager.Data{"containerID": containerID})
	client, err := n.getClient()
	if err != nil {
		return err
	}
	container, err := client.NewContainer(n.context, containerID, containerd.WithSpec(spec))
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
	log.Debug("loading-container", lager.Data{"containerID": containerID})
	client, err := n.getClient()
	if err != nil {
		return err
	}
	container, err := client.LoadContainer(n.context, containerID)
	if err != nil {
		return err
	}

	log.Debug("loading-task", lager.Data{"containerID": containerID})
	task, err := container.Task(n.context, nil)
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
	log.Debug("loading-container", lager.Data{"containerID": containerID})
	client, err := n.getClient()
	if err != nil {
		return 0, "", err
	}
	container, err := client.LoadContainer(n.context, containerID)
	if err != nil {
		return 0, "", err
	}

	log.Debug("loading-task", lager.Data{"containerID": containerID})
	task, err := container.Task(n.context, nil)
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

func (n *Nerd) Exec(log lager.Logger, containerID, processID string, spec *specs.Process, io garden.ProcessIO) error {
	log.Debug("loading-container", lager.Data{"containerID": containerID})
	client, err := n.getClient()
	if err != nil {
		return err
	}
	container, err := client.LoadContainer(n.context, containerID)
	if err != nil {
		return err
	}

	log.Debug("loading-task", lager.Data{"containerID": containerID})
	task, err := container.Task(n.context, nil)
	if err != nil {
		return err
	}

	log.Debug("execing-task", lager.Data{"containerID": containerID, "processID": processID})
	process, err := task.Exec(n.context, processID, spec, cio.NewCreator(cio.WithStdio, withGardenProcessIO(io)))
	if err != nil {
		return err
	}

	return process.Start(n.context)
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
