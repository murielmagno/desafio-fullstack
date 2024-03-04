package test

import (
	"bytes"
	"desafio-fullstack/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"desafio-fullstack/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(person *models.Person) error {
	args := m.Called(person)
	return args.Error(0)
}

func (m *MockDB) Update(person *models.Person) error {
	args := m.Called(person)
	return args.Error(0)
}

func TestPostPerson(t *testing.T) {
	mockDB := &MockDB{}
	mockDB.On("Create", mock.AnythingOfType("*models.Person")).Return(nil)

	router := gin.Default()
	router.POST("/account/person", func(c *gin.Context) {
		var person models.Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := mockDB.Create(&person); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
			return
		}
		c.JSON(http.StatusCreated, person)
	})
	requestBody := models.Person{}
	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/account/person", bytes.NewBuffer(jsonBody))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockDB.AssertCalled(t, "Create", mock.AnythingOfType("*models.Person"))
}

func TestPutPerson(t *testing.T) {
	mockDB := &MockDB{}
	mockDB.On("Update", mock.AnythingOfType("*models.Person")).Return(nil)
	router := gin.Default()
	router.PUT("/account/person/:id", controllers.PutPerson)
	examplePerson := models.Person{}
	jsonBody, _ := json.Marshal(examplePerson)
	req, err := http.NewRequest("PUT", "/account/person/1", bytes.NewBuffer(jsonBody))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
