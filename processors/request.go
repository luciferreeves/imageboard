package processors

import (
	"github.com/gofiber/fiber/v2"
)

type QueryParam struct {
	Key   string
	Value string
}

type Request struct {
	Path        string
	Method      string
	Query       []QueryParam
	Params      []QueryParam
	QueryString string
	IP          string
	URL         string
}

func RequestContextProcessor(ctx *fiber.Ctx) error {
	queryParams := []QueryParam{}
	for k, v := range ctx.Queries() {
		queryParams = append(queryParams, QueryParam{Key: k, Value: v})
	}

	routeParams := []QueryParam{}
	for k, v := range ctx.AllParams() {
		routeParams = append(routeParams, QueryParam{Key: k, Value: v})
	}

	request := Request{
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
