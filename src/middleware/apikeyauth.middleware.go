package middleware

import (
	"crypto/subtle"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

func unauthorizedApiKey(ctx fiber.Ctx) error {
	return ctx.Status(401).JSON(api.Response[*string, map[string][]string]{
		Status:  "error",
		Message: "Anda belum terotentikasi",
		Errors: map[string][]string{
			"auth": {
				"Gagal otorisasi API key",
			},
		},
	})
}

func ApiKeyAuth(ctx fiber.Ctx) error {
	apiKey := ctx.Get("X-API-KEY")
	appKey := environment.APP_KEY

	if appKey == "" || apiKey == "" {
		return unauthorizedApiKey(ctx)
	}

	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(appKey)) == 1 {
		return ctx.Next()
	}

	return unauthorizedApiKey(ctx)
}
