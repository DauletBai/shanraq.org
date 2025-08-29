package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"shanraq.org/internal/store"
)

func main() {
	log.Println("Starting database seeder...")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Clearing old data...")
	if _, err := db.Exec(`TRUNCATE TABLE awards, competitions, categories, countries RESTART IDENTITY CASCADE`); err != nil {
		log.Fatalf("Failed to clear tables: %v", err)
	}

	log.Println("Seeding countries...")
	countries := store.GetCountrySeedData()
	countryIdMap := make(map[string]string)
	for _, country := range countries {
		var countryID string
		err := db.QueryRow("INSERT INTO countries (name, population, flag_code) VALUES ($1, $2, $3) RETURNING id",
			country.Name, country.Population, country.FlagCode).Scan(&countryID)
		if err != nil {
			log.Fatalf("Failed to insert country %s: %v", country.Name, err)
		}
		countryIdMap[country.Name] = countryID
	}
	log.Printf("Successfully seeded %d countries.", len(countries))

	log.Println("Seeding competition data...")
	competitionData := store.GetCompetitionSeedData()

	for categoryName, competitions := range competitionData {
		var categoryID string
		err := db.QueryRow("INSERT INTO categories (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = $1 RETURNING id", categoryName).Scan(&categoryID)
		if err != nil {
			log.Fatalf("Failed to get or create category %s: %v", categoryName, err)
		}

		for competitionName, medals := range competitions {
			yearStr := strings.Fields(competitionName)[len(strings.Fields(competitionName))-1]
			year, _ := strconv.Atoi(yearStr)

			var competitionID string
			err := db.QueryRow("INSERT INTO competitions (category_id, name, year) VALUES ($1, $2, $3) RETURNING id", categoryID, competitionName, year).Scan(&competitionID)
			if err != nil {
				log.Fatalf("Failed to insert competition %s: %v", competitionName, err)
			}

			for _, medal := range medals {
				countryID, ok := countryIdMap[medal.CountryName]
				if !ok {
					log.Printf("Could not find country %s in map, skipping.", medal.CountryName)
					continue
				}
				_, err = db.Exec(`
					INSERT INTO awards (competition_id, country_id, gold_medals, silver_medals, bronze_medals)
					VALUES ($1, $2, $3, $4, $5)`,
					competitionID, countryID, medal.Gold, medal.Silver, medal.Bronze)
				if err != nil {
					log.Fatalf("Failed to insert award for %s in %s: %v", medal.CountryName, competitionName, err)
				}
			}
			log.Printf("Successfully seeded %d records for %s", len(medals), competitionName)
		}
	}
	log.Println("Database seeding completed successfully!")
}