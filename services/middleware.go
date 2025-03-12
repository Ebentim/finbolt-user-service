package services

import (
	"github.com/rs/cors"
)

func CorsMiddleware() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://finbolt-lab.vercel.app", "http://localhost:3000", "https://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours for preflight cache

	})

	return c
}
