// main.go
package main

import (
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"io"
	"log"
	"net/http"
)

func main() {
	app := pocketbase.New()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// or you can also use the shorter e.Router.GET("/articles/:slug", handler, middlewares...)
		_, err := e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/custom/chatamount",
			Handler: func(c echo.Context) error {
				record, err := app.Dao().FindFirstRecordByData("chats", "message", "Hello!")
				if err != nil {
					return apis.NewNotFoundError("The chat does not exist.", err)
				}

				// enable ?expand query param support
				err = apis.EnrichRecord(c, app.Dao(), record)
				if err != nil {
					return err
				}

				return c.JSON(http.StatusOK, record.Get("id"))
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
	resp, err := http.Get("http://194.195.243.228/api/custom/chatamount")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	//read data from b
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(b))
}
