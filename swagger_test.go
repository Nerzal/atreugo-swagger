package atreugoswagger

import (
	"strings"
	"testing"

	"github.com/savsgio/atreugo/v11"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
	"github.com/valyala/fasthttp"
)

func TestWrapHandler(t *testing.T) {
	config := atreugo.Config{
		Addr: "0.0.0.0:1337",
	}

	a := atreugo.New(config)

	a.GET("/docs/{doc:*}", AtreugoWrapHandler())
	a.GET("/api/v1/hello", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("hello", fasthttp.StatusOK)
	})

	go a.ListenAndServe()

	t.Run("Get Index.html", func(t *testing.T) {
		status, result, err := fasthttp.Get(nil, "http://localhost:1337/docs/index.html")
		if err != nil {
			t.Error("failed to call hello: ", err)
		}

		if status != fasthttp.StatusOK {
			t.Error("received wrong statuscode")
		}

		stringResult := string(result)
		if !strings.Contains(stringResult, "<title>Swagger UI</title>") {
			t.Error("Could not find swagger title")
		}
	})

	t.Run("Get doc.json", func(t *testing.T) {
		status, result, err := fasthttp.Get(nil, "http://localhost:1337/docs/doc.json")
		if err != nil {
			t.Error("failed to call hello: ", err)
		}

		if status != fasthttp.StatusOK {
			t.Error("received wrong statuscode")
		}

		stringResult := string(result)
		if !strings.Contains(stringResult, `"title": "Swagger Example API",`) {
			t.Error("Could not find swagger title")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		status, _, err := fasthttp.Get(nil, "http://localhost:1337/notfound")
		if err != nil {
			t.Error("failed to call hello: ", err)
		}

		if status != fasthttp.StatusNotFound {
			t.Error("received wrong statuscode")
		}
	})

	t.Run("Get Hello", func(t *testing.T) {
		status, response, err := fasthttp.Get(nil, "http://localhost:1337/api/v1/hello")
		if err != nil {
			t.Error("failed to call hello: ", err)
		}

		if status != fasthttp.StatusOK {
			t.Error("received wrong statuscode")
		}

		if string(response) != "hello" {
			t.Errorf("received wrong response: %s", string(response))
		}
	})
}
