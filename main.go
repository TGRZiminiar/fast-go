// //package main

// // import (
// // 	"github.com/TGRZiminiar/based/examples/handler"
// // 	"github.com/TGRZiminiar/based/launch"
// // )

// // func main() {
// // 	launch.Post("/post", handler.CreatePost)
// // 	launch.Get("/post/:id/app/:mangaid", handler.GetPost)
// // 	launch.Start()
// // }

package main

import (
	"github.com/TGRZiminiar/based/examples/handler"
	"github.com/TGRZiminiar/based/fast"
)

func main() {

	app := fast.Init(fast.Config{
		Port:   8080,
		Logger: true,
	})

	app.Get("/post/:id/:mangaid", handler.GetPost)
	app.Post("/post", handler.CreatePost)

	app.Start()
	// r := gin.Default()

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"api": "golang",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (or "PORT" env var)
	// launch.Post("/post", handler.CreatePost)
	// launch.Get("/post/:id", handler.GetPost)

	// launch.Start()

}
