package main

import (
	"github.com/TGRZiminiar/based/examples/handler"
	"github.com/TGRZiminiar/based/fast"
)

func main() {
	app := fast.Init(fast.Config{
		Port:   5000,
		Logger: true,
	})

	app.UpdateCorsConfig(fast.CorsConfig{
		AllowOrigins:     "*",
		AllowCredentials: true,
		MaxAge:           3600,
		AllowCookies:     true,
	})

	app.Get("/post/:id/:mangaid", handler.UserAuth, handler.GetPost)
	app.Post("/posts", handler.CreatePost)

	app.Start()
}
