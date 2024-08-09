package event

import (
	"github.com/RobinHoodArmyHQ/robin-api/models"
)

type GetEventFeedRequest struct {
	Offset int
	Limit  int
	CityId int32
}

type GetEventFeedResponse struct {
	Events []*models.Event
	Count  int64
}
