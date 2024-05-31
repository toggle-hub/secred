package main

import (
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/pkg/api"
	"github.com/xsadia/secred/pkg/database"
	"github.com/xsadia/secred/pkg/utils"
)

const defaultPort = "6969"

func main() {
	godotenv.Load()
	port := utils.Or(os.Getenv("PORT"), defaultPort)

	// db, err := database.New(
	// 	"localhost",
	// 	"root",
	// 	"root",
	// 	"secred",
	// 	"disable",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// database.ConfigDB(db)
	// if err := db.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := database.Migrate(db); err != nil {
	// 	log.Fatal(err)
	// }
	_, err := database.GetInstance()
	if err != nil {
		log.Fatal(err)
	}

	app := api.New()
	app.Listen(":" + port)
}
