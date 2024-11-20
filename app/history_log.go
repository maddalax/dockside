package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"log/slog"
	"paas/app/subject"
	"paas/app/util/must"
	"time"
)

func LogChange(locator *service.Locator, subject subject.Subject, data map[string]any) {
	client := service.Get[KvClient](locator)
	err := client.CreateHistoryStream()
	if err != nil {
		slog.Error("failed to create history stream: %v", err)
		return
	}
	data["created_at"] = time.Now().Format(time.Stamp)
	data["subject"] = subject
	err = client.Publish(subject, must.Serialize(data))
	if err != nil {
		slog.Error("failed to publish history: %v", err)
	}
}