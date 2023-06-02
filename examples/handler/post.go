package handler

import (
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

func GetPost(ctx *fast.Ctx) error {

	id := ctx.Param("id").AsInt(ctx)

	cook := &http.Cookie{
		Name:     "test",
		Value:    "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(ctx.Writer(), cook)

	// ctx.Cookie()
	return ctx.Status(200).JSON(map[string]interface{}{
		"id": id,
		// "cookie": cook,
	})
}
