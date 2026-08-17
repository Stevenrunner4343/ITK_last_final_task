package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func TestInsertBatch(t *testing.T) {
// 	connStr := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable",
// 		os.Getenv("POSTGRES_USER"),
// 		os.Getenv("POSTGRES_PASSWORD"),
// 		os.Getenv("POSTGRES_DB"))

// 	db, err := pgxpool.New(context.Background(), connStr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	testUserID := 999999 // РАНДОМНЫЙ ID (нАДО ПОСМОТРЕТЬ МАКИМАЛЬНЫЙ В бд И взять больше)
// 	defer db.Exec(context.Background(), "DELETE FROM operations WHERE user_id = $1", testUserID)

// 	batch := []Data{
// 		{Id: testUserID, Operation: "deposit", Amount: 100},
// 		{Id: testUserID, Operation: "deposit", Amount: 50},
// 	}

// 	if err := insertBatch(db, batch); err != nil {
// 		t.Fatalf("Оштбка при ЗАписи в батча а бд: %v", err)
// 	}

// 	var count int
// 	db.QueryRow(context.Background(),
// 		"SELECT COUNT(*) FROM operations WHERE user_id = $1", testUserID).Scan(&count)

// 	if count != len(batch) {
// 		t.Errorf("Ожидалось %d строк, получено %d", len(batch), count)
// 	}
// }
