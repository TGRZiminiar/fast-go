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
	"time"

	"github.com/TGRZiminiar/based/examples/handler"
	"github.com/TGRZiminiar/based/fast"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	app := fast.Init(fast.Config{
		Port:   8080,
		Logger: true,
	})

	app.UpdateCorsConfig(fast.CorsConfig{
		AllowOrigins:     "not",
		AllowCredentials: true,
		MaxAge:           3600,
		AllowCookies:     true,
	})

	app.Get("/post/:id/:mangaid", handler.Test, handler.GetPost)
	app.Post("/post", handler.CreatePost)

	app.Start()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/post/:id/:postid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"api": "golang",
		})
		// c.Cookie()
	})
	r.Run()
	// launch.Post("/post", handler.CreatePost)
	// launch.Get("/post/:id", handler.GetPost)

	// launch.Start()
	// http.HandleFunc("/data", getData)
	// http.ListenAndServe(":8080", nil)

}

// func getData(w http.ResponseWriter, r *http.Request) {
// 	// Create your data
// 	data := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 	}

// 	// Marshal the data into JSON manually
// 	jsonData := "{"
// 	for key, value := range data {
// 		jsonData += fmt.Sprintf(`"%s":"%s",`, key, value)
// 	}
// 	jsonData = jsonData[:len(jsonData)-1] + "}"

// 	// Set the Content-Type header to application/json
// 	w.Header().Set("Content-Type", "application/json")

// 	// Write the JSON response to the response writer
// 	w.Write([]byte(jsonData))
// }
