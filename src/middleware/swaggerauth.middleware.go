package middleware

import (
	"encoding/base64"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/mapping/api"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func unauthorizedSwagger(ctx fiber.Ctx) error {
	ctx.Set("WWW-Authenticate", "Basic realm=Restricted Area")
	return ctx.Status(401).JSON(api.Response[*string, map[string][]string]{
		Status:  "error",
		Message: "Anda belum terotentikasi",
		Errors: map[string][]string{
			"auth": {
				"Gagal otorisasi Swagger",
			},
		},
	})
}

func SwaggerBasicAuth(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return unauthorizedSwagger(ctx)
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return unauthorizedSwagger(ctx)
	}

	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")

	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		log.Printf("Error decoding Basic Auth string: %v", err)
		return unauthorizedSwagger(ctx)
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return unauthorizedSwagger(ctx)
	}

	// fetch basic creds
	username := creds[0]
	password := creds[1]

	if username == environment.SWAGGER_USERNAME && password == environment.SWAGGER_PASSWORD {
		return ctx.Next()
	}

	return unauthorizedSwagger(ctx)
}
