package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// GetCustomWords 获取当前用户的自定义单词列表
func GetCustomWords(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var customWords []models.CustomWord
	database.DB.Where("user_id = ?", uid).Find(&customWords)

	c.JSON(http.StatusOK, customWords)
}

// GetCustomWord 获取单个自定义单词详情
func GetCustomWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id := c.Param("id")

	var customWord models.CustomWord
	result := database.DB.Where("id = ? AND user_id = ?", id, uid).First(&customWord)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单词不存在"})
		return
	}

	c.JSON(http.StatusOK, customWord)
}

// AddCustomWord 添加自定义单词
func AddCustomWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var customWord models.CustomWord
	if err := c.ShouldBindJSON(&customWord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customWord.UserID = uid

	result := database.DB.Create(&customWord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// UpdateCustomWord 更新自定义单词
func UpdateCustomWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id := c.Param("id")

	var customWord models.CustomWord
	if err := c.ShouldBindJSON(&customWord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingWord models.CustomWord
	result := database.DB.Where("id = ? AND user_id = ?", id, uid).First(&existingWord)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单词不存在"})
		return
	}

	customWord.ID = existingWord.ID
	customWord.UserID = uid

	result = database.DB.Save(&customWord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteCustomWord 删除自定义单词
func DeleteCustomWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id := c.Param("id")

	var customWord models.CustomWord
	result := database.DB.Where("id = ? AND user_id = ?", id, uid).First(&customWord)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单词不存在"})
		return
	}

	result = database.DB.Delete(&customWord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}