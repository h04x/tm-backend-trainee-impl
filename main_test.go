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

func TestClearStat(t *testing.T) {
	app, err := app.Testing()
	assert.Equal(t, err, nil)
	defer app.Db.Close()

	rh := RequestHelper{
		t,
		app.Router,
	}

	// call clear
	body, actual_code := rh.DoRequest("DELETE", "/clear_stat", "")
	assert.Equal(t, []byte{}, body)
	assert.Equal(t, 200, actual_code)

	// add some row
	body, actual_code = rh.DoRequest("POST", "/save_stat",
		`{ "date": "1234-12-12" }`)
	assert.Equal(t, []byte{}, body)
	assert.Equal(t, 202, actual_code)

	// make sure that row is there
	body, actual_code = rh.DoRequest("POST", "/get_stat",
		`{ "from": "1234-12-12", "to": "1234-12-12" }`)
	assert.Equal(t, `[{"Date":"1234-12-12","Views":0,"Clicks":0,"Cost":"0","Cpc":0,"Cpm":0}]`, string(body))
	assert.Equal(t, 200, actual_code)

	// call clear
	body, actual_code = rh.DoRequest("DELETE", "/clear_stat", "")
	assert.Equal(t, []byte{}, body)
	assert.Equal(t, 200, actual_code)

	// expected empty array
	body, actual_code = rh.DoRequest("POST", "/get_stat",
		`{ "from": "1234-12-12", "to": "1234-12-12" }`)
	assert.Equal(t, `[]`, string(body))
	assert.Equal(t, 200, actual_code)
}

func TestGetStat(t *testing.T) {
	app, err := app.Testing()
	assert.Equal(t, err, nil)
	defer app.Db.Close()

	rh := RequestHelper{
		t,
		app.Router,
	}

	body, actual_code := rh.DoRequest("DELETE", "/clear_stat", "")
	assert.Equal(t, []byte{}, body)
	assert.Equal(t, 200, actual_code)

	input := []struct {
		requested_json string
		expected_code  int
	}{
		{`{ "date": "1111-11-11", "views": 15, "clicks": 0, "cost": "15" }`, 202},
		{`{ "date": "2222-12-12", "views": 2, "clicks": 4, "cost": "0" }`, 202},
		{`{ "date": "2222-12-12", "views": 1, "clicks": 1, "cost": "0.01" }`, 202},
		{`{ "date": "2222-12-13", "views": 9, "clicks": 9, "cost": "0" }`, 202},
		{`{ "date": "2222-12-12", "views": 1, "clicks": 1, "cost": "2" }`, 202},
	}

	for _, row := range input {
		body, actual_code := rh.DoRequest("POST", "/save_stat", row.requested_json)
		b := assert.Equal(t, []byte{}, body)
		c := assert.Equal(t, row.expected_code, actual_code)
		if !(b && c) {
			t.Errorf("Failed at: %s", row.requested_json)
		}
	}

	test_case_get_stat := []struct {
		requested_json string
		expected_code  int
		expected_resp  string
	}{
		{`{ "from": "0000-00-00", "to": "0000-00-00" }`, 400, ""},
		{`{ "from": "0001-01-01", "to": "0001-13-01" }`, 400, ""},
		{`{ "from": "0001-01-01", "to": "0001-12-1" }`, 400, ""},
		{`{ "from": "0001-01-01", "to": "0001-01-01" }`, 200, "[]"},
		{`{ "from": "0001-01-01", "to": "0001-02-01" }`, 200, "[]"},

		{`{ "from": "2222-12-12", "to": "2222-12-12" }`, 200, `[{"Date":"2222-12-12","Views":4,"Clicks":6,"Cost":"2.01","Cpc":0.34,"Cpm":502.5}]`},

		{`{ "from": "1111-11-11", "to": "2222-12-13" }`, 200,
			`[{"Date":"1111-11-11","Views":15,"Clicks":0,"Cost":"15","Cpc":0,"Cpm":1000},` +
				`{"Date":"2222-12-12","Views":4,"Clicks":6,"Cost":"2.01","Cpc":0.34,"Cpm":502.5},` +
				`{"Date":"2222-12-13","Views":9,"Clicks":9,"Cost":"0","Cpc":0,"Cpm":0}]`},

		{`{ "from": "1111-11-11", "to": "2222-12-13", "order": "Views" }`, 200,
			`[{"Date":"2222-12-12","Views":4,"Clicks":6,"Cost":"2.01","Cpc":0.34,"Cpm":502.5},` +
				`{"Date":"2222-12-13","Views":9,"Clicks":9,"Cost":"0","Cpc":0,"Cpm":0},` +
				`{"Date":"1111-11-11","Views":15,"Clicks":0,"Cost":"15","Cpc":0,"Cpm":1000}]`},
	}

	for _, row := range test_case_get_stat {
		body, actual_code := rh.DoRequest("POST", "/get_stat", row.requested_json)
		b := assert.Equal(t, row.expected_resp, string(body))
		c := assert.Equal(t, row.expected_code, actual_code)
		if !(b && c) {
			t.Errorf("Failed at: %s", row.requested_json)
		}
	}
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
		{`{ "date": "2021-12-31", "views": 9, "clicks": 1, "cost": "3.00" }`, 202},
		{`{ "date": "2022-11-30", "views": 9, "clicks": 1, "cost": "2.0" }`, 202},
		{`{ "date": "2022-12-30", "views": 9, "clicks": 1, "cost": "1" }`, 202},
		{`{ "date": "2023-12-30", "views": 9, "clicks": 1, "cost": "1.23" }`, 202},
		{`{ "date": "2013-12-30", "views": 9, "clicks": 1, "cost": "0.0" }`, 202},
		{`{ "date": "2013-12-30", "views": 9, "clicks": 1, "cost": "0.00" }`, 202},
		{``, 400},
		{`{ "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "202a-12-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-13-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-13-31", "views": 9, "clicks": 1, "cost": "2,00" }`, 400},
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
