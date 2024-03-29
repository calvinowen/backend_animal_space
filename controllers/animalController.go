package controllers

import (
	"net/http"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateAnimal(c *gin.Context) {
	// single file

	var body struct {
		Name        string
		Image       string
		Gender      string
		Type        string
		Age         int
		Month       int
		Description string
		Quantity    int
		CategoryID  int
		ShelterID   int
		Price       int
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	animal := models.Animal{
		Name:        body.Name,
		Image:       body.Image,
		Gender:      body.Gender,
		Type:        body.Type,
		Age:         body.Age,
		Month:       body.Month,
		Description: body.Description,
		Quantity:    1,
		Status:      "",
		CategoryID:  body.CategoryID,
		ShelterID:   body.ShelterID,
		Price:       body.Price,
	}
	resultInsert := initializers.DB.Create(&animal)

	if resultInsert.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "image": result.Location,
		"messege": "success create animal",
	})
}

func GetAllAnimalByShelter(c *gin.Context) {
	var animal []models.Animal
	shelter_id := c.Query("shelter_id")
	result := initializers.DB.Where("shelter_id = ?", shelter_id).Find(&animal)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data shelter",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    animal,
	})

}

func GetAllAnimalAdmin(c *gin.Context) {
	var animal []models.Animal
	result := initializers.DB.Preload("Shelter").Find(&animal)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    animal,
	})

}

func GetAnimal(c *gin.Context) {
	var animal models.Animal
	id := c.Param("id")
	result := initializers.DB.Where("id = ?", id).First(&animal)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data animal",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    animal,
	})

}

func UpdateAnimal(c *gin.Context) {

	var body struct {
		Name        string
		Image       string
		Gender      string
		Type        string
		Age         int
		Description string
		Quantity    int
		CategoryID  int
		ShelterID   int
		Price       int
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	id := c.Param("id")

	resultUpdate := initializers.DB.Exec(`
		UPDATE animals 
		SET name = ?, 
			gender = ?, 
			type = ?, 
			age = ?, 
			description = ?, 
			category_id = ?, 
			quantity = ?, 
			image = ?, 
			updated_at = ?,
			price = ?
		WHERE id = ?`,
		body.Name,
		body.Gender,
		body.Type,
		body.Age,
		body.Description,
		body.CategoryID,
		body.Quantity,
		body.Image,
		time.Now(),
		body.Price,
		id,
	)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update table animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "ssuccess update data",
	})

}

func DeleteAnimal(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	var animal models.Animal
	// resultUpdate := initializers.DB.Exec(`
	// 	DELETE FROM animals
	// 	WHERE id = ?`,
	// 	id,
	// )
	resultUpdate := initializers.DB.Where("id = ?", id).Delete(&animal)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delet animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success delete data",
	})

}

func GetAllAnimal(c *gin.Context) {
	var animal []models.GetAllAnimal
	search := c.Query("search")
	category := c.Query("category")
	from_age := c.Query("from_age")
	to_age := c.Query("to_age")
	from_month := c.Query("from_month")
	to_month := c.Query("to_month")
	var queryRange string
	var queryRangeMonth string
	var querySearch string
	var queryLogic string
	var queryWhere string
	var queryCategory string
	if category == "" {
		category = "%"
	}
	queryCategory = "and LOWER(b.name) LIKE LOWER('" + category + "')"

	if from_age == "" && to_age != "" {
		from_age = to_age
	}
	if to_age == "" && from_age != "" {
		to_age = from_age
	}
	if to_age != "" || from_age != "" {
		queryRange = " a.age Between " + from_age + " AND " + to_age + " "
	}

	if from_month == "" && to_month != "" {
		from_month = to_month
	}
	if to_month == "" && from_month != "" {
		to_month = from_month
	}
	if to_month != "" || from_month != "" {
		queryRangeMonth = "AND a.month Between " + from_month + " AND " + to_month + " "
	}

	if search != "" {
		querySearch = " LOWER(a.name) LIKE LOWER('%" + search + "%') OR LOWER(a.type) LIKE LOWER('%" + search + "%') "
	}
	if queryRange != "" && querySearch != "" {
		queryLogic = "AND"

	}
	if queryRange != "" || querySearch != "" {
		queryWhere = "AND"
	}

	result := initializers.DB.Raw(` 
		SELECT 
			a.id, 
			a.name, 
			a.gender, 
			a.type, 
			a.age as age, 
			a.description, 
			a.quantity, 
			a.status, 
			a.price, 
			a.image, 
			b.id as category_id, 
			b.name as category_name, 
			c.id as shelter_id,
			c.name as shelter_name,
			a.created_at 
		FROM animals a 
		JOIN categories b 
			on a.category_id = b.id 
			` + queryCategory + `
		JOIN shelters c
			on a.shelter_id = c.id
		WHERE a.deleted_at is NULL and a.status NOT LIKE 'sold' ` + queryWhere + ` ` + querySearch + ` ` + queryLogic + ` ` + queryRange + ` ` + queryRangeMonth +
		` ORDER BY a.created_at DESC
	`).Scan(&animal)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    animal,
	})

}

func GetSingelAnimal(c *gin.Context) {
	var animal models.GetAllAnimal
	id := c.Param("id")

	result := initializers.DB.Raw(` 
		SELECT 
			a.id, 
			a.name, 
			a.gender, 
			a.type, 
			a.age as age,
   			a.month as month, 
			a.description, 
			a.quantity, 
			a.status, 
			a.price, 
			a.image, 
			b.id as category_id, 
			b.name as category_name, 
			c.id as shelter_id,
			c.name as shelter_name,
			c.phone as shelter_phone,
			c.address as shelter_address,
			a.created_at 
		FROM animals a 
		JOIN categories b 
			on a.category_id = b.id 
		JOIN shelters c
			on a.shelter_id = c.id
		WHERE a.id = ?
	`, id).Scan(&animal)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    animal,
	})

}
