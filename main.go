package main // import "hello-wv"

import (
	"embed"
	"net/http"

	"github.com/inkeliz/gowebview"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ncruces/zenity"
)

var (
	//go:embed static
	content embed.FS
)

const defaultName = ``

func selectFileSave(c echo.Context) error {
	zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename(defaultName),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})

	return c.String(http.StatusOK, "ok")
}

func setupServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(
		middleware.CORS(),
		middleware.Recover(),
	)

	contentHandler := echo.WrapHandler(http.FileServer(http.FS(content)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	e.GET("/save-sel", selectFileSave)
	e.GET("/*", contentHandler, contentRewrite)

	return e
}

func main() {
	e := setupServer()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE,
			echo.HEAD, echo.OPTIONS,
		},
	}))

	go e.Start("127.0.0.1:2918")

	w, err := gowebview.New(&gowebview.Config{
		URL: "http://localhost:2918/index.html",
		WindowConfig: &gowebview.WindowConfig{
			Title: "안녕 세상!!",
			// Size: &gowebview.Point{
			// 	X: 640, Y: 720,
			// },
		},
	})
	if err != nil {
		panic(err)
	}

	defer w.Destroy()
	w.Run()
}
