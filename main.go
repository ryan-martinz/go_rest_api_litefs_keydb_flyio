package main

import (
	"context"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Record struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Data string `json:"data"`
}

var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

func main() {
	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	gin.SetMode(gin.ReleaseMode)

	initializeDatabase()
	initializeRedis()

	setupRoutes(r)

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}

func initializeDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("/litefs/db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database: " + err.Error())
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	db.AutoMigrate(&Record{})
}

func initializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "keydb-internal:6379",
		Password: "yourpassword",
		DB:       0,
	})
}

func setupRoutes(r *gin.Engine) {
	r.GET("/", indexHandler)
	r.POST("/record", createRecordHandler)
	r.GET("/record/:id", getRecordHandler)
	r.GET("/records", getAllRecordsHandler)
}

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func createRecordHandler(c *gin.Context) {
	var record Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Explicitly ignore any ID that might have been provided in the request
	record.ID = 0

	if err := db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func getRecordHandler(c *gin.Context) {
	id := c.Param("id")
	var record Record

	if err := db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func getAllRecordsHandler(c *gin.Context) {
	var records []Record

	if err := db.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving records"})
		return
	}

	c.JSON(http.StatusOK, records)
}
