package controllers

import (
	"fmt"
	config "matel/configs"
	"matel/exceptions"
	"matel/helper"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func GetProfile(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	newUser, err := repository.UserProfile(c, UserID)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	var user models.User

	user.StartSubscription = newUser.StartSubscription
	user.EndSubscription = newUser.EndSubscription
	user.CreatedAt = newUser.CreatedAt

	newUser.Status = uint(helper.GetUserStatus(user))
	user.Status = newUser.Status

	if user.Status == 0 {
		var endDate = newUser.CreatedAt.Add(1 * 24 * time.Hour)
		user.EndSubscription = endDate.Format("2006-01-02")
		newUser.EndSubscription = user.EndSubscription
	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
}

func GetMember(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)
	search := c.Query("search")

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	var newUser []models.User

	user, err := repository.GetMember(c, search)

	for _, v := range user {
		v.Status = uint(helper.GetUserStatus(v))
		newUser = append(newUser, v)
	}

	for _, v := range newUser {
		logrus.Info(v.SubscriptionMonth)
	}

	if len(user) == 0 {
		payloads.HandleSuccess(c, nil, "User tidak ditemukan", http.StatusOK)
		return
	}

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
}

func MemberChange(c *gin.Context) {

	user, err := repository.MemberChange(c)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = exportMemberToExcel(user, c)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, nil, "Excel file generated and sent for download", 200)

}

func exportMemberToExcel(data []models.UserChangeExport, c *gin.Context) error {
	file := excelize.NewFile()
	sheetName := "Sheet1"

	// Set header row
	headers := []string{"Nama Pengguna", "Sebelum Diubah", "Setelah Diubah", "Mulai Berlangganan"}
	for col, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+col, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	// Set data rows
	for row, user := range data {
		cell := fmt.Sprintf("A%d", row+2)
		file.SetCellValue(sheetName, cell, user.UserName)
		cell = fmt.Sprintf("B%d", row+2)
		formattedUnupdatedStatus := fmt.Sprintf("%d hari", user.UnupdatedStatus)
		file.SetCellValue(sheetName, cell, formattedUnupdatedStatus)
		cell = fmt.Sprintf("C%d", row+2)
		formattedUpdatedStatus := fmt.Sprintf("%d hari", user.UpdatedStatus)
		file.SetCellValue(sheetName, cell, formattedUpdatedStatus)
		cell = fmt.Sprintf("D%d", row+2)
		file.SetCellValue(sheetName, cell, user.TimeUpdated)
	}

	// Set the content type and headers for the response
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=cabang_export.xlsx")

	// Save the Excel file to the response writer
	err := file.Write(c.Writer)
	if err != nil {
		return err
	}

	return nil
}

func SetUser(c *gin.Context) {
	type SetUserReq struct {
		UserID            uint   `json:"user_id" validate:"required"`
		SubscriptionMonth string `json:"subscription_month" validate:"required"`
	}
	var req SetUserReq
	c.BindJSON(&req)

	u, err := repository.UserProfile(c, req.UserID)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	sub, e := strconv.Atoi(req.SubscriptionMonth)

	if e != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = repository.SetUser(c, req.UserID, uint(sub))

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	var user models.UserChange

	currentTime := time.Now()

	// Define the desired format
	format := "02 Jan 2006"

	// Format the date and time to the desired format
	formattedDateTime := currentTime.Format(format)

	user.UserID = req.UserID
	user.TimeUpdated = formattedDateTime
	user.UnupdatedStatus = u.SubscriptionMonth
	user.UpdatedStatus = uint(sub)

	err = repository.UserChange(c, user)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}

func DeleteMember(c *gin.Context) {
	id := c.Param("id")

	user := models.User{}
	if err := config.InitDB().Where("id = ?", id).First(&user).Error; err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err := config.InitDB().Unscoped().Delete(&user).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}
