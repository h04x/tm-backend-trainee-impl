package main

import (
	//"database/sql"
	"fmt"
	"gin-helloworld/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	//_ "github.com/lib/pq"
	"github.com/jackc/pgx/v5"
	"regexp"
)

func main() {

	config, err := config.NewConfig("config/config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", config.Pg.Url)
	if err != nil {
		panic(err)
	}

	var valid_currency = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)

	var TwoDigitAfterPointNumber validator.Func = func(fl validator.FieldLevel) bool {
		return valid_currency.MatchString(fl.Field().String())
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("TwoDigitAfterPointNumber", TwoDigitAfterPointNumber)
	}

	router := gin.Default()

	router.POST("/save_stat", saveStat(db))
	router.POST("/get_stat", getStat(db))

	router.Run(config.App.Listen)
}
