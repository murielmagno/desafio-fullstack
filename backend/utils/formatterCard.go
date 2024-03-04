package utils

import (
	"desafio-fullstack/models"
	"github.com/gin-gonic/gin"
)

func FormatCardData(cards []models.Card) []gin.H {
	var formattedCards []gin.H
	for _, card := range cards {
		formattedCard := gin.H{
			"title":         card.Title,
			"pan":           card.PAN,
			"expiry_mm":     card.ExpiryMM,
			"expiry_yyyy":   card.ExpiryYYYY,
			"security_code": card.SecurityCode,
			"date":          card.Date.Format("2006-01-02"), // Formatar a data como "AAAA-MM-DD"
		}
		formattedCards = append(formattedCards, formattedCard)
	}
	return formattedCards
}
