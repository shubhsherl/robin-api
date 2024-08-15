package checkin

import (
	"github.com/RobinHoodArmyHQ/robin-api/internal/env"
	"github.com/RobinHoodArmyHQ/robin-api/internal/handler/contract"
	"github.com/RobinHoodArmyHQ/robin-api/internal/repositories/checkin"
	"github.com/RobinHoodArmyHQ/robin-api/internal/repositories/event"
	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateCheckInHandler(c *gin.Context) {
	req := &checkin.CreateCheckInRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &contract.Response{Message: err.Error()})
		return
	}

	// check if event exists
	eventRepo := env.FromContext(c).EventRepository
	eventResp, err := eventRepo.GetEvent(&event.GetEventRequest{EventID: req.CheckIn.EventID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &contract.Response{Message: err.Error()})
		return
	}

	if eventResp == nil || eventResp.Event == nil {
		c.JSON(http.StatusNotFound, &contract.Response{Message: "event not found"})
		return
	}

	// check if event has started
	if !eventResp.Event.Started() {
		c.JSON(http.StatusBadRequest, &contract.Response{Message: "event not started"})
		return
	}

	checkInRepo := env.FromContext(c).CheckInRepository
	checkInResp, err := checkInRepo.CreateCheckIn(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &contract.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkInResp)
}

func GetCheckInHandler(c *gin.Context) {
	checkInID := c.Param("checkin_id")
	if checkInID == "" {
		c.JSON(http.StatusBadRequest, &contract.Response{Message: "checkin_id is required"})
		return
	}

	checkInRepo := env.FromContext(c).CheckInRepository
	checkInResp, err := checkInRepo.GetCheckIn(&checkin.GetCheckInRequest{CheckInID: nanoid.NanoID(checkInID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &contract.Response{Message: err.Error()})
		return
	}

	if checkInResp == nil || checkInResp.CheckIn == nil {
		c.JSON(http.StatusNotFound, &contract.Response{Message: "check_in not found"})
		return
	}

	c.JSON(http.StatusOK, checkInResp)
}

func GetUserCheckInsHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &contract.Response{Message: "user_id is required"})
		return
	}

	checkInRepo := env.FromContext(c).CheckInRepository
	checkInResp, err := checkInRepo.GetUserCheckIns(&checkin.GetUserCheckInsRequest{UserID: nanoid.NanoID(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &contract.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkInResp)
}
