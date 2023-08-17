package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
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

// saveStat
// @Summary      List accounts
// @Description  save stat
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        date    query     string  true  "date"  Format(string)
// @Success      200  {array}   SaveStat
// @Router       /save_stat [get]
func saveStat(db *pgxpool.Pool) gin.HandlerFunc {
	f := func(c *gin.Context) {
		// init special default Cost field
		t := SaveStat{
			Cost: "0",
		}
		if err := c.ShouldBindBodyWith(&t, binding.JSON); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		_, err := db.Exec(c, SaveStatdbReq, t.Date, t.Views, t.Clicks, t.Cost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusAccepted)
	}
	return f
}
