package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func clearStat(db *pgxpool.Pool) gin.HandlerFunc {
	f := func(c *gin.Context) {
		_, err := db.Exec(c, "truncate clicks")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusOK)
	}
	return f
}
