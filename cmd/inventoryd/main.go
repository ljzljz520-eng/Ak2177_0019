package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"inventorychain/internal/api"
	"inventorychain/internal/service"
	"inventorychain/internal/store"
)

func main() {
	dbPath := flag.String("db", "inventorychain.db", "bbolt database path")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	database, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	svc := service.New(database, service.FixedClock{Value: "2026-01-01T00:00:00Z"})
	fmt.Printf("inventory chain service listening on %s\n", *addr)
	if err := http.ListenAndServe(*addr, api.CommandHandler(svc)); err != nil {
		log.Fatal(err)
	}
}
