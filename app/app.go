package app

import (
	"context"
	"fmt"
	"gin-helloworld/config"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router *gin.Engine
	config *config.Config
	Db     *pgxpool.Pool
}

func Default() (*App, error) {
	config, err := config.NewConfig("config/config.yaml")
	if err != nil {
		return nil, err
	}
	return New(config)
}

func New(config *config.Config) (*App, error) {
	r := strings.NewReplacer("{user}", config.Pg.User, "{password}",
		config.Pg.Password, "{host}", config.Pg.Host)

	dbpool, err := pgxpool.New(context.Background(), r.Replace(config.Pg.Url))
	if err != nil {
		return nil, err
	}

	var valid_currency = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	var TwoDigitAfterPointNumber validator.Func = func(fl validator.FieldLevel) bool {
		return valid_currency.MatchString(fl.Field().String())
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("error while obtaining validator engine")
	}

	err = v.RegisterValidation("TwoDigitAfterPointNumber", TwoDigitAfterPointNumber)
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(errorHandler)

	router.POST("/save_stat", saveStat(dbpool))
	router.POST("/get_stat", getStat(dbpool))
	router.DELETE("/clear_stat", clearStat(dbpool))

	c := &App{
		router,
		config,
		dbpool,
	}
	return c, nil
}

func errorHandler(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		log.Println(c.Request.URL, "Error Msg:", err)

		// Log request Body on error
		// To get request Body request handler must use ShouldBindBodyWith() instead Bind()
		// https://github.com/gin-gonic/gin/issues/1377
		raw, _ := c.Get(gin.BodyBytesKey)
		body, _ := raw.([]byte)

		// truncate too long request bodies
		l := 100
		if len(body) < 100 {
			l = len(body)
		}
		log.Println("Error Body:", string(body[:l]))
	}
}

func (app *App) Run() error {
	listen := ":" + app.config.App.Port
	fmt.Println("Run server on " + listen)
	return app.Router.Run(listen)
}
