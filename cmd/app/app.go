package main

import (
	"database/sql"
	"os"

	"dev-challenge/internal/cell"
	"dev-challenge/internal/database"
	"dev-challenge/internal/router"
	"dev-challenge/internal/sheet"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		panic("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	server := App(db)

	log.Println("Starting server...")
	server.ListenAndServe()
}

func App(db *sql.DB) *http.Server {
	// for the sake of simpicity here we do manual dependecy injection
	cellRepo := database.NewCellRepository(db)
	cellRepo.CreateTableIfNotExists()

	cellService := cell.NewService(cellRepo)
	sheetService := sheet.NewService(cellService)

	router := router.New(sheetService, cellService)

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
