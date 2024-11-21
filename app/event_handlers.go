package app

import (
	"dockside/app/logger"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
)

type EventHandler struct {
	locator *service.Locator
}

func NewEventHandler(locator *service.Locator) *EventHandler {
	return &EventHandler{
		locator: locator,
	}
}

func (eh *EventHandler) OnJobStarted(job *Job) {
	registry := GetServiceRegistry(eh.locator)
	logger.DebugWithFields("job started", map[string]any{
		"job_name": job.name,
	})
	registry.GetJobMetricsManager().OnJobStarted(job)
}

func (eh *EventHandler) OnJobFinished(job *Job) {
	registry := GetServiceRegistry(eh.locator)
	logger.DebugWithFields("job finished", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
		"duration":   fmt.Sprintf("%dms", job.lastRunDuration.Milliseconds()),
	})
	registry.GetJobMetricsManager().OnJobFinished(job)
}

func (eh *EventHandler) OnJobStopped(job *Job) {
	registry := GetServiceRegistry(eh.locator)
	logger.DebugWithFields("job stopped", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
	})
	registry.GetJobMetricsManager().OnJobStopped(job)
}

func (eh *EventHandler) OnServerDisconnected(server *Server) {
	logger.InfoWithFields("server disconnected", map[string]any{
		"server_id": server.Id,
		"name":      server.FormattedName(),
	})
}

func (eh *EventHandler) OnServerConnected(server *Server) {
	logger.InfoWithFields("server connected", map[string]any{
		"server_id": server.Id,
		"name":      server.FormattedName(),
	})
}

func (eh *EventHandler) OnServerDetached(serverId string, resource *Resource) {
	logger.InfoWithFields("server detached from resource", map[string]any{
		"server_id":     serverId,
		"resource_id":   resource.Id,
		"resource_name": resource.Name,
	})
}

func (eh *EventHandler) OnResourceStatusChange(resource *Resource, status RunStatus) {
	logger.InfoWithFields("resource status changed", map[string]any{
		"resource_id":   resource.Id,
		"resource_name": resource.Name,
		"new_status":    status,
	})
}