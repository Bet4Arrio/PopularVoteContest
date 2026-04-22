package main

import (
	"context"
	"log"

	"github.com/PopularVote/config"

	// infraauth "github.com/PopularVote/internal/infrastructure/auth"

	infraauth "github.com/PopularVote/internal/infrastructure/auth"
	sqlrepo "github.com/PopularVote/internal/infrastructure/sql"
	apihandles "github.com/PopularVote/internal/interfaces/api/handles"
	apimiddleware "github.com/PopularVote/internal/interfaces/api/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg := config.Load()
	// 1. Inicializar o pool de conexões
	db, err := sqlrepo.NewDB(context.Background(), cfg.DB.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 2. Inicializar Fiber
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	// 3. Compartilhar o pool com as rotas
	app.Use(func(c fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})
	jwtSvc := infraauth.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Welcome to the Popular Vote API!")
	})
	api := app.Group("/api")
	authHandler := apihandles.NewAuthAPIHandler(jwtSvc)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)
	authGroup.Get("/me", apimiddleware.JWTProtected(jwtSvc), authHandler.Me)

	contestGroupp := api.Group("/contests", apimiddleware.JWTProtected(jwtSvc))
	contestHandler := apihandles.NewContestHandler()
	contestGroupp.Post("/", contestHandler.CreateContest)
	contestGroupp.Get("/:id", contestHandler.GetContest)
	contestGroupp.Get("/", contestHandler.ListContests)
	contestGroupp.Put("/:id", contestHandler.UpdateContest)
	contestGroupp.Post("/:id/participants", contestHandler.AddParticipant)
	contestGroupp.Delete("/:id/participants/:participantId", contestHandler.RemoveParticipant)
	contestGroupp.Get("/:id/pairs", contestHandler.ListPairtoVote)
	contestGroupp.Post("/:id/vote", contestHandler.VoteParticipant)
	// Here you would set up your API server and handlers, passing the database connection as needed.
	log.Fatal(app.Listen(":" + cfg.App.Port))

}
