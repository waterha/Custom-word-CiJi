package controllers

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// SearchWords 搜索单词（支持部分匹配，同时搜索单词和中文释义）
func SearchWords(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	likePattern := "%" + query + "%"

	// 搜索系统词库（同时匹配单词和中文释义）
	var systemWords []models.Word
	database.DB.Where("word LIKE ? OR translation LIKE ?", likePattern, likePattern).Find(&systemWords)

	// 搜索用户自定义词库
	var customWords []models.CustomWord
	database.DB.Where("user_id = ? AND (word LIKE ? OR translation LIKE ?)", uid, likePattern, likePattern).Find(&customWords)

	// 预分配容量，组合结果
	results := make([]gin.H, 0, len(systemWords)+len(customWords))

	for _, w := range systemWords {
		results = append(results, gin.H{
			"id":                         w.ID,
			"word":                       w.Word,
			"translation":                w.Translation,
			"example_sentence":           w.ExampleSentence,
			"example_sentence_translation": w.ExampleSentenceTranslation,
			"level":                      w.Level,
			"source":                     "system",
		})
	}

	for _, w := range customWords {
		results = append(results, gin.H{
			"id":                         w.ID,
			"word":                       w.Word,
			"translation":                w.Translation,
			"example_sentence":           w.ExampleSentence,
			"example_sentence_translation": w.ExampleSentenceTranslation,
			"source":                     "custom",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"total":   len(results),
	})
}
