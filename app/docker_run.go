package app

import (
	"context"
	"dockman/app/logger"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
	"io"
	"strconv"
	"strings"
)

type RunOptions struct {
	Stdout io.WriteCloser
	// If we should kill the existing container that's running first
	RemoveExisting bool
	// Whether we should just return if the container is already running
	IgnoreIfRunning bool
}

func (c *DockerClient) GetContainer(resource *Resource, index int) (types.ContainerJSON, error) {
	containerName := fmt.Sprintf("%s-%s-container-%d", resource.Name, resource.Id, index)
	return c.cli.ContainerInspect(context.Background(), containerName)
}

func (c *DockerClient) Stop(resource *Resource) error {
	for i := range resource.InstancesPerServer {
		containerName := fmt.Sprintf("%s-%s-container-%d", resource.Name, resource.Id, i)
		err := c.cli.ContainerStop(context.Background(), containerName, container.StopOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DockerClient) Run(resource *Resource, opts RunOptions) error {
	instances := resource.InstancesPerServer
	if instances == 0 {
		instances = 1
	}

	c.ReduceToMatchResourceCount(resource, instances)

	for i := range instances {
		err := c.doRun(resource, i, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DockerClient) doRun(resource *Resource, index int, opts RunOptions) error {
	ctx := context.Background()
	imageName := fmt.Sprintf("%s-%s", resource.Name, resource.Id)
	containerName := fmt.Sprintf("%s-%s-container-%d", resource.Name, resource.Id, index)

	err := c.LoadImage(imageName)

	if err != nil {
		return err
	}

	if opts.IgnoreIfRunning {
		exists, err := c.GetContainer(resource, index)
		// if the container exists and is running, we can just return
		if err == nil && exists.State.Running {
			return nil
		}
	}

	err = c.cli.ContainerStop(ctx, containerName, container.StopOptions{})

	if err != nil {
		switch err.(type) {
		case errdefs.ErrNotFound:
			// don't need to worry about it if the container doesn't exist
			err = nil
		default:
			return err
		}
	}

	if opts.RemoveExisting {
		err = c.cli.ContainerRemove(ctx, containerName, container.RemoveOptions{
			Force: true,
		})
	}

	if err != nil {
		switch err.(type) {
		case errdefs.ErrNotFound:
			// don't need to worry about it if the container doesn't exist
			err = nil
		default:
			return err
		}
	}

	hostPort, err := FindOpenPort()

	if err != nil {
		return err
	}

	exposedPort := 0

	switch b := resource.BuildMeta.(type) {
	case *DockerBuildMeta:
		exposedPort = b.ExposedPort
	}

	if exposedPort == 0 {
		return ResourceExposedPortNotSetError
	}

	exposedPortFmt := nat.Port(fmt.Sprintf("%d/tcp", exposedPort))

	// Define port bindings
	portBindings := nat.PortMap{
		exposedPortFmt: []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(hostPort),
			},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyUnlessStopped,
		},
		LogConfig: container.LogConfig{
			Type: "fluentd",
			Config: map[string]string{
				"fluentd-address": "localhost:24224",
				"fluentd-async":   "true",
				"labels":          "dockman.resource.id,dockman.build.id",
			},
		},
	}

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		ExposedPorts: map[nat.Port]struct{}{
			exposedPortFmt: {},
		},
		AttachStdout: true,
		AttachStderr: true,
	}, hostConfig, nil, nil, containerName)

	if err != nil {
		switch err.(type) {
		case errdefs.ErrNotFound:
			return ResourceNotFoundError
		case errdefs.ErrConflict:
			// container already exists, it failed to get killed for some reason
			if opts.RemoveExisting {
				return ContainerExistsError
			} else {
				// we don't want to remove existing, so lets run the current one
				err = nil
			}
		default:
			return err
		}
	}

	err = c.cli.ContainerStart(ctx, containerName, container.StartOptions{})

	if err != nil {
		// another container may have taken the port, lets try a different one
		if strings.Contains(err.Error(), "port is already allocated") {
			logger.ErrorWithFields("Port is already allocated, trying a different one", err, map[string]any{
				"container_name": containerName,
			})
			for i := 0; i < 50; i++ {
				err = c.doRun(resource, index, opts)
				if err == nil {
					return nil
				}
			}
		}
		// the port this container is trying to bind to is already in use
		// this can happen if we reboot the container and something else took it
		// kind of edge case, but it can happen, ideally we should be able to kill the container
		// and start it again, but we can't do that if opts.RemoveExisting is false
		if strings.Contains(err.Error(), "address already in use") {
			return ResourcePortInUseError(strconv.Itoa(hostPort))
		}

		return err
	}

	if opts.Stdout != nil {
		return c.StreamLogs(resp.ID, ctx, StreamLogsOptions{
			Stdout: opts.Stdout,
		})
	}

	return nil
}
