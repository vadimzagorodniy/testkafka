package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testjun/models"
	util "testjun/utils/postgres"
)

var user models.User
var users []models.User

func Update(c *gin.Context) {

	db := util.GetDbInstance()
	c.BindJSON(&user)

	if isValid(&user) {
		db.Model(&user).Updates(user)
	}

}

func Delete(c *gin.Context) {
	c.BindJSON(&user)
	db := util.GetDbInstance()
	db.Where("uid", user.UID).Delete(&user)
}

func Create(c *gin.Context) {

	db := util.GetDbInstance()
	c.BindJSON(&user)

	//user.UID = 0
	if isValid(&user) {
		newUser := db.Create(&user)
		if newUser.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something is wrong with the request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "true"})
		return
	} else {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

}

func Get(c *gin.Context) {

	db := util.GetDbInstance()

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		// set to default
		limit = 5
	}

	offset, err := strconv.Atoi(c.Query("offset"))

	if err != nil {
		// set to default
		offset = 0
	}

	user := &models.User{
		Name:        c.Query("name"),
		Surname:     c.Query("surname"),
		Patronymic:  c.Query("patronymic"),
		Age:         c.Query("age"),
		Gender:      c.Query("gender"),
		Nationality: c.Query("nationality"),
	}
	query := db.Where(user)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(limit)
	}

	query.Find(&users)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to fetch results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func isValid(u *models.User) bool {
	if u.Name != "" && u.Surname != "" {
		return true
	} else {
		return false
	}
}
