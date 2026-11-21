package main

import (
	"Message-Queue/handlers"
	"Message-Queue/messaging"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	var broker = messaging.NewBroker()
	var handler = handlers.NewHandler(broker)

	e.POST("/v1/subscribe", handler.Subscribe)
	e.POST("/v1/unsunscribe", handler.Unsubscribe)
	e.POST("/v1/publish", handler.Publish)
	e.GET("/v1/retrieve", handler.Retrieve)
	e.GET("/v1/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// START THE SERVER
	fmt.Println("Server is running at port 3000")
	log.Fatal(e.Start(":3000"))
}
