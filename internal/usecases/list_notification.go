package usecases

import (
	"context"
	"fmt"

	"github.com/diegorezm/reminder/internal/store"
)

type NOTIFICATION_TYPE string

const (
	NOTIFICATION_TYPE_ALL             NOTIFICATION_TYPE = "all"
	NOTIFICATION_TYPE_ALL_IGNORE_DATE NOTIFICATION_TYPE = "all-ignore-date"
	NOTIFICATION_TYPE_ABOUT_TO_EXPIRE NOTIFICATION_TYPE = "about-to-expire"
	NOTIFICATION_TYPE_DISMISSED       NOTIFICATION_TYPE = "dismissed"
	NOTIFICATION_TYPE_EXPIRED         NOTIFICATION_TYPE = "expired"
)

type ListNotificationInput struct {
	ReminderID int64
	Type       NOTIFICATION_TYPE
}

func ListNotification(input ListNotificationInput, s *store.Queries, ctx context.Context) ([]store.Notification, error) {
	var notifications []store.Notification
	var err error

	if input.Type == NOTIFICATION_TYPE_ALL {
		notifications, err = s.GetNotifications(ctx, input.ReminderID)
	} else if input.Type == NOTIFICATION_TYPE_ABOUT_TO_EXPIRE {
		notifications, err = s.GetAboutToExpireNotifications(ctx, input.ReminderID)
	} else if input.Type == NOTIFICATION_TYPE_DISMISSED {
		notifications, err = s.GetDismissedNotifications(ctx)
	} else if input.Type == NOTIFICATION_TYPE_ALL_IGNORE_DATE {
		notifications, err = s.GetAllNotifications(ctx, input.ReminderID)
	} else if input.Type == NOTIFICATION_TYPE_EXPIRED {
		notifications, err = s.GetExpiredNotifications(ctx)
	} else {
		return nil, fmt.Errorf("unknown notification type: %s", input.Type)
	}

	return notifications, err
}
