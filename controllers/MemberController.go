package controllers

import (
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"motor/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateMember(c *gin.Context) {

	memberId := c.Param("id")

	var body struct {
		Status uint `json:"status"`
	}

	c.ShouldBindJSON(&body)

	ID, _ := strconv.ParseUint(memberId, 10, 64)

	uintID := uint(ID)

	var member = models.Member{
		ID:     uintID,
		Status: body.Status,
	}

	findMemberFromDB, err := repository.UpdateMember(c, member)

	if err != nil {
		exceptions.AppException(c, "Member Not Found")
		return
	}

	payloads.HandleSuccess(c, findMemberFromDB, "Member updated", http.StatusOK)
}
