package main

import (
	"fmt"
	"log"
	"github.com/Ceinl/Go-dashboard/internal/storage"
)

func testDB() {
	db, err := storage.InitDB("file:test.db?_foreign_keys=on")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	workspaces, err := storage.GetAllWorkspaces(db)
	if err != nil {
		log.Fatalf("failed to get all workspaces: %v", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found in the database.")
		return
	}

	fmt.Println("Workspaces in the database:")
	for _, ws := range workspaces {
		fmt.Printf("ID: %s, Name: %s, Color: %s, Created At: %s\n", ws.ID, ws.Name, ws.Color, ws.CreatedAt)
	}
}