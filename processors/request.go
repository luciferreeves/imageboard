package processors

import (
	"imageboard/config"

	"github.com/gofiber/fiber/v2"
)

func RequestContextProcessor(ctx *fiber.Ctx) error {
	queryParams := []config.QueryParam{}
	queryArgs := ctx.Request().URI().QueryArgs()
	queryArgs.VisitAll(func(key, value []byte) {
		queryParams = append(queryParams, config.QueryParam{Key: string(key), Value: string(value)})
	})

	routeParams := []config.QueryParam{}
	for k, v := range ctx.AllParams() {
		routeParams = append(routeParams, config.QueryParam{Key: k, Value: v})
	}

	request := config.Request{
		Path:        ctx.Path(),
		Method:      ctx.Method(),
		Query:       queryParams,
		Params:      routeParams,
		QueryString: string(ctx.Request().URI().QueryString()),
		IP:          ctx.IP(),
		URL:         ctx.OriginalURL(),
	}

	ctx.Locals("Request", request)
	return ctx.Next()
}
