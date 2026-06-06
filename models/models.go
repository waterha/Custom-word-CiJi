package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100" json:"email"`
	Password  string         `gorm:"size:255" json:"-"`
	Role      string         `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	LastLogin time.Time      `json:"last_login"`
}

// Word 单词模型
type Word struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
	Word                     string         `gorm:"uniqueIndex;size:100" json:"word"`
	Translation              string         `gorm:"size:255" json:"translation"`
	ExampleSentence          string         `gorm:"type:text" json:"example_sentence"`
	ExampleSentenceTranslation string       `gorm:"type:text" json:"example_sentence_translation"`
	Level                    string         `gorm:"type:enum('all','cet4');default:'all'" json:"level"`
}

// UserWordProgress 学习进度模型
type UserWordProgress struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        uint           `json:"user_id"`
	WordID        uint           `json:"word_id"`
	Status        string         `gorm:"type:enum('known','unknown','unlearned');default:'unlearned'" json:"status"`
	LastReviewed  time.Time      `json:"last_reviewed"`
}

// WrongWord 错词模型
type WrongWord struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `json:"user_id"`
	WordID     uint           `json:"word_id"`
	WrongCount int            `gorm:"default:0" json:"wrong_count"`
}

// VisitLog 访问日志模型
type VisitLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `json:"user_id"`
	IP        string         `gorm:"size:45" json:"ip"`
	VisitTime time.Time      `json:"visit_time"`
	Page      string         `gorm:"size:100" json:"page"`
}

// CustomWord 用户自定义单词模型
type CustomWord struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
	UserID                   uint           `json:"user_id"`
	Word                     string         `gorm:"size:100" json:"word"`
	Translation              string         `gorm:"size:255" json:"translation"`
	ExampleSentence          string         `gorm:"type:text" json:"example_sentence"`
	ExampleSentenceTranslation string       `gorm:"type:text" json:"example_sentence_translation"`
}
