package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"actiondelta/internal/model"
	"actiondelta/internal/repository"
)

func GetAdminStats(c *gin.Context) {
	cursor, err := repository.DB().Collection("users").Find(c, bson.M{})
	if err != nil {
		respond(c, http.StatusInternalServerError, "query error", nil)
		return
	}
	var userInfoList []model.User
	if err := cursor.All(c, &userInfoList); err != nil {
		respond(c, http.StatusInternalServerError, "decode error", nil)
		return
	}
	respond(c, http.StatusOK, "success", gin.H{"users": userInfoList})
}
