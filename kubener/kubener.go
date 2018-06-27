package kubener

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
)

type Kubener struct {
	Logger     lager.Logger
	containers []KubeContainer
}

func (k *Kubener) Ping() error { return nil }

func (k *Kubener) Capacity() (garden.Capacity, error) {
	return garden.Capacity{}, errors.New("capacity not implemented")
}

func writeToFile(file, msg string) {
	d1 := []byte(fmt.Sprintf("%s\n", msg))
	ioutil.WriteFile(file, d1, 0644)
}

func (k *Kubener) Create(containerSpec garden.ContainerSpec) (garden.Container, error) {
	logger := k.Logger.Session("start")

	image := "busybox"
	if containerSpec.Image.URI != "" {
		image = containerSpec.Image.URI
	}

	kubeCommand := exec.Command("kubectl", "--kubeconfig=/root/.kube/config", "create", "-f", "/tmp/manifest.json")
	manifest := fmt.Sprintf(
		`{"apiVersion": "v1", "kind": "Pod", "metadata": {"name": "%s"}, "spec":{"containers":[{"name":"%s", "image":"%s", "args": ["sleep", "100000"]}]}}`,
		containerSpec.Handle,
		containerSpec.Handle,
		image,
	)
	writeToFile("/tmp/manifest.json", manifest)

	out, err := kubeCommand.CombinedOutput()
	logger.Info("kubener-create", lager.Data{"output": string(out), "manifest": manifest})
	writeToFile("/tmp/dat1", fmt.Sprintf("output: %s\nmanifest: %s\n", string(out), manifest))
	if err != nil {
		return nil, err
	}

	return NewKubeContainer(containerSpec.Handle), nil
}

func (k *Kubener) Destroy(handle string) error { return errors.New("destroy not implemented") }

func (k *Kubener) Containers(garden.Properties) ([]garden.Container, error) {
	return []garden.Container{}, nil
}

func (k *Kubener) BulkInfo(handles []string) (map[string]garden.ContainerInfoEntry, error) {
	return nil, errors.New("bulkinfo not implemented")
}

func (k *Kubener) BulkMetrics(handles []string) (map[string]garden.ContainerMetricsEntry, error) {
	return nil, errors.New("bulkmetrics not implemented")
}

func (k *Kubener) Lookup(handle string) (garden.Container, error) {
	return nil, errors.New("lookup not implemented")
}

func (k *Kubener) Start() error { return nil }
func (k *Kubener) Stop()        {}

func (k *Kubener) GraceTime(garden.Container) time.Duration {
	return time.Hour * 24
}
