package controllers

import (
	"net/http"
	"regexp"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"xiuyanzhe/database"
	"xiuyanzhe/models"
	"xiuyanzhe/middleware"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "服务运行正常",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Register 注册
func Register(c *gin.Context) {
	var registerData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户名长度 (3-12个字符)
	if len(registerData.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名必须大于2个字符"})
		return
	}
	if len(registerData.Username) > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能超过12个字符"})
		return
	}

	// 验证邮箱
	if registerData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入邮箱地址"})
		return
	}
	// 简单的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(registerData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入有效的邮箱地址"})
		return
	}

	// 验证密码长度
	if len(registerData.Password) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码必须大于6个字符"})
		return
	}

	// 验证密码不能包含中文
	for _, r := range registerData.Password {
		if r >= '\u4e00' && r <= '\u9fff' {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能包含中文字符"})
			return
		}
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ?", registerData.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	result = database.DB.Where("email = ?", registerData.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username:  registerData.Username,
		Email:     registerData.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		LastLogin: time.Now(),
	}

	result = database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 登录
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.DB.Where("username = ?", loginData.Username).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	database.DB.Save(&user)

	// 记录访问日志（1小时内不重复）
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	var recentVisit models.VisitLog
	database.DB.Where("user_id = ? AND visit_time > ?", user.ID, oneHourAgo).First(&recentVisit)
	if recentVisit.ID == 0 {
		clientIP := c.ClientIP()
		visitLog := models.VisitLog{
			UserID:    user.ID,
			IP:        clientIP,
			VisitTime: time.Now(),
			Page:      "login",
		}
		database.DB.Create(&visitLog)
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
