package controllers

import (
	"motor/payloads"
	"motor/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {

	leasingTotal, err := repository.GetLeasingTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Leasing not found", http.StatusOK)
		return
	}

	userTotal, err := repository.GetUserTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "User not found", http.StatusOK)
		return
	}

	data := make(map[string]interface{})
	data["leasing"] = leasingTotal
	data["user"] = userTotal

	payloads.HandleSuccess(c, data, "Data found", http.StatusOK)
}
