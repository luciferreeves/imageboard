package shortcuts

import (
	"reflect"
	"strings"

	"maps"

	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, name string, bind any) error {
	finalData := fiber.Map{}

	ctx.Context().VisitUserValues(func(key []byte, value any) {
		finalData[string(key)] = value
	})

	if bind != nil {
		switch v := bind.(type) {
		case fiber.Map:
			maps.Copy(finalData, v)
		case map[string]any:
			maps.Copy(finalData, v)
		default:
			structData := structToMap(bind)
			maps.Copy(finalData, structData)
		}
	}

	return ctx.Render(name, finalData)
}

func RenderWithStatus(ctx *fiber.Ctx, name string, bind any, statusCode int) error {
	ctx.Status(statusCode)
	return Render(ctx, name, bind)
}

func structToMap(obj any) map[string]any {
	result := make(map[string]any)

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := range v.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		key := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			if commaIdx := strings.Index(tag, ","); commaIdx > 0 {
				key = tag[:commaIdx]
			} else if commaIdx == -1 {
				key = tag
			}
		}

		result[key] = v.Field(i).Interface()
	}

	return result
}
