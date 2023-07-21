package controllers

import (
	"matel/payloads"
	"matel/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetHome(c *gin.Context) {

	leasing := c.Query("leasing")
	cabang := c.Query("cabang")

	kendaraanTotal, err := repository.GetKendaraanTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Something went wrong", http.StatusOK)
		return
	}

	var kendaraanTotalPerCabang uint = 0

	if leasing != "" {

		k, err := repository.GetKendaraanPerCabangTotal(c, leasing, cabang)

		if err != nil {
			payloads.HandleSuccess(c, nil, "Something went wrong", http.StatusOK)
			return
		}

		kendaraanTotalPerCabang = k
		logrus.Info(kendaraanTotalPerCabang)
	}

	leasingChart, err := repository.GetLeasingChart(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Something went wrong", http.StatusOK)
		return
	}

	leasingTotal, err := repository.GetLeasingTotal(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Something went wrong", http.StatusOK)
		return
	}

	userTotal, err := repository.GetUserTotalInfo(c)

	if err != nil {
		payloads.HandleSuccess(c, nil, "Somethig went wrong", http.StatusOK)
		return
	}

	data := make(map[string]interface{})
	data["leasing"] = leasingTotal
	data["kendaraan"] = kendaraanTotal
	data["kendaraan_per_cabang"] = kendaraanTotalPerCabang
	data["trial_members"] = userTotal.TrialMembers
	data["premium_members"] = userTotal.PremiumMembers
	data["expired_members"] = userTotal.ExpiredMembers
	data["leasing_chart"] = leasingChart

	payloads.HandleSuccess(c, data, "Data found", http.StatusOK)
}
