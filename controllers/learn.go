package controllers

import (
	"net/http"
	"strconv"
	"time"
	"math/rand"
	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
	"xiuyanzhe/redis"
)

// GetWord 获取单个单词详情
func GetWord(c *gin.Context) {
	id := c.Param("id")
	var word models.Word
	result := database.DB.First(&word, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单词不存在"})
		return
	}
	c.JSON(http.StatusOK, word)
}

// GetNextWord 获取下一个学习单词
func GetNextWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	level := c.DefaultQuery("level", "all")

	// 检查Redis中是否有未完成的轮次
	roundKey := "round:" + strconv.Itoa(int(uid)) + ":" + level
	wordIDStr, err := redis.RPop(roundKey)

	var wordID uint
	if err == nil {
		// 从Redis中获取单词ID
		wordID64, _ := strconv.ParseUint(wordIDStr, 10, 32)
		wordID = uint(wordID64)
	} else {
		// 轮次已完成，重新生成
		var wordIDs []uint
		database.DB.Model(&models.Word{}).Where("level = ?", level).Pluck("id", &wordIDs)

		// 如果没有单词
		if len(wordIDs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "本轮学习已完成"})
			return
		}

		// 随机打乱
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(wordIDs), func(i, j int) {
			wordIDs[i], wordIDs[j] = wordIDs[j], wordIDs[i]
		})

		// 存入Redis
		for _, id := range wordIDs {
			redis.LPush(roundKey, strconv.Itoa(int(id)))
		}

		// 获取第一个单词
		wordIDStr, _ = redis.RPop(roundKey)
		wordID64, _ := strconv.ParseUint(wordIDStr, 10, 32)
		wordID = uint(wordID64)
	}

	// 获取单词信息
	var word models.Word
	database.DB.First(&word, wordID)

	// 返回完整的单词信息，包括中文和例句
	c.JSON(http.StatusOK, word)
}

// SubmitAnswer 提交认识/不认识结果
func SubmitAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var answerData struct {
		WordID uint `json:"word_id"`
		Known  bool `json:"known"`
	}

	if err := c.ShouldBindJSON(&answerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新学习进度
	var progress models.UserWordProgress
	result := database.DB.Where("user_id = ? AND word_id = ?", uid, answerData.WordID).First(&progress)

	status := "known"
	if !answerData.Known {
		status = "unknown"
		// 更新错词记录
		var wrongWord models.WrongWord
		result := database.DB.Where("user_id = ? AND word_id = ?", uid, answerData.WordID).First(&wrongWord)
		if result.RowsAffected == 0 {
			wrongWord = models.WrongWord{
				UserID:     uid,
				WordID:     answerData.WordID,
				WrongCount: 1,
			}
			database.DB.Create(&wrongWord)
		} else {
			wrongWord.WrongCount++
			database.DB.Save(&wrongWord)
		}
	}

	if result.RowsAffected == 0 {
		progress = models.UserWordProgress{
			UserID:       uid,
			WordID:       answerData.WordID,
			Status:       status,
			LastReviewed: time.Now(),
		}
		database.DB.Create(&progress)
	} else {
		progress.Status = status
		progress.LastReviewed = time.Now()
		database.DB.Save(&progress)
	}

	c.JSON(http.StatusOK, gin.H{"message": "提交成功"})
}

// GetProgress 获取当前轮进度
func GetProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	level := c.DefaultQuery("level", "all")

	roundKey := "round:" + strconv.Itoa(int(uid)) + ":" + level
	remaining, _ := redis.LLen(roundKey)

	// 计算总单词数
	var total int64
	database.DB.Model(&models.Word{}).Where("level = ?", level).Count(&total)

	completed := total - int64(remaining)
	if remaining == 0 {
		completed = total
	}
	if completed < 0 {
		completed = 0
	}

	var progress float64
	if total > 0 {
		progress = float64(completed) / float64(total)
	} else {
		progress = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": completed,
		"total":     total,
		"progress":  progress,
	})
}

// GetNextCustomWord 获取下一个自定义单词
func GetNextCustomWord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	roundKey := "round:custom:" + strconv.Itoa(int(uid))
	wordIDStr, err := redis.RPop(roundKey)

	var wordID uint
	if err == nil {
		wordID64, _ := strconv.ParseUint(wordIDStr, 10, 32)
		wordID = uint(wordID64)
	} else {
		var wordIDs []uint
		database.DB.Model(&models.CustomWord{}).Where("user_id = ?", uid).Pluck("id", &wordIDs)

		if len(wordIDs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "本轮学习已完成"})
			return
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(wordIDs), func(i, j int) {
			wordIDs[i], wordIDs[j] = wordIDs[j], wordIDs[i]
		})

		for _, id := range wordIDs {
			redis.LPush(roundKey, strconv.Itoa(int(id)))
		}

		wordIDStr, _ = redis.RPop(roundKey)
		wordID64, _ := strconv.ParseUint(wordIDStr, 10, 32)
		wordID = uint(wordID64)
	}

	var word models.CustomWord
	database.DB.First(&word, wordID)

	c.JSON(http.StatusOK, word)
}

// SubmitCustomAnswer 提交自定义单词学习结果
func SubmitCustomAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var answerData struct {
		WordID uint `json:"word_id"`
		Known  bool `json:"known"`
	}

	if err := c.ShouldBindJSON(&answerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "known"
	if !answerData.Known {
		status = "unknown"
	}

	var progress models.UserWordProgress
	result := database.DB.Where("user_id = ? AND word_id = ?", uid, answerData.WordID).First(&progress)

	if result.RowsAffected == 0 {
		progress = models.UserWordProgress{
			UserID:       uid,
			WordID:       answerData.WordID,
			Status:       status,
			LastReviewed: time.Now(),
		}
		database.DB.Create(&progress)
	} else {
		progress.Status = status
		progress.LastReviewed = time.Now()
		database.DB.Save(&progress)
	}

	c.JSON(http.StatusOK, gin.H{"message": "提交成功"})
}

// GetCustomProgress 获取自定义单词学习进度
func GetCustomProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	roundKey := "round:custom:" + strconv.Itoa(int(uid))
	remaining, _ := redis.LLen(roundKey)

	var total int64
	database.DB.Model(&models.CustomWord{}).Where("user_id = ?", uid).Count(&total)

	completed := total - int64(remaining)
	if remaining == 0 {
		completed = total
	}
	if completed < 0 {
		completed = 0
	}

	var progress float64
	if total > 0 {
		progress = float64(completed) / float64(total)
	} else {
		progress = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": completed,
		"total":     total,
		"progress":  progress,
	})
}
