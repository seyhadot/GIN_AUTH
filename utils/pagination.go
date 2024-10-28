package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	return
}

func SendPaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	c.JSON(200, gin.H{
		"data": data,
		"pagination": gin.H{
			"total":       total,
			"total_pages": totalPages,
			"page":        page,
			"limit":       limit,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	})
}
