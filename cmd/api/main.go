package api

import (
	"context"

	sqlrepo "github.com/PopularVote/internal/infrastructure/sql"
	"github.com/gofiber/fiber/v3"
)

func main() {

	db, err := sqlrepo.NewDB(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	// 2. Inicializar Fiber
	app := fiber.New()

	// 3. Compartilhar o pool com as rotas
	app.Use(func(c fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Here you would set up your API server and handlers, passing the database connection as needed.

}
