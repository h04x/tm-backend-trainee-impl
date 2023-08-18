package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
)

const GetStatSql = `select date::text, views, clicks, cost, 
round(coalesce(cost / NULLIF(clicks, 0), 0), 2) as cpc, 
round(coalesce(cost / NULLIF(views, 0) * 1000, 0), 2) as cpm
from clicks
where date between $1 and $2 
order by %s`

type GetStatReq struct {
	From  string `binding:"required,datetime=2006-01-02"`
	To    string `binding:"required,datetime=2006-01-02"`
	Order string `binding:"oneof=Date Views Clicks Cost Cpc Cpm"`
}

type GetStatRow struct {
	Date   string
	Views  uint
	Clicks uint
	Cost   string
	Cpc    float32
	Cpm    float32
}

func getStat(db *pgxpool.Pool) gin.HandlerFunc {
	f := func(c *gin.Context) {

		r := GetStatReq{Order: "Date"}
		if err := c.ShouldBindBodyWith(&r, binding.JSON); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		rows, _ := db.Query(c, fmt.Sprintf(GetStatSql, r.Order), r.From, r.To)
		defer rows.Close()

		stats := []GetStatRow{}
		for rows.Next() {
			var row GetStatRow
			if err := rows.Scan(&row.Date, &row.Views, &row.Clicks,
				&row.Cost, &row.Cpc, &row.Cpm); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			stats = append(stats, row)
		}

		if err := rows.Err(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, &stats)
	}
	return f
}
