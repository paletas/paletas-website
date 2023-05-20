package webapp

import (
	"github.com/gofiber/fiber/v2"
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
	app.Static("/", "./website", fiber.Static{
		Compress: true,
		Index:    "index.html",
	})

	configureRoutes(app)
}

func configureRoutes(app *fiber.App) {
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
}

func (app WebApp) Start(listenPort string) error {
	err := app.internalServer.Listen(listenPort)
	if err != nil {
		return err
	}
	return nil
}
