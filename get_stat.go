package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const GetStatReq = `select date::text, views, clicks, cost, round(cost / clicks, 2) as cpc, 
round(cost / views * 1000, 2) as cpm from clicks
where date between $1 and $2 order by date`

type ShowStatReq struct {
	From string `binding:"required,datetime=2006-01-02"`
	To   string `binding:"required,datetime=2006-01-02"`
}

type ShowStatResp struct {
	Date   string
	Views  uint
	Clicks uint
	Cost   string
	Cpc    float32
	Cpm    float32
}

func getStat(db *sql.DB) gin.HandlerFunc {
	f := func(c *gin.Context) {
		r := ShowStatReq{}

		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		rows, err := db.Query(GetStatReq, r.From, r.To)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var stats []ShowStatResp
		for rows.Next() {
			var showStat ShowStatResp
			if err := rows.Scan(&showStat.Date, &showStat.Views, &showStat.Clicks,
				&showStat.Cost, &showStat.Cpc, &showStat.Cpm); err != nil {
				log.Fatal(err)
			}
			stats = append(stats, showStat)
		}

		c.JSON(http.StatusAccepted, &stats)
	}
	return f
}
