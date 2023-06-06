package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TGRZiminiar/based/fast"
)

func CreatePost(ctx *fast.Ctx) error {
	file, header, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]string{
			"msg": "No file found in the form data",
		})
	}
	defer file.Close()

	return ctx.JSON(map[string]interface{}{
		"msg": header,
	})
	// return nil
}

// func CreatePost(ctx *launch.PostCtx[any]) error {
// 	// params := ctx.RequestParams()
// 	fmt.Println("FUCKING SHIT")
// 	return nil
// 	// val := ctx.FormValue("name")

// 	// fmt.Println(val)
// 	// return ctx.JSON(map[string]string{
// 	// 	"mes": val,
// 	// })
// }

type Post struct {
	id int
}

func GetPost(ctx *fast.Ctx) error {

	fmt.Println(ctx.R.Temp)

	id := ctx.Param("id").AsInt(ctx)

	ctx.SetCookie(http.Cookie{
		Name:     "test",
		Value:    "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	return ctx.Status(200).JSON(map[string]interface{}{
		"id":     id,
		"cookie": ctx.W.Header().Get("Content-Type"),
	})
}

func UserAuth(ctx *fast.Ctx) error {

	ctx.R.Temp = "user"
	return nil
}
