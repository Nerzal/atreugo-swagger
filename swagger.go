package atreugoswagger

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/Nerzal/atreugo-swagger/assets"
	"github.com/savsgio/atreugo/v11"
	"github.com/swaggo/swag"
	"github.com/valyala/fasthttp"
)

const fileRegex = `(.*)(redoc\.html|index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[\?|.]*`

// Config stores atreugoswagger configuration variables.
type Config struct {
	// The url pointing to API definition (normally swagger.json or swagger.yaml). Default is `doc.json`.
	URL string
	// Title is used for the redoc documentation
	Title string
}

// Title presents the title of the tab
func Title(title string) func(c *Config) {
	return func(c *Config) {
		c.Title = title
	}
}

// URL presents the url pointing to API definition (normally swagger.json or swagger.yaml).
func URL(url string) func(c *Config) {
	return func(c *Config) {
		c.URL = url
	}
}

// AtreugoWrapHandler is a handler which serves swagger files
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

	var re = regexp.MustCompile(fileRegex)

	return func(ctx *atreugo.RequestCtx) error {
		var matches []string

		if matches = re.FindStringSubmatch(string(ctx.RequestURI())); len(matches) != 3 {
			return ctx.TextResponse("404 page not found", fasthttp.StatusNotFound)
		}

		path := matches[2]
		setContentType(ctx, path)

		switch path {
		case "index.html":
			return index.Execute(ctx.Response.BodyWriter(), config)
		case "redoc.html":
			redocHTML := strings.Replace(assets.RedocDocumentation, "{title}", config.Title, 1)
			return ctx.HTTPResponse(redocHTML, fasthttp.StatusOK)
		case "doc.json":
			doc, _ := swag.ReadDoc()
			return ctx.TextResponse(doc, fasthttp.StatusOK)
		case "favicon-16x16.png":
			return ctx.RawResponseBytes(assets.FileFavicon16x16Png, fasthttp.StatusOK)
		case "favicon-32x32.png":
			return ctx.RawResponseBytes(assets.FileFavicon32x32Png, fasthttp.StatusOK)
		case "oauth2-redirect.html":
			return ctx.RawResponseBytes(assets.FileOauth2RedirectHTML, fasthttp.StatusOK)
		case "swagger-ui-bundle.js":
			return ctx.RawResponseBytes(assets.FileSwaggerUIBundleJs, fasthttp.StatusOK)
		case "swagger-ui-bundle.js.map":
			return ctx.RawResponseBytes(assets.FileSwaggerUIBundleJsMap, fasthttp.StatusOK)
		case "swagger-ui-standalone-preset.js":
			return ctx.RawResponseBytes(assets.FileSwaggerUIStandalonePresetJs, fasthttp.StatusOK)
		case "swagger-ui-standalone-preset.js.map":
			return ctx.RawResponseBytes(assets.FileSwaggerUIStandalonePresetJsMap, fasthttp.StatusOK)
		case "swagger-ui.css":
			return ctx.RawResponseBytes(assets.FileSwaggerUICSS, fasthttp.StatusOK)
		case "swagger-ui.css.map":
			return ctx.RawResponseBytes(assets.FileSwaggerUICSSMap, fasthttp.StatusOK)
		case "swagger-ui.js":
			return ctx.RawResponseBytes(assets.FileSwaggerUIJs, fasthttp.StatusOK)
		case "swagger-ui.js.map":
			return ctx.RawResponseBytes(assets.FileSwaggerUIJsMap, fasthttp.StatusOK)
		}

		return ctx.TextResponse("404 page not found", fasthttp.StatusNotFound)
	}
}

func setContentType(ctx *atreugo.RequestCtx, path string) {
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
}
