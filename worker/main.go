package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	LoadConfig()
	connectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello World!"})
	})

	app.Post("/api/record", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user data"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := DB.Collection("users").InsertOne(ctx, user)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "user added successfully", "data": user})

	})

	app.Get("/api/record", func(c *fiber.Ctx) error {
		var records []User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := DB.Collection("users").Find(ctx, bson.M{})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if err := cursor.All(ctx, &records); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal sever error"})
		}

		return c.Status(http.StatusOK).JSON(records)

	})

	PORT := AppConfig.PORT
	if PORT == "" {
		PORT = ":8080"
	}

	if err := app.Listen(PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
