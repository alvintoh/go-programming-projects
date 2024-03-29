package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alvintoh/go-programming-projects/go-fiber-crm-basic/internal/app/middleware"
	"github.com/alvintoh/go-programming-projects/go-fiber-crm-basic/internal/app/model"
	"github.com/alvintoh/go-programming-projects/go-fiber-crm-basic/internal/app/route"
	"github.com/alvintoh/go-programming-projects/go-fiber-crm-basic/internal/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func setupConfig() {
	viper.SetConfigName("application")
	viper.AddConfigPath("./resource")
	// use for development debug reference path
	viper.AddConfigPath("../../resource")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if selectedConfigEnv := os.Getenv("CONFIG_ENV"); selectedConfigEnv != "" {
		viper.SetConfigName(selectedConfigEnv)
		viper.MergeInConfig()
	}
}

func connectDatabase() {
	database.Connect()
	db := database.GetDB()
	model.SetDB(db)

	if err := db.AutoMigrate(&model.Lead{}); err != nil {
		panic("Failed to auto migrate Lead: " + err.Error())
	}
}

func setupRoutes(app *fiber.App) {
	app.Static("/swagger/ui", "./swagger/ui")
	app.Static("/swagger", "./swagger")
	app.Use(middleware.LoggingMiddleware)
	route.SetupRoutes(app)
}

func startServer(app *fiber.App) {
	addr := viper.GetString("service.host") + ":" + viper.GetString("service.port")
	log.Println("Starting server at " + addr)
	log.Fatal(app.Listen(addr))
}

func main() {
	setupConfig()
	connectDatabase()
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})
	setupRoutes(app)
	startServer(app)
}
