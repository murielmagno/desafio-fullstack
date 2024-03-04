package controllers

import (
	"desafio-fullstack/database"
	"desafio-fullstack/models"
	"desafio-fullstack/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

// @Summary Create Person
// @Description Create a new person
// @ID create-person
// @Accept json
// @Produce json
// @Tags Person
// @Param person body models.Person true "Person object"
// @Success 201 {object} models.Person "Success"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /account/person [post]
func PostPerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	person.Password = string(hashedPassword)
	if err := database.DB.Create(&person).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}
	c.JSON(http.StatusCreated, person)
}

// @Summary Get Person
// @Description Get a person by ID
// @ID get-person
// @Accept json
// @Produce json
// @Tags Person
// @Param id path string true "Person ID"
// @Success 200 {object} models.Person "Successfully retrieved person"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /person/{id} [get]
func GetPerson(c *gin.Context) {
	var person models.Person
	person = FindPerson(c)
	c.JSON(http.StatusOK, person)
}

// @Summary Update Person
// @Description Update an existing person
// @ID update-person
// @Accept json
// @Produce json
// @Tags Person
// @Param id path string true "Person ID"
// @Param body body models.Person true "Person object"
// @Success 200 {object} models.Person "Success"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /account/person/{id} [put]
func PutPerson(c *gin.Context) {
	var person models.Person
	person = FindPerson(c)
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	person.UpdatedAt = time.Now()
	database.DB.Model(&person).UpdateColumns(person)
	c.JSON(http.StatusOK, person)
}

func FindPerson(c *gin.Context) models.Person {
	var person models.Person
	idString := c.Params.ByName("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	result := database.DB.First(&person, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
	}
	return person
}

// @Summary Add Friend
// @Description Add friends to a person
// @ID add-friend
// @Accept json
// @Produce json
// @Tags Person
// @Param id path string true "Person ID"
// @Success 200 "Success"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /account/person/friend/{id} [post]
func AddFriend(c *gin.Context) {
	var person models.Person
	person = FindPerson(c)
	var requestBody struct {
		Friends []string `json:"friends"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	existingFriends, err := database.GetPersonFriendsByID(uuid.UUID(person.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing friends"})
		return
	}

	existingFriendsMap := make(map[string]bool)
	for _, friend := range existingFriends {
		id := friend["id"].(string)
		existingFriendsMap[id] = true
	}
	for _, friendID := range requestBody.Friends {
		if existingFriendsMap[friendID] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend ID already exists"})
			return
		}
		id, err := uuid.Parse(friendID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
			return
		}
		var friend models.Person
		result := database.DB.First(&friend, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
			return
		}
		person.Friends = append(person.Friends, &friend)
	}

	if err := database.DB.Save(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save changes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friends added successfully"})
}

func FindFriends(c *gin.Context) {
	var person models.Person
	person = FindPerson(c)
	result, err := database.GetPersonFriendsByID(uuid.UUID(person.ID))
	if err != nil {
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}
	c.JSON(http.StatusOK, result)
}

// @Summary Create Card
// @Description Create a new card for a person
// @ID create-card
// @Accept json
// @Produce json
// @Tags Card
// @Param id path string true "Person ID"
// @Param body body models.Card true "Card details"
// @Success 201 {object} models.Card "Created"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /account/card/{id} [post]
func PostCard(c *gin.Context) {
	var card models.Card
	id := c.Params.ByName("id")
	person, _ := database.GetPersonById(id)
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateCard(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card.UserId = person.ID
	card.Date = time.Now()
	if err := database.DB.Create(&card).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pan already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}
	c.JSON(http.StatusCreated, card)
}

// @Summary Get Cards
// @Description Get all cards of a person
// @ID get-cards
// @Accept json
// @Produce json
// @Tags Card
// @Param id path string true "Person ID"
// @Success 200 {array} models.Card "Success"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /account/cards/{id} [get]
func GetCards(c *gin.Context) {
	var cards []models.Card
	var cardsF []gin.H
	id := c.Params.ByName("id")
	cards, _ = database.GetCardsByUserId(id)
	cardsF = utils.FormatCardData(cards)
	c.JSON(http.StatusCreated, cardsF)
}

func validateCard(card *models.Card) error {
	if err := utils.ValidateSecurityCode(card.SecurityCode); err != nil {
		return err
	}
	if err := utils.ValidateExpiryMonth(card.ExpiryMM); err != nil {
		return err
	}
	if err := utils.ValidateExpiryYear(card.ExpiryYYYY); err != nil {
		return err
	}
	if err := utils.ValidateExpiryDate(card.ExpiryYYYY, card.ExpiryMM); err != nil {
		return err
	}
	return nil
}
