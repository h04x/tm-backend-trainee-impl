package main

import (
	"bytes"
	"gin-helloworld/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestHelper struct {
	t      *testing.T
	router *gin.Engine
}

// Retrun body and http code
func (rq *RequestHelper) DoRequest(method string, url string, body string) ([]byte, int) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	assert.Equal(rq.t, err, nil)
	rr := httptest.NewRecorder()
	rq.router.ServeHTTP(rr, req)
	responseData, err := io.ReadAll(rr.Body)
	assert.Equal(rq.t, err, nil)
	return responseData, rr.Code
}

func TestSaveStat(t *testing.T) {
	app, err := app.Testing()
	assert.Equal(t, err, nil)
	defer app.Db.Close()

	rh := RequestHelper{
		t,
		app.Router,
	}

	input := []struct {
		requested_json string
		expected_code  int
	}{
		{`{ "date": "2022-12-31", "views": 9, "cost": "2.00", "clicks": 1  }`, 202},
		{`{ "date": "2022-12-31" }`, 202},
		{`{ "date": "2022-12-31", "cost": "2.00" }`, 202},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 202},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.0" }`, 202},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2" }`, 202},
		{``, 400},
		{`{ "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "202a-12-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-13-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-32", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": "9", "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": -9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": "1", "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": -1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": 2.00 }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "-2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.001" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.a0" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2." }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2." }`, 400},
	}

	for _, row := range input {
		body, actual_code := rh.DoRequest("POST", "/save_stat", row.requested_json)
		b := assert.Equal(t, []byte{}, body)
		c := assert.Equal(t, row.expected_code, actual_code)
		if !(b && c) {
			t.Errorf("Failed at: %s", row.requested_json)
		}

	}
}
