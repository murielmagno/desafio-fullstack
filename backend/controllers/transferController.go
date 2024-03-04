package controllers

import (
	"desafio-fullstack/database"
	"desafio-fullstack/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Create Transfer
// @Description Create a new transfer
// @ID create-transfer
// @Accept json
// @Produce json
// @Tags Transfer
// @Param body body models.Transfer true "Transfer Request Body"
// @Success 200 {object} models.Transfer "Successfully created transfer"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /account/transfer [post]
func PostTransfer(c *gin.Context) {
	var transfer models.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := MakeTransfer(transfer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transfer)
}

func MakeTransfer(transferRequest models.Transfer) error {
	if !AreFriends(transferRequest.SenderID, transferRequest.FriendID) {
		return errors.New("sender and friend are not friends")
	}
	var card models.Card
	card, _ = FindCardById(transferRequest.Pan)
	if card.Limit < transferRequest.TotalToTransfer {
		return errors.New("insufficient funds")
	}
	card.Limit -= transferRequest.TotalToTransfer
	if err := database.DB.Create(&transferRequest).Error; err != nil {
		return err
	}
	if err := database.DB.Save(&card).Error; err != nil {
		return err
	}
	return nil
}

func AreFriends(senderID, friendID string) bool {
	var friendIDs []string
	friendIDs, _ = database.GetFriendIDs(senderID)
	for _, id := range friendIDs {
		if id == friendID {
			return true
		}
	}
	return false
}

func FindCardById(cardID string) (models.Card, error) {
	var card models.Card
	card, _ = database.GetCardById(cardID)
	return card, nil
}

// @Summary Get Transfers
// @Description Get transfers by user ID
// @ID get-transfers
// @Accept json
// @Produce json
// @Tags Transfer
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Transfer "Successfully retrieved transfers"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /bank-statement/{user_id} [get]
func GetTransfers(c *gin.Context) {
	userID := c.Params.ByName("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	var transfers []models.Transfer
	transfers, _ = database.FindTransfers(userID)
	if transfers == nil {
		c.JSON(http.StatusNotFound, "")
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// @Summary Get transfers from friends
// @Description Get transfers from friends based on their IDs
// @Tags Transfers
// @Accept json
// @Produce json
// @Success 200 {array} models.Transfer
// @Failure 400 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /account/bank-statement [get]
func GetTransfersFriends(c *gin.Context) {
	var friendIDs []string

	if err := database.DB.Table("person_friends").
		Select("people.id").
		Joins("JOIN people ON people.id = person_friends.person_id").
		Group("people.id").
		Pluck("person_ids", &friendIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var transfers []models.Transfer
	for _, friendID := range friendIDs {
		friendTransfers, err := database.FindTransfers(friendID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transfers = append(transfers, friendTransfers...)
	}
	if transfers == nil {
		c.JSON(http.StatusNotFound, "")
	}
	c.JSON(http.StatusOK, transfers)
}
