package controller

import (
	"actiondelta/internal/model"
	"actiondelta/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAdminStats(c *gin.Context) {
	// Get admin stats from database
	cursor, err := repository.DB().Collection("users").Find(c, bson.M{})
	if err != nil {
		panic(err)
	}

	var userInfoList []model.User

	if err := cursor.All(c, &userInfoList); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, userInfoList)
}
