package redis

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"xiuyanzhe/database"
	"xiuyanzhe/models"
)

// QueueItem 异步任务队列项
type QueueItem struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// VisitLogData 访问日志数据
type VisitLogData struct {
	UserID    uint      `json:"user_id"`
	IP        string    `json:"ip"`
	VisitTime time.Time `json:"visit_time"`
	Page      string    `json:"page"`
}

// WrongWordData 错词更新数据
type WrongWordData struct {
	UserID uint `json:"user_id"`
	WordID uint `json:"word_id"`
	Known  bool `json:"known"`
}

var (
	queueOnce      sync.Once
	flushTicker    *time.Ticker
	stopFlush      chan struct{}
	batchSize      = 50
	flushInterval  = 10 * time.Second
)

// StartAsyncQueue 启动异步队列处理器
// 将非关键写入（访问日志、错词更新）先推入Redis队列，
// 后台协程定期批量刷入数据库，减少数据库写入压力
func StartAsyncQueue() {
	queueOnce.Do(func() {
		stopFlush = make(chan struct{})
		flushTicker = time.NewTicker(flushInterval)

		go func() {
			for {
				select {
				case <-flushTicker.C:
					flushVisitLogs()
				case <-stopFlush:
					return
				}
			}
		}()

		log.Println("[AsyncQueue] 异步任务队列已启动，刷新间隔:", flushInterval)
	})
}

// StopAsyncQueue 停止异步队列处理器
func StopAsyncQueue() {
	if flushTicker != nil {
		flushTicker.Stop()
	}
	if stopFlush != nil {
		close(stopFlush)
	}
}

// EnqueueVisitLog 将访问日志加入异步队列
func EnqueueVisitLog(userID uint, ip, page string) {
	data := VisitLogData{
		UserID:    userID,
		IP:        ip,
		VisitTime: time.Now(),
		Page:      page,
	}
	item := QueueItem{Type: "visit_log", Data: data}
	bytes, err := json.Marshal(item)
	if err != nil {
		return
	}

	// 推入Redis列表，左进右出
	_ = Client.LPush(Ctx, "queue:visit_logs", string(bytes))
}

// EnqueueWrongWord 将错词更新加入异步队列
func EnqueueWrongWord(userID, wordID uint, known bool) {
	data := WrongWordData{
		UserID: userID,
		WordID: wordID,
		Known:  known,
	}
	item := QueueItem{Type: "wrong_word", Data: data}
	bytes, err := json.Marshal(item)
	if err != nil {
		return
	}

	_ = Client.LPush(Ctx, "queue:wrong_words", string(bytes))
}

// flushVisitLogs 批量刷新访问日志到数据库
func flushVisitLogs() {
	logs := make([]models.VisitLog, 0, batchSize)

	for i := 0; i < batchSize; i++ {
		result, err := Client.RPop(Ctx, "queue:visit_logs").Result()
		if err != nil {
			break
		}

		var item QueueItem
		if err := json.Unmarshal([]byte(result), &item); err != nil {
			continue
		}

		if item.Type == "visit_log" {
			dataBytes, _ := json.Marshal(item.Data)
			var data VisitLogData
			if err := json.Unmarshal(dataBytes, &data); err != nil {
				continue
			}

			logs = append(logs, models.VisitLog{
				UserID:    data.UserID,
				IP:        data.IP,
				VisitTime: data.VisitTime,
				Page:      data.Page,
			})
		}
	}

	if len(logs) > 0 {
		// 批量插入
		if err := database.DB.CreateInBatches(logs, 100).Error; err != nil {
			log.Printf("[AsyncQueue] 批量写入访问日志失败: %v", err)
			// 失败时重新入队
			for _, log := range logs {
				EnqueueVisitLog(log.UserID, log.IP, log.Page)
			}
		}
	}
}
