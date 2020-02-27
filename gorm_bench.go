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

type City struct {
	Name        string
	CountryCode string
	District    string
	Population  int32
}

func (City) TableName() string {
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

func sql_query(b *testing.B) {
	db, err := sql.Open("postgres", connString())
	panicIf(err)
	defer db.Close()

	query := "SELECT name, country_code, district, population FROM city"
	rows, err := db.Query(query)
	panicIf(err)

	cities := make([]City, 0)
	for rows.Next() {
		var city City
		err = rows.Scan(&city.Name, &city.CountryCode, &city.District, &city.Population)
		panicIf(err)
		cities = append(cities, city)
	}

	if len(cities) != cityCount {
		panic(fmt.Sprintf("invalid result: %d", len(cities)))
	}
}

func pgx_query(b *testing.B) {
	conn, err := pgx.Connect(context.Background(), connString())
	panicIf(err)
	defer conn.Close(context.Background())

	query := "SELECT name, country_code, district, population FROM city"
	rows, err := conn.Query(context.Background(), query)
	panicIf(err)
	defer rows.Close()

	cities := make([]City, 0)
	for rows.Next() {
		var city City
		err = rows.Scan(&city.Name, &city.CountryCode, &city.District, &city.Population)
		panicIf(err)
		cities = append(cities, city)
	}

	if len(cities) != cityCount {
		panic(fmt.Sprintf("invalid result: %d", len(cities)))
	}
}

func gorm_query(b *testing.B) {
	db, err := gorm.Open("postgres", connString())
	panicIf(err)
	defer db.Close()

	var cities []City
	err = db.Find(&cities).Error
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

func main() {
	fmt.Println("pgx: ", testing.Benchmark(pgx_query))
	fmt.Println("sql: ", testing.Benchmark(sql_query))
	fmt.Println("gorm:", testing.Benchmark(gorm_query))
}
