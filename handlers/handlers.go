package handlers

import (
	"Message-Queue/messaging"
	"net/http"

	"github.com/labstack/echo"
)

type GenericResponse struct {
	Error   string
	Data    interface{}
	Message string
}

type Handler struct {
	Broker *messaging.Broker
}

func NewHandler(broker *messaging.Broker) *Handler {
	return &Handler{
		Broker: broker,
	}
}

func (h *Handler) Subscribe(c echo.Context) error {
	var req messaging.SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid request Format",
		})
	}
	subID, err := h.Broker.Subscribe(req.Endpoint, req.Topics)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &GenericResponse{
		Data: map[string]string{
			"subscriber_id": subID,
		},
		Message: "Subscribed successfully",
	})
}

func (h *Handler) Unsubscribe(c echo.Context) error {
	var req messaging.UnsubscribeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid request Format",
		})
	}

	err := h.Broker.Unsubscribe(req.SubscriberID, req.Topics)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &GenericResponse{
		Message: "Unsubscribed successfully to the Specified Topics",
	})
}

func (h *Handler) Publish(c echo.Context) error {
	var req messaging.PublishRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid Request Format",
		})
	}
	err := h.Broker.Publish(req.Topic, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &GenericResponse{
		Message: "Message published successfully to all Topic subscribers",
	})
}

func (h *Handler) Retrieve(c echo.Context) error {
	var req messaging.RetrieveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid Request Format",
		})
	}

	msg, err := h.Broker.Retrieve(req.SubscriberId, req.Topic, req.StartDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &GenericResponse{
		Data: msg,
	})
}
