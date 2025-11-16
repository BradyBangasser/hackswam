package environment

import (
	"hackswam/m/src/db"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type EnvImagePayload struct {
	PM25 float64 `form:"pm25"`
	PM10 float64 `form:"pm10"`
	PH   float64 `form:"ph"`
	DB   float64 `form:"db"`
	Lat  float64 `form:"lat"`
	Lon  float64 `form:"lon"`
}

func POST(c *gin.Context) {
	var data EnvImagePayload
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read uploaded image file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing image"})
		return
	}
	defer file.Close()

	// Ensure upload directory exists
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Create unique filename
	filename := time.Now().Format("20060102_150405") + "_" + header.Filename
	path := filepath.Join(uploadDir, filename)

	// Save file to disk
	out, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert metadata + geolocation into Postgres (PostGIS required)
	_, err = db.DB.Exec(`
        INSERT INTO env_samples (pm25, pm10, ph, db, image_path, geom)
        VALUES ($1, $2, $3, $4, $5,
            ST_SetSRID(ST_MakePoint($6, $7), 4326))
    `, data.PM25, data.PM10, data.PH, data.DB, path, data.Lon, data.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "environment + image uploaded"})
}
