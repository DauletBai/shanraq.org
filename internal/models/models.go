package models

import "github.com/google/uuid"

// Country остается без изменений
type Country struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Population int64     `json:"population"`
	FlagCode   string    `json:"flag_code"`
}

// RankedCountry - это структура для итогового рейтинга
type RankedCountry struct {
	Rank        int
	Country     Country
	TotalPoints int
	HPI         float64 
}

// AwardInfo - структура для отображения детальной информации о награде
type AwardInfo struct {
	CompetitionName string
	CategoryName    string
	Year            int
	Gold            int
	Silver          int
	Bronze          int
}

// CountryDetails - полная информация о стране для детальной страницы
type CountryDetails struct {
	Country Country
	Awards  []AwardInfo
	ChartData map[string]int
}