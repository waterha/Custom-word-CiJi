package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
	"xiuyanzhe/models"
)

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var DB *gorm.DB

func InitDB() {
	// 从环境变量获取数据库连接信息
	host := getEnv("MYSQL_HOST", "mysql") // 本地测试改成localhost
	port := getEnv("MYSQL_PORT", "3306")
	user := getEnv("MYSQL_USER", "root")
	password := getEnv("MYSQL_PASSWORD", "123456")
	dbname := getEnv("MYSQL_DATABASE", "xiuyanzhe")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// 配置数据库连接池
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic("获取数据库连接池失败")
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Word{},
		&models.UserWordProgress{},
		&models.WrongWord{},
		&models.VisitLog{},
		&models.CustomWord{},
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 添加全文索引以优化搜索性能
	var idxCount int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.STATISTICS WHERE table_schema = DATABASE() AND table_name = 'words' AND index_name = 'idx_words_word'").Scan(&idxCount)
	if idxCount == 0 {
		DB.Exec("CREATE FULLTEXT INDEX idx_words_word ON words(word)")
	}
	DB.Raw("SELECT COUNT(*) FROM information_schema.STATISTICS WHERE table_schema = DATABASE() AND table_name = 'custom_words' AND index_name = 'idx_custom_words_word'").Scan(&idxCount)
	if idxCount == 0 {
		DB.Exec("CREATE FULLTEXT INDEX idx_custom_words_word ON custom_words(word)")
	}

	// 初始化管理员账号
	initAdmin()

	// 初始化示例单词
	initWords()
}

func initAdmin() {
	var admin models.User
	result := DB.Where("username = ?", "admin").First(&admin)
	if result.RowsAffected == 0 {
		admin = models.User{
			Username:  "admin",
			Password:  "$2b$12$mX8oWvNnba4O6yDQu7NuG.yS.LmLuxtJyCjyfyQ31WKfL0dhqqJb.", // admin
			Role:      "admin",
			LastLogin: time.Now(),
		}
		DB.Create(&admin)
	} else {
		DB.Model(&admin).Update("password", "$2b$12$mX8oWvNnba4O6yDQu7NuG.yS.LmLuxtJyCjyfyQ31WKfL0dhqqJb.") // admin
	}
}

func initWords() {
	var count int64
	DB.Model(&models.Word{}).Count(&count)
	if count == 0 {
		words := []models.Word{
			{Word: "apple", Translation: "苹果", ExampleSentence: "I have an <span class=\"highlight\">apple</span>", Level: "all"},
			{Word: "banana", Translation: "香蕉", ExampleSentence: "I like <span class=\"highlight\">banana</span>", Level: "all"},
			{Word: "cat", Translation: "猫", ExampleSentence: "The <span class=\"highlight\">cat</span> is cute", Level: "all"},
			{Word: "dog", Translation: "狗", ExampleSentence: "The <span class=\"highlight\">dog</span> is loyal", Level: "all"},
			{Word: "egg", Translation: "鸡蛋", ExampleSentence: "I eat an <span class=\"highlight\">egg</span> every morning", Level: "all"},
			{Word: "fish", Translation: "鱼", ExampleSentence: "The <span class=\"highlight\">fish</span> swims in the river", Level: "all"},
			{Word: "goat", Translation: "山羊", ExampleSentence: "The <span class=\"highlight\">goat</span> eats grass", Level: "all"},
			{Word: "horse", Translation: "马", ExampleSentence: "The <span class=\"highlight\">horse</span> runs fast", Level: "all"},
			{Word: "ice", Translation: "冰", ExampleSentence: "The <span class=\"highlight\">ice</span> is cold", Level: "all"},
			{Word: "juice", Translation: "果汁", ExampleSentence: "I drink <span class=\"highlight\">juice</span> every day", Level: "all"},
		}
		DB.Create(&words)
	}
}
