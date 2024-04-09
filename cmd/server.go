package main

import (
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/pkg/database"
	"github.com/xsadia/secred/pkg/handlers"
	"github.com/xsadia/secred/pkg/utils"
)

const defaultPort = "6969"

func main() {
	godotenv.Load()
	port := utils.Or(os.Getenv("PORT"), defaultPort)

	db, err := database.New(
		"localhost",
		"root",
		"root",
		"secred",
		"disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.ConfigDB(db)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.POST("graphql", handlers.GraphQLHandler(db))
	e.GET("/", handlers.PlaygroundHandler())

	log.Panic(e.Start(":" + port))
}
