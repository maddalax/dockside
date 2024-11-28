package app

import (
	"dockside/app/logger"
	"dockside/app/util/syncutil"
)

func (a *Agent) monitorInstanceCount() {
	resources, err := GetResourcesForServer(a.locator, a.serverId)
	if err != nil {
		logger.ErrorWithFields("Failed to get resources for server", err, map[string]any{
			"server_id": a.serverId,
		})
		return
	}
	wg := syncutil.NewWaitGroupWithConcurrency(3)
	for _, resource := range resources {
		wg.Add()
		go func() {
			defer wg.Done()
			switch resource.BuildMeta.(type) {
			case *DockerBuildMeta, *DockerRegistryMeta:
				a.monitorDockerInstanceCount(resource)
			}
		}()
	}
	wg.Wait()
}

func (a *Agent) monitorDockerInstanceCount(resource *Resource) {
	client, err := DockerConnect(a.locator)
	if err != nil {
		logger.Error("Failed to connect to docker", err)
		return
	}
	containers, err := client.GetRunningContainers(resource)
	if err != nil {
		logger.Error("Failed to get running containers", err)
		return
	}
	// matches, all good
	if len(containers) == resource.InstancesPerServer {
		logger.DebugWithFields("instance count running matches expected", map[string]any{
			"resource_id": resource.Id,
			"count":       len(containers),
		})
		return
	}

	logger.InfoWithFields("instance count running does not match expected, attempting to fix", map[string]any{
		"resource_id": resource.Id,
		"count":       len(containers),
		"expected":    resource.InstancesPerServer,
	})

	// client.Run will automatically scale to the correct number of instances
	err = client.Run(resource, RunOptions{
		IgnoreIfRunning: true,
	})

	if err != nil {
		logger.Error("Failed to run resource", err)
	}
}