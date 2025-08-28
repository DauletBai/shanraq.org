package service

import (
	"shanraq.org/internal/models"
	"sort"
)

type Store interface {
	GetAllCountryStats() ([]CountryStats, error)
	GetCountryDetailsByName(name string) (models.CountryDetails, error) 
}

// CountryStats - дублируем структуру для удобства
type CountryStats struct {
	models.Country
	TotalGold   int
	TotalSilver int
	TotalBronze int
}

// IndexService содержит зависимости, например, store
type IndexService struct {
	store Store
}

// NewIndexService создает новый сервис
func NewIndexService(store Store) *IndexService {
	return &IndexService{store: store}
}

// GetRankedCountries - главный метод, который выполняет все расчеты
func (s *IndexService) GetRankedCountries() ([]models.RankedCountry, error) {
	stats, err := s.store.GetAllCountryStats()
	if err != nil {
		return nil, err
	}

	var rankedCountries []models.RankedCountry

	for _, stat := range stats {
		// Шаг 1: Считаем взвешенные очки
		totalPoints := (stat.TotalGold * 3) + (stat.TotalSilver * 2) + (stat.TotalBronze * 1)

		// Шаг 2: Считаем HPI
		var hpi float64
		if stat.Population > 0 {
			hpi = (float64(totalPoints) / float64(stat.Population)) * 1000000
		}

		rankedCountries = append(rankedCountries, models.RankedCountry{
			Country:     stat.Country,
			TotalPoints: totalPoints,
			HPI:         hpi,
		})
	}

	// Шаг 3: Сортируем страны по HPI (от большего к меньшему)
	sort.Slice(rankedCountries, func(i, j int) bool {
		return rankedCountries[i].HPI > rankedCountries[j].HPI
	})

	// Шаг 4: Присваиваем ранг
	for i := range rankedCountries {
		rankedCountries[i].Rank = i + 1
	}

	return rankedCountries, nil
}

func (s *IndexService) GetCountryDetails(name string) (models.CountryDetails, error) {
	details, err := s.store.GetCountryDetailsByName(name)
	if err != nil {
		return models.CountryDetails{}, err
	}

	// Инициализируем карту со всеми нашими основными категориями и нулевыми значениями
	chartData := map[string]int{
		"Science":            0,
		"Sport":              0,
		"Arts & Culture":     0,
		"Professional Skills": 0,
	}

	// Проходим по наградам страны и суммируем очки только для тех категорий, что у нее есть
	for _, award := range details.Awards {
		points := (award.Gold * 3) + (award.Silver * 2) + (award.Bronze * 1)
		// Проверяем, что категория есть в нашем списке, чтобы избежать ошибок
		if _, ok := chartData[award.CategoryName]; ok {
			chartData[award.CategoryName] += points
		}
	}

	details.ChartData = chartData // Добавляем посчитанные данные в результат

	return details, nil
}