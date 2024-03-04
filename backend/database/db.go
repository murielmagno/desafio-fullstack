package database

import (
	"desafio-fullstack/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=postgres user=root password=root dbname=root port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Opção para criar automaticamente a tabela se ela ainda não existir
	err = DB.AutoMigrate(&models.Person{}, &models.Card{}, &models.Transfer{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func GetCardsByUserId(id string) ([]models.Card, error) {
	var cards []models.Card
	if err := DB.Where("user_id = ?", id).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func GetPersonFriendsByID(personID uuid.UUID) ([]map[string]interface{}, error) {
	var rawResults []map[string]interface{}
	if err := DB.Table("person_friends").
		Select(" people.* ").
		Joins("JOIN people ON people.id = person_friends.friend_id").
		Where("person_friends.person_id = ?", personID).
		Find(&rawResults).Error; err != nil {
		return nil, err
	}
	var results []map[string]interface{}
	for _, res := range rawResults {
		result := make(map[string]interface{})
		result["id"] = res["id"]
		result["first_name"] = res["first_name"]
		result["last_name"] = res["last_name"]
		result["birthday"] = res["birthday"]
		result["user_name"] = res["user_name"]
		results = append(results, result)
	}
	return results, nil

}

func GetPersonById(personIDStr string) (models.Person, error) {
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		return models.Person{}, err
	}
	var person models.Person
	result := DB.First(&person, personID)
	if result.Error != nil {
		return person, result.Error
	}
	return person, nil
}

func GetFriendIDs(senderID string) ([]string, error) {
	var friendIDs []string

	// Consulta SQL personalizada para buscar os IDs dos amigos do remetente
	if err := DB.Table("person_friends").
		Select(" people.id ").
		Joins("JOIN people ON people.id = person_friends.friend_id").
		Where("person_friends.person_id = ?", senderID).
		Group("people.id").
		Pluck("friend_ids", &friendIDs).Error; err != nil {
		return nil, err
	}

	return friendIDs, nil
}

func GetCardById(id string) (models.Card, error) {
	var card models.Card
	if err := DB.Where("pan = ?", id).First(&card).Error; err != nil {
		return models.Card{}, err
	}
	return card, nil
}

func FindTransfers(userID string) ([]models.Transfer, error) {
	var transfers []models.Transfer
	if err := DB.Where("sender_id = ?", userID).Find(&transfers).Error; err != nil {
		return []models.Transfer{}, err
	}
	return transfers, nil
}
