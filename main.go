package main

import (
	"context"
	"fmt"
	"log"

	// "net/http"
	"os"

	// "github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

// NOTE: Making db a global variable simplifies this example.
// In production, you’d avoid the global variable,
// such as by passing the variable to functions that need it or by wrapping it in a struct.
var db *pgx.Conn
var err error

type Album struct {
	ID        int64
	Title     string
	Artist    string
	Price     float32
	CreatedAt pgtype.Date
	UpdatedAt pgtype.Date
}

func main() {
	// var config = "user=vscode password=vscode host=localhost port=5432 dbname=recordings sslmode=verify-ca pool_max_conns=10"
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	pingErr := db.Ping(context.Background())
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
}

func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("1 albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price, &alb.CreatedAt, &alb.UpdatedAt); err != nil {
			return nil, fmt.Errorf("2 albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("3 albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}
