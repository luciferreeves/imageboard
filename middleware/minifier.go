package middleware

import (
	"bytes"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/xml"
)

func Minifier(context *fiber.Ctx) error {
	var (
		minifyHTML = true
		minifyCSS  = true
		minifyJS   = true
		minifyXML  = true
		minifyJSON = true
	)

	m := minify.New()

	if minifyHTML {
		m.Add("text/html", &html.Minifier{
			KeepEndTags:      true,
			KeepDocumentTags: true,
		})
	}

	if minifyCSS {
		m.Add("text/css", &css.Minifier{})
	}

	if minifyJS {
		m.Add("application/javascript", &js.Minifier{})
		m.Add("application/x-javascript", &js.Minifier{})
		m.Add("text/javascript", &js.Minifier{})
		m.AddRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), &js.Minifier{})
	}

	if minifyXML {
		m.AddRegexp(regexp.MustCompile("xml$"), &xml.Minifier{})
	}

	if minifyJSON {
		m.AddRegexp(regexp.MustCompile("json$"), &json.Minifier{})
	}

	if err := context.Next(); err != nil {
		return err
	}

	statusCode := context.Response().StatusCode()
	if statusCode != fiber.StatusOK && statusCode != fiber.StatusNotModified {
		return nil
	}

	if statusCode == fiber.StatusNotModified {
		return nil
	}
	origBody := context.Response().Body()
	if len(origBody) == 0 {
		return nil
	}

	context.Response().ResetBody()

	mimetype := string(context.Response().Header.Peek("Content-Type"))

	if err := m.Minify(mimetype, context.Response().BodyWriter(), bytes.NewReader(origBody)); err != nil {
		context.Response().BodyWriter().Write(origBody)
	}

	return nil
}
