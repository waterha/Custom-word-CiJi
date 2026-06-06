package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// GetWords 获取单词列表（管理员）- 支持分页和关键字搜索
func GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	db := database.DB.Model(&models.Word{})

	// 关键字过滤（模糊搜索单词或释义）
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("word LIKE ? OR translation LIKE ?", like, like)
	}

	// 级别过滤
	level := c.Query("level")
	if level == "cet4" || level == "all" {
		db = db.Where("level = ?", level)
	}

	db.Count(&total)

	var words []models.Word
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Find(&words)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取单词列表失败"})
		return
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        words,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// AddWord 新增单词
func AddWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&word)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// UpdateWord 更新单词
func UpdateWord(c *gin.Context) {
	id := c.Param("id")
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Model(&models.Word{}).Where("id = ?", id).Updates(word)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteWord 删除单词
func DeleteWord(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Word{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
