package store

import (
	"database/sql"
	"shanraq.org/internal/models"
)

// PostgresStore содержит пул соединений с БД
type PostgresStore struct {
	DB *sql.DB
}

// NewPostgresStore создает новый экземпляр PostgresStore
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{DB: db}
}

// GetCountries извлекает все страны из базы данных
func (s *PostgresStore) GetCountries() ([]models.Country, error) {
	// Готовим SQL-запрос
	query := `SELECT id, name, population, flag_code FROM countries ORDER BY name ASC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country

	// Проходим по всем строкам результата
	for rows.Next() {
		var c models.Country
		// Сканируем значения из строки в поля нашей структуры
		if err := rows.Scan(&c.ID, &c.Name, &c.Population, &c.FlagCode); err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	return countries, nil
}