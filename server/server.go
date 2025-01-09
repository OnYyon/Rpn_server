package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/OnYyon/Rpn_server.git/logic"
)

type Request struct {
	Expression string `json:"expression"`
}

type ErrorsJson struct {
	Err string `json:"error"`
}

func ExpressionCalcHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	res, err := logic.Calc(req.Expression)
	if err != nil {
		t, _ := json.Marshal(ErrorsJson{Err: "Expression is not valid"})
		w.WriteHeader(422)
		fmt.Fprint(w, string(t))
		return
	}
	w.WriteHeader(200)
	fmt.Println(res)
}

func CheckMidlware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if !checkExpression(req.Expression) {
			t, _ := json.Marshal(ErrorsJson{Err: "Expression is not valid"})
			w.WriteHeader(422)
			fmt.Fprint(w, string(t))
			return
		}
		next(w, r)
	}
}

func checkExpression(exp string) bool {
	re := regexp.MustCompile("^[0-9+-/*()/s]+$")
	return re.MatchString(exp)
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, guest!")
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", CheckMidlware(ExpressionCalcHandler))
	mux.HandleFunc("/", SimpleHandler)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

/*
curl http://127.0.0.1:8080/api/v1/calculate --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}'
*/
