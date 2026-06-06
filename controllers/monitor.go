package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// GetMonitorOverview 获取监控数据
func GetMonitorOverview(c *gin.Context) {
	// 今日访问量（只统计登录相关的访问，即page为login的记录）
	var todayVisits int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.VisitLog{}).Where("DATE(visit_time) = ? AND page = ?", today, "login").Count(&todayVisits)

	// 错词排行（按所有用户的错词总数统计）
	type WrongWordRank struct {
		Word        string `json:"word"`
		Translation string `json:"translation"`
		WrongCount  int    `json:"wrong_count"`
	}

	var wrongWords []WrongWordRank
	database.DB.Table("wrong_words").
		Select("words.word, words.translation, SUM(wrong_words.wrong_count) as wrong_count").
		Joins("JOIN words ON wrong_words.word_id = words.id").
		Group("words.id").
		Order("wrong_count DESC").
		Limit(10).
		Scan(&wrongWords)

	c.JSON(http.StatusOK, gin.H{
		"today_visits": todayVisits,
		"wrong_words":  wrongWords,
	})
}

// GetHourlyVisits 获取24小时各时段访问量（只统计登录相关）
func GetHourlyVisits(c *gin.Context) {
	hourlyVisits := make([]int, 24)
	now := time.Now()
	today := now.Format("2006-01-02")

	for hour := 0; hour < 24; hour++ {
		var count int64
		database.DB.Model(&models.VisitLog{}).
			Where("DATE(visit_time) = ? AND HOUR(visit_time) = ? AND page = ?", today, hour, "login").
			Count(&count)
		hourlyVisits[hour] = int(count)
	}

	c.JSON(http.StatusOK, hourlyVisits)
}

// GetDailyRegistrations 获取最近7天注册量
func GetDailyRegistrations(c *gin.Context) {
	dailyRegistrations := make([]int, 7)
	now := time.Now()

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.User{}).
			Where("DATE(created_at) = ?", date).
			Count(&count)
		dailyRegistrations[6-i] = int(count)
	}

	c.JSON(http.StatusOK, dailyRegistrations)
}
