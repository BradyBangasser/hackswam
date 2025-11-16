package environment

import (
	"hackswam/m/src/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon query parameters required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
		return
	}

	// 1 mile ≈ 1609.34 meters
	rows, err := db.DB.Query(`
        SELECT id, image_path, pm25, pm10, ph, db, ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lon
        FROM env_samples
        WHERE ST_DWithin(
            geom,
            ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
            1609.34
        )
    `, lon, lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Sample struct {
		ID        int     `json:"id"`
		ImagePath string  `json:"image_path"`
		PM25      float64 `json:"pm25"`
		PM10      float64 `json:"pm10"`
		PH        float64 `json:"ph"`
		DB        float64 `json:"db"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
	}

	var results []Sample
	for rows.Next() {
		var s Sample
		if err := rows.Scan(&s.ID, &s.ImagePath, &s.PM25, &s.PM10, &s.PH, &s.DB, &s.Lat, &s.Lon); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}
