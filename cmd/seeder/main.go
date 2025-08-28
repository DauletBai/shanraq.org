package main

import (
	"database/sql"
	//"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
	"shanraq.org/internal/store"
)

func main() {
	log.Println("Starting database seeder...")

	// 1. Подключаемся к БД (берем URL из переменных окружения)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	// 2. Очищаем старые данные (чтобы избежать дубликатов при повторном запуске)
	log.Println("Clearing old data...")
	if _, err := db.Exec(`TRUNCATE TABLE awards, competitions, categories, countries RESTART IDENTITY CASCADE`); err != nil {
		log.Fatalf("Failed to clear tables: %v", err)
	}

	// 3. Загружаем данные по странам
	log.Println("Seeding countries...")
	countries := store.GetCountrySeedData()
	for _, country := range countries {
		_, err := db.Exec("INSERT INTO countries (name, population, flag_code) VALUES ($1, $2, $3)",
			country.Name, country.Population, country.FlagCode)
		if err != nil {
			log.Fatalf("Failed to insert country %s: %v", country.Name, err)
		}
	}

	log.Printf("Successfully seeded %d countries.", len(countries))

	// В будущем здесь будет код для загрузки данных об Олимпиадах

	log.Println("Database seeding completed successfully!")
}