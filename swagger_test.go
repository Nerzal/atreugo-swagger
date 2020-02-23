package atreugoswagger

import (
	"testing"

	"github.com/savsgio/atreugo/v10"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
	"github.com/valyala/fasthttp"
)

func TestWrapHandler(t *testing.T) {

	config := &atreugo.Config{
		Addr: "0.0.0.0:1337",
	}

	a := atreugo.New(config)

	a.GET("/*", WrapHandler)

	go a.ListenAndServe()

	t.Run("Get Index.html", func(t *testing.T) {
		status, _, err := fasthttp.Get(nil, "http://localhost:1337/index.html")
		if err != nil {
			t.Error("failed to call hello: ", err)
		}

		if status != 200 {
			t.Error("received wrong statuscode")
		}
	})

	// w1 := performRequest("GET", "/index.html", a)
	// assert.Equal(t, 200, w1.Code)

	// w2 := performRequest("GET", "/doc.json", a)
	// assert.Equal(t, 200, w2.Code)

	// w3 := performRequest("GET", "/favicon-16x16.png", a)
	// assert.Equal(t, 200, w3.Code)

	// w4 := performRequest("GET", "/notfound", a)
	// assert.Equal(t, 404, w4.Code)

}
