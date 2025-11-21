package main

import (
	"fmt"
	"log"

	"github.com/Alex322322/short-url/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// TODO init config - cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO init logger - slog log/slog
	// TODO init storage - sqlite/postgres
	// TODO init router - chi
	// TODO run server - net/http
}
