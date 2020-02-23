package atreugoswagger

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/savsgio/atreugo/v10"
	"github.com/swaggo/swag"
	"github.com/valyala/fasthttp"
)

// Config stores atreugoswagger configuration variables.
type Config struct {
	//The url pointing to API definition (normally swagger.json or swagger.yaml). Default is `doc.json`.
	URL string
}

// URL presents the url pointing to API definition (normally swagger.json or swagger.yaml).
func URL(url string) func(c *Config) {
	return func(c *Config) {
		c.URL = url
	}
}

// WrapHandler serves swagger files
var WrapHandler = AtreugoWrapHandler()

// EchoWrapHandler wraps `http.Handler` into `atreugo.Middleware`.
func AtreugoWrapHandler(confs ...func(c *Config)) func(ctx *atreugo.RequestCtx) error {
	config := &Config{
		URL: "doc.json",
	}

	for _, c := range confs {
		c(config)
	}

	// create a template with name
	t := template.New("swagger_index.html")
	index, _ := t.Parse(indexTempl)

	type pro struct {
		Host string
	}

	var re = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[\?|.]*`)

	return func(ctx *atreugo.RequestCtx) error {
		var matches []string

		if matches = re.FindStringSubmatch(string(ctx.RequestURI())); len(matches) != 3 {
			return ctx.TextResponse("404 page not found", fasthttp.StatusNotFound)
		}

		path := matches[2]

		if strings.Contains(path, "png") {
			ctx.Response.Header.Add("Content-Type", "image/png")
		}

		if strings.Contains(path, "jpg") {
			ctx.Response.Header.Add("Content-Type", "image/jpeg")
		}

		if strings.Contains(path, "html") {
			ctx.Response.Header.Add("Content-Type", "text/html")
		}

		if strings.Contains(path, "css") {
			ctx.Response.Header.Add("Content-Type", "text/css")
		}

		if strings.Contains(path, "js") {
			ctx.Response.Header.Add("Content-Type", "text/javascript")
		}

		switch path {
		case "index.html":
			return index.Execute(ctx.Response.BodyWriter(), config)
		case "doc.json":
			doc, _ := swag.ReadDoc()
			return ctx.TextResponse(doc, fasthttp.StatusOK)
		case "favicon-16x16.png":
			return ctx.RawResponseBytes(FileFavicon16x16Png, fasthttp.StatusOK)
		case "favicon-32x32.png":
			return ctx.RawResponseBytes(FileFavicon32x32Png, fasthttp.StatusOK)
		case "oauth2-redirect.html":
			return ctx.RawResponseBytes(FileOauth2RedirectHTML, fasthttp.StatusOK)
		case "swagger-ui-bundle.js":
			return ctx.RawResponseBytes(FileSwaggerUIBundleJs, fasthttp.StatusOK)
		case "swagger-ui-bundle.js.map":
			return ctx.RawResponseBytes(FileSwaggerUIBundleJsMap, fasthttp.StatusOK)
		case "swagger-ui-standalone-preset.js":
			return ctx.RawResponseBytes(FileSwaggerUIStandalonePresetJs, fasthttp.StatusOK)
		case "swagger-ui-standalone-preset.js.map":
			return ctx.RawResponseBytes(FileSwaggerUIStandalonePresetJsMap, fasthttp.StatusOK)
		case "swagger-ui.css":
			return ctx.RawResponseBytes(FileSwaggerUICSS, fasthttp.StatusOK)
		case "swagger-ui.css.map":
			return ctx.RawResponseBytes(FileSwaggerUICSSMap, fasthttp.StatusOK)
		case "swagger-ui.js":
			return ctx.RawResponseBytes(FileSwaggerUIJs, fasthttp.StatusOK)
		case "swagger-ui.js.map":
			return ctx.RawResponseBytes(FileSwaggerUIJsMap, fasthttp.StatusOK)
		}

		return ctx.TextResponse("404 page not found", fasthttp.StatusNotFound)
	}
}
