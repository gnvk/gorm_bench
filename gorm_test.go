package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const cityCount = 4079

type city struct {
	Name        string
	CountryCode string
	District    string
	Population  int32
}

func (city) TableName() string {
	return "city"
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func connString() string {
	host := getenv("PGHOST", "localhost")
	port := getenv("PGPORT", "35432")
	user := getenv("PGUSER", "world")
	password := getenv("PGPASSWORD", "world123")
	dbname := getenv("PGDATABASE", "world-db")

	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		host, port, user, dbname, password)
}

func sqlQuery(db *sql.DB) {
	query := "SELECT name, country_code, district, population FROM city"
	rows, err := db.Query(query)
	panicIf(err)

	cities := make([]city, 0)
	for rows.Next() {
		var c city
		err = rows.Scan(&c.Name, &c.CountryCode, &c.District, &c.Population)
		panicIf(err)
		cities = append(cities, c)
	}

	if len(cities) != cityCount {
		panic(fmt.Sprintf("invalid result: %d", len(cities)))
	}
}

func pgxQuery(conn *pgx.Conn) {
	query := "SELECT name, country_code, district, population FROM city"
	rows, err := conn.Query(context.Background(), query)
	panicIf(err)
	defer rows.Close()

	cities := make([]city, 0)
	for rows.Next() {
		var c city
		err = rows.Scan(&c.Name, &c.CountryCode, &c.District, &c.Population)
		panicIf(err)
		cities = append(cities, c)
	}

	if len(cities) != cityCount {
		panic(fmt.Sprintf("invalid result: %d", len(cities)))
	}
}

func gormQuery(db *gorm.DB) {
	var cities []city
	err := db.Find(&cities).Error
	panicIf(err)

	if len(cities) != cityCount {
		panic(fmt.Sprintf("invalid result: %d", len(cities)))
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func BenchmarkSqlQuery(b *testing.B) {
	db, err := sql.Open("postgres", connString())
	panicIf(err)
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sqlQuery(db)
	}
}

func BenchmarkPgxQuery(b *testing.B) {
	conn, err := pgx.Connect(context.Background(), connString())
	panicIf(err)
	defer conn.Close(context.Background())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pgxQuery(conn)
	}
}

func BenchmarkGormQuery(b *testing.B) {
	db, err := gorm.Open("postgres", connString())
	panicIf(err)
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gormQuery(db)
	}
}
