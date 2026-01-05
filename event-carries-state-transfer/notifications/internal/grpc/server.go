package grpc

import (
	"eda-in-golang/notifications/internal/application"
	"eda-in-golang/notifications/notificationspb"
)

type server struct {
	app application.App
	notificationspb.UnimplementedNotificationsServiceServer
}
