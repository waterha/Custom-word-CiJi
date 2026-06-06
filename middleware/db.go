package middleware

import (
	"log"
	"time"

	"xiuyanzhe/database"
)

// QueryTimeout 数据库查询超时时间
const QueryTimeout = 5 * time.Second

// EnsureIndexes 确保关键数据库索引存在
// MySQL 不支持 CREATE INDEX IF NOT EXISTS，所以先查 information_schema
func EnsureIndexes() {
	type idxResult struct {
		Count int64
	}

	indexes := []struct {
		table string
		name  string
		sql   string
	}{
		{"user_word_progresses", "idx_uwp_user_status", "CREATE INDEX idx_uwp_user_status ON user_word_progresses(user_id, status)"},
		{"user_word_progresses", "idx_uwp_user_reviewed", "CREATE INDEX idx_uwp_user_reviewed ON user_word_progresses(user_id, last_reviewed)"},
		{"user_word_progresses", "idx_uwp_word_user", "CREATE INDEX idx_uwp_word_user ON user_word_progresses(word_id, user_id)"},
		{"wrong_words", "idx_ww_user_count", "CREATE INDEX idx_ww_user_count ON wrong_words(user_id, wrong_count DESC)"},
		{"visit_logs", "idx_vl_user_time", "CREATE INDEX idx_vl_user_time ON visit_logs(user_id, visit_time)"},
		{"visit_logs", "idx_vl_visit_time", "CREATE INDEX idx_vl_visit_time ON visit_logs(visit_time)"},
		{"custom_words", "idx_cw_user", "CREATE INDEX idx_cw_user ON custom_words(user_id)"},
		{"words", "idx_words_level", "CREATE INDEX idx_words_level ON words(level)"},
	}

	for _, idx := range indexes {
		var result idxResult
		database.DB.Raw(
			"SELECT COUNT(*) as count FROM information_schema.STATISTICS WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name,
		).Scan(&result)

		if result.Count == 0 {
			if err := database.DB.Exec(idx.sql).Error; err != nil {
				log.Printf("[索引警告] 创建索引 %s 失败: %v", idx.name, err)
			} else {
				log.Printf("[索引] 已创建索引 %s ON %s", idx.name, idx.table)
			}
		}
	}
	log.Println("[索引] 数据库索引检查完成")
}
