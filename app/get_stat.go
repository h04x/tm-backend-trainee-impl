package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

const GetStatSql = `select date::text, views, clicks, cost, round(cost / clicks, 2) as cpc, 
round(cost / views * 1000, 2) as cpm from clicks
where date between $1 and $2 order by date`

type GetStatReq struct {
	From string `binding:"required,datetime=2006-01-02"`
	To   string `binding:"required,datetime=2006-01-02"`
}

type GetStatResp struct {
	Date   string
	Views  uint
	Clicks uint
	Cost   string
	Cpc    float32
	Cpm    float32
}

func getStat(db *pgxpool.Pool) gin.HandlerFunc {
	f := func(c *gin.Context) {

		r := GetStatReq{}
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		rows, err := db.Query(c, GetStatSql, r.From, r.To)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		var stats []GetStatResp
		for rows.Next() {
			var showStat GetStatResp
			if err := rows.Scan(&showStat.Date, &showStat.Views, &showStat.Clicks,
				&showStat.Cost, &showStat.Cpc, &showStat.Cpm); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			stats = append(stats, showStat)
		}

		c.JSON(http.StatusOK, &stats)
	}
	return f
}
