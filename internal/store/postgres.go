package store

import (
	"database/sql"
	"shanraq.org/internal/models"
	"shanraq.org/internal/service"
)

// PostgresStore содержит пул соединений с БД
type PostgresStore struct {
	DB *sql.DB
}

// CountryStats - временная структура для получения агрегированных данных из БД
type CountryStats struct {
	models.Country
	TotalGold   int
	TotalSilver int
	TotalBronze int
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

// GetAllCountryStats получает все страны с суммой их наград
func (s *PostgresStore) GetAllCountryStats() ([]service.CountryStats, error) {
	query := `
		SELECT
			c.id, c.name, c.population, c.flag_code,
			COALESCE(SUM(a.gold_medals), 0) as total_gold,
			COALESCE(SUM(a.silver_medals), 0) as total_silver,
			COALESCE(SUM(a.bronze_medals), 0) as total_bronze
		FROM
			countries c
		LEFT JOIN
			awards a ON c.id = a.country_id
		GROUP BY
			c.id
		ORDER BY
			c.name ASC;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем срез того типа, который ожидает вызывающая сторона (сервис)
	var serviceStats []service.CountryStats

	for rows.Next() {
		// Используем локальную структуру для сканирования данных
		var localStat CountryStats
		if err := rows.Scan(
			&localStat.ID, &localStat.Name, &localStat.Population, &localStat.FlagCode,
			&localStat.TotalGold, &localStat.TotalSilver, &localStat.TotalBronze,
		); err != nil {
			return nil, err
		}
		
		// Явно преобразуем локальный тип в тип сервиса и добавляем в результат
		serviceStats = append(serviceStats, service.CountryStats{
			Country:     localStat.Country,
			TotalGold:   localStat.TotalGold,
			TotalSilver: localStat.TotalSilver,
			TotalBronze: localStat.TotalBronze,
		})
	}
	
	// Возвращаем корректно типизированный срез
	return serviceStats, nil
}

// GetCountryDetailsByName находит страну по имени и все ее награды
func (s *PostgresStore) GetCountryDetailsByName(name string) (models.CountryDetails, error) {
	var details models.CountryDetails

	// 1. Находим страну
	countryQuery := `SELECT id, name, population, flag_code FROM countries WHERE name = $1`
	row := s.DB.QueryRow(countryQuery, name)
	err := row.Scan(&details.Country.ID, &details.Country.Name, &details.Country.Population, &details.Country.FlagCode)
	if err != nil {
		return models.CountryDetails{}, err // Возвращаем ошибку, если страна не найдена
	}

	// 2. Находим все ее награды
	awardsQuery := `
		SELECT cat.name, comp.name, comp.year, a.gold_medals, a.silver_medals, a.bronze_medals
		FROM awards a
		JOIN competitions comp ON a.competition_id = comp.id
		JOIN categories cat ON comp.category_id = cat.id
		WHERE a.country_id = $1
		ORDER BY comp.year DESC, cat.name;
	`
	rows, err := s.DB.Query(awardsQuery, details.Country.ID)
	if err != nil {
		return models.CountryDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var award models.AwardInfo
		if err := rows.Scan(
			&award.CategoryName, &award.CompetitionName, &award.Year,
			&award.Gold, &award.Silver, &award.Bronze,
		); err != nil {
			return models.CountryDetails{}, err
		}
		details.Awards = append(details.Awards, award)
	}

	return details, nil
}