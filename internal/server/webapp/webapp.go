package webapp

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/paletas/paletas_website/internal/server/webapp/ws"
)

type WebApp struct {
	internalServer *fiber.App
}

func NewWebApp() *WebApp {
	fiberServer := fiber.New()
	configureFiberServer(fiberServer)

	return &WebApp{
		internalServer: fiberServer,
	}
}

func configureFiberServer(app *fiber.App) {
	// Serve static files
	app.Static("/", "./website")

	configureRoutes(app)
	ws.ConfigureChat(app)
}

func configureRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

func (app WebApp) Start(listenPort string) error {
	err := app.internalServer.Listen(listenPort)
	if err != nil {
		return err
	}
	return nil
}
