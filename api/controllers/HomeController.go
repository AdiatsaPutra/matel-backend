package controllers

import (
	"matel/payloads"
	"matel/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {

	kendaraanTotal, err := repository.GetKendaraanTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Leasing not found", http.StatusOK)
		return
	}

	leasingTotal, err := repository.GetLeasingTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Leasing not found", http.StatusOK)
		return
	}

	userTotal, err := repository.GetUserTotalInfo(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "User not found", http.StatusOK)
		return
	}

	data := make(map[string]interface{})
	data["leasing"] = leasingTotal
	data["kendaraan"] = kendaraanTotal
	data["trial_members"] = userTotal.TrialMembers
	data["premium_members"] = userTotal.PremiumMembers
	data["expired_members"] = userTotal.ExpiredMembers

	payloads.HandleSuccess(c, data, "Data found", http.StatusOK)
}
