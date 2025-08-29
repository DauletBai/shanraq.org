package service

import (
	"math"
	"shanraq.org/internal/models"
	"sort"
)

const PageSize = 10

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
func (s *IndexService) GetRankedCountries(page int) (models.PaginatedRanking, error) {
	stats, err := s.store.GetAllCountryStats()
	if err != nil {
		return models.PaginatedRanking{}, err
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

	// НАЧАЛО ЛОГИКИ ПАГИНАЦИИ
    totalItems := len(rankedCountries)
    totalPages := int(math.Ceil(float64(totalItems) / float64(PageSize)))
    
    if page < 1 {
        page = 1
    }
    if page > totalPages {
        page = totalPages
    }

    start := (page - 1) * PageSize
    end := start + PageSize
    if end > totalItems {
        end = totalItems
    }
    
    // "Нарезаем" нужный кусок от общего отсортированного списка
    paginatedCountries := rankedCountries[start:end]

    // Присваиваем ранг только для стран на текущей странице
    for i := range paginatedCountries {
        paginatedCountries[i].Rank = start + i + 1
    }

    pagination := models.PaginationData{
        CurrentPage: page,
        TotalPages:  totalPages,
        HasNext:     page < totalPages,
        HasPrev:     page > 1,
        NextPage:    page + 1,
        PrevPage:    page - 1,
    }

    return models.PaginatedRanking{
        Countries:  paginatedCountries,
        Pagination: pagination,
    }, nil
}

func (s *IndexService) GetCountryDetails(name string) (models.CountryDetails, error) {
	details, err := s.store.GetCountryDetailsByName(name)
	if err != nil {
		return models.CountryDetails{}, err
	}

	// Инициализируем карту со всеми нашими НОВЫМИ категориями
	chartData := map[string]int{
		"Sport":                    0,
		"Science":      0,
		"IT & Engineering": 0,
		"Arts & Culture":           0,
		"Professional Skills":      0,
	}

	// Проходим по наградам страны и суммируем очки только для тех категорий, что у нее есть
	for _, award := range details.Awards {
		points := (award.Gold * 3) + (award.Silver * 2) + (award.Bronze * 1)
		// Проверяем, что категория есть в нашем списке, чтобы избежать ошибок
		if _, ok := chartData[award.CategoryName]; ok {
			chartData[award.CategoryName] += points
		}
	}

	details.ChartData = chartData 

	return details, nil
}

// GetAllRankedCountries возвращает полный список стран для карты
func (s *IndexService) GetAllRankedCountries() ([]models.RankedCountry, error) {
    // Эта логика дублирует GetRankedCountries, но без пагинации
    stats, err := s.store.GetAllCountryStats()
    if err != nil {
        return nil, err
    }

    var rankedCountries []models.RankedCountry
    for _, stat := range stats {
        totalPoints := (stat.TotalGold * 3) + (stat.TotalSilver * 2) + (stat.TotalBronze * 1)
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

    sort.Slice(rankedCountries, func(i, j int) bool {
        return rankedCountries[i].HPI > rankedCountries[j].HPI
    })

    for i := range rankedCountries {
        rankedCountries[i].Rank = i + 1
    }

    return rankedCountries, nil
}