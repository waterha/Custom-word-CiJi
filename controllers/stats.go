package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// TodayStats 今日统计
type TodayStats struct {
	Total   int64 `json:"total"`
	Known   int64 `json:"known"`
	Unknown int64 `json:"unknown"`
}

// OverallStats 累计统计
type OverallStats struct {
	TotalWords    int64   `json:"total_words"`
	LearnedWords  int64   `json:"learned_words"`
	KnownWords    int64   `json:"known_words"`
	UnknownWords  int64   `json:"unknown_words"`
	WrongWords    int64   `json:"wrong_words"`
	Progress      float64 `json:"progress"`
	Accuracy      float64 `json:"accuracy"`
}

// DailyCount 每日学习数量
type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// WrongWordItem 常错单词
type WrongWordItem struct {
	Word       string `json:"word"`
	Translation string `json:"translation"`
	WrongCount int    `json:"wrong_count"`
}

// LearningStatsResponse 学习状况响应
type LearningStatsResponse struct {
	Today      TodayStats      `json:"today"`
	Overall    OverallStats    `json:"overall"`
	Weekly     []DailyCount    `json:"weekly"`
	TopWrong   []WrongWordItem `json:"top_wrong"`
	CustomWords int64          `json:"custom_words"`
	StreakDays int             `json:"streak_days"`
}

// GetLearningStats 获取学习状况统计
func GetLearningStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	today := time.Now().Truncate(24 * time.Hour)

	var resp LearningStatsResponse

	// 1. 今日统计
	database.DB.Raw(`
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'known' THEN 1 ELSE 0 END), 0) as known,
			COALESCE(SUM(CASE WHEN status = 'unknown' THEN 1 ELSE 0 END), 0) as unknown
		FROM user_word_progresses
		WHERE user_id = ? AND last_reviewed >= ?
	`, uid, today).Scan(&resp.Today)

	// 2. 累计统计
	var totalWords int64
	database.DB.Model(&models.Word{}).Count(&totalWords)

	var learnedWords int64
	database.DB.Model(&models.UserWordProgress{}).
		Where("user_id = ?", uid).
		Distinct("word_id").Count(&learnedWords)

	var knownWords, unknownWords int64
	database.DB.Model(&models.UserWordProgress{}).
		Where("user_id = ? AND status = 'known'", uid).
		Distinct("word_id").Count(&knownWords)
	database.DB.Model(&models.UserWordProgress{}).
		Where("user_id = ? AND status = 'unknown'", uid).
		Distinct("word_id").Count(&unknownWords)

	var wrongWords int64
	database.DB.Model(&models.WrongWord{}).
		Where("user_id = ?", uid).Count(&wrongWords)

	resp.Overall = OverallStats{
		TotalWords:   totalWords,
		LearnedWords: learnedWords,
		KnownWords:   knownWords,
		UnknownWords: unknownWords,
		WrongWords:   wrongWords,
		Progress:     float64(learnedWords) / float64(max(totalWords, 1)),
		Accuracy:     float64(knownWords) / float64(max(knownWords+unknownWords, 1)),
	}

	// 3. 最近7天趋势
	sevenDaysAgo := today.AddDate(0, 0, -6)
	database.DB.Raw(`
		SELECT DATE(last_reviewed) as date, COUNT(*) as count
		FROM user_word_progresses
		WHERE user_id = ? AND last_reviewed >= ?
		GROUP BY DATE(last_reviewed)
		ORDER BY date
	`, uid, sevenDaysAgo).Scan(&resp.Weekly)

	// 补全没有学习记录的天数为0
	dateMap := make(map[string]int64)
	for _, d := range resp.Weekly {
		dateMap[d.Date] = d.Count
	}
	resp.Weekly = nil
	for i := 6; i >= 0; i-- {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		resp.Weekly = append(resp.Weekly, DailyCount{
			Date:  day,
			Count: dateMap[day],
		})
	}

	// 4. 常错单词 Top 5
	database.DB.Raw(`
		SELECT w.word, w.translation, ww.wrong_count
		FROM wrong_words ww
		JOIN words w ON w.id = ww.word_id
		WHERE ww.user_id = ?
		ORDER BY ww.wrong_count DESC
		LIMIT 5
	`, uid).Scan(&resp.TopWrong)

	// 5. 自定义单词数
	database.DB.Model(&models.CustomWord{}).
		Where("user_id = ?", uid).Count(&resp.CustomWords)

	// 6. 连续学习天数
	resp.StreakDays = calcStreak(uid, today)

	c.JSON(http.StatusOK, resp)
}

func calcStreak(uid uint, today time.Time) int {
	type studyDay struct {
		Date string
	}
	var days []studyDay
	database.DB.Raw(`
		SELECT DISTINCT DATE(last_reviewed) as date
		FROM user_word_progresses
		WHERE user_id = ? AND last_reviewed >= DATE_SUB(?, INTERVAL 30 DAY)
		ORDER BY date DESC
	`, uid, today).Scan(&days)

	if len(days) == 0 {
		return 0
	}

	// 如果今天没有学习记录，连续天数为0
	todayStr := today.Format("2006-01-02")
	if days[0].Date != todayStr {
		return 0
	}

	streak := 0
	checkDate := today
	for _, d := range days {
		if d.Date == checkDate.Format("2006-01-02") {
			streak++
			checkDate = checkDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak
}
