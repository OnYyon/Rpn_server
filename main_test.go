package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OnYyon/Rpn_server.git/server"
)

type Tests struct {
	exp    string
	output float64
}

func TestFormatexpression(t *testing.T) {
	var temp server.ResultJson
	testCases := []Tests{
		{exp: "2+2*2", output: 6},
		{exp: "(2+2)*2", output: 8},
		{exp: "((7+1)/(2+2)*4)/8*(32-((4+12)*2))-1", output: -1},
		{exp: "2 + 2 * 2", output: 6},
		{exp: "10/2+5", output: 10},
		{exp: "10/2-5", output: 0},
		{exp: "10*2+5", output: 25},
		{exp: "10+2*5", output: 20},
	}
	for _, tt := range testCases {
		reqBody, _ := json.Marshal(server.Request{Expression: tt.exp})
		req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CheckMidlware(server.ExpressionCalcHandler))
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != 200 {
			t.Errorf("handler returned wrong status code: got %v want %v", status, 200)
		}
		err = json.Unmarshal(rr.Body.Bytes(), &temp)
		if err != nil {
			t.Fatal(err)
		}
		if temp.Result != tt.output {
			t.Errorf("want %v, but got %v", tt.output, temp.Result)
		}
	}
}

func TestCheckMidlware_NoBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/calculate", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.CheckMidlware(server.ExpressionCalcHandler))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 500 {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
