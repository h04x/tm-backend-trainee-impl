package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

const SaveStatdbReq = `insert into clicks (date, views, clicks, cost) 
values ($1, $2, $3, $4)
on conflict (date) 
do update 
set views = clicks.views + $2, clicks = clicks.clicks + $3, cost = clicks.cost + $4`

type SaveStat struct {
	Date   string `binding:"required,datetime=2006-01-02"`
	Views  uint
	Clicks uint
	Cost   string `binding:"TwoDigitAfterPointNumber"`
}

func saveStat(db *sql.DB) gin.HandlerFunc {
	f := func(c *gin.Context) {
		t := SaveStat{}
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		_, err := db.Query(SaveStatdbReq, t.Date, t.Views, t.Clicks, t.Cost)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusAccepted, &t)
	}
	return f
}
