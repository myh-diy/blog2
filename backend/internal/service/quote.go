package service

import (
	"math/rand"
	"time"

	"gorm.io/gorm"

	"blog-backend/internal/model"
)

func GetRandomQuotes(db *gorm.DB, limit int) []string {
	var quotes []model.Quote
	if err := db.Find(&quotes).Error; err != nil || len(quotes) == 0 {
		// Fallback: return nil so caller can provide defaults
		return nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(quotes), func(i, j int) {
		quotes[i], quotes[j] = quotes[j], quotes[i]
	})

	n := limit
	if n > len(quotes) {
		n = len(quotes)
	}

	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = quotes[i].Text
	}
	return result
}

func GetAllQuotes(db *gorm.DB) ([]model.Quote, error) {
	var quotes []model.Quote
	err := db.Order("created_at DESC").Find(&quotes).Error
	return quotes, err
}

func CreateQuote(db *gorm.DB, text string) (*model.Quote, error) {
	q := &model.Quote{Text: text}
	err := db.Create(q).Error
	return q, err
}

func DeleteQuote(db *gorm.DB, id uint) error {
	return db.Delete(&model.Quote{}, id).Error
}
