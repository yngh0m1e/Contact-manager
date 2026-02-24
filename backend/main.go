package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var db *gorm.DB

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string `json:"-"`
	Email    string
	Name     string
}

type Contact struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `json:"user_id"`
	Name   string `gorm:"index" json:"name"`
	Email  string `gorm:"index" json:"email"`
	Phone  string `gorm:"index" json:"phone"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func main() {
	initDB()

	r := gin.Default()

	// ==================== CORS (ОБЯЗАТЕЛЬНО ПЕРВЫМ!) ====================
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		fmt.Println("CORS middleware сработал для:", c.Request.Method, c.Request.URL.Path) // ← для отладки

		if c.Request.Method == "OPTIONS" {
			fmt.Println("→ OPTIONS запрос обработан, возвращаем 204")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ====================== РОУТЫ ======================
	r.POST("/register", register)
	r.POST("/login", login)

	protected := r.Group("/")
	protected.Use(authMiddleware())
	protected.GET("/profile", getProfile)
	protected.PUT("/profile", updateProfile)
	protected.GET("/contacts", getContacts)
	protected.POST("/contacts", createContact)
	protected.GET("/contacts/:id", getContact)
	protected.PUT("/contacts/:id", updateContact)
	protected.DELETE("/contacts/:id", deleteContact)

	r.Run(":8080")
}

func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Contact{})

	// тестовый пользователь (логин: test, пароль: password)
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&User{Username: "test", Password: string(hashed), Email: "test@example.com", Name: "Тестовый Пользователь"})
	}
}

func generateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID:           userID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется токен"})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}
		c.Set("userID", token.Claims.(*Claims).UserID)
		c.Next()
	}
}

func getUserID(c *gin.Context) uint { return c.GetUint("userID") }

// ====================== AUTH ======================
func register(c *gin.Context) {
	type Req struct{ Username, Password, Email, Name string }
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := User{Username: req.Username, Password: string(hashed), Email: req.Email, Name: req.Name}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь уже существует"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешна"})
}

func login(c *gin.Context) {
	type Req struct{ Username, Password string }
	var req Req
	c.ShouldBindJSON(&req)
	var user User
	if db.Where("username = ?", req.Username).First(&user).Error != nil ||
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}
	token, _ := generateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ====================== PROFILE ======================
func getProfile(c *gin.Context) {
	var u User
	db.First(&u, getUserID(c))
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "username": u.Username, "email": u.Email, "name": u.Name})
}

func updateProfile(c *gin.Context) {
	var u User
	db.First(&u, getUserID(c))
	type Req struct{ Email, Name string }
	var req Req
	c.ShouldBindJSON(&req)
	u.Email = req.Email
	u.Name = req.Name
	db.Save(&u)
	c.JSON(http.StatusOK, gin.H{"message": "Профиль обновлён"})
}

// ====================== CONTACTS ======================
func getContacts(c *gin.Context) {
	userID := getUserID(c)
	search := c.Query("search")
	var contacts []Contact
	query := db.Where("user_id = ?", userID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}
	query.Find(&contacts)
	c.JSON(http.StatusOK, contacts)
}

func createContact(c *gin.Context) {
	userID := getUserID(c)
	var contact Contact
	c.ShouldBindJSON(&contact)
	contact.UserID = userID
	db.Create(&contact)
	c.JSON(http.StatusCreated, contact)
}

func getContact(c *gin.Context) {
	var contact Contact
	db.First(&contact, c.Param("id"))
	c.JSON(http.StatusOK, contact)
}

func updateContact(c *gin.Context) {
	var contact Contact
	db.First(&contact, c.Param("id"))
	c.ShouldBindJSON(&contact)
	db.Save(&contact)
	c.JSON(http.StatusOK, contact)
}

func deleteContact(c *gin.Context) {
	db.Delete(&Contact{}, c.Param("id"))
	c.JSON(http.StatusNoContent, nil)
}
