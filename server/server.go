package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/OnYyon/Rpn_server.git/logic"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResultJson struct {
	Result float64 `json:"result"`
}

type ErrorsJson struct {
	Err string `json:"error"`
}

var Req Request

func ExpressionCalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := logic.Calc(Req.Expression)
	if err != nil {
		t, err := json.Marshal(ErrorsJson{Err: "Expression is not valid"})
		if err != nil {
			http.Error(w, "Oppps somthing went wrong", 500)
			return
		}
		http.Error(w, string(t), 422)
		fmt.Println(422, 1)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ResultJson{Result: res})
}

func CheckMidlware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkExpression(Req.Expression) {
			fmt.Println(Req.Expression)
			t, err := json.Marshal(ErrorsJson{Err: "Expression is not valid"})
			if err != nil {
				http.Error(w, "Oppps something went wrong", 500)
				return
			}
			http.Error(w, string(t), 422)
			fmt.Println(422, 2)
			return
		}
		next(w, r)
	}
}

func MidlwareLogging(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				http.Error(w, "Oppps somthing went wrong", 500)
				return
			}
			err := json.NewDecoder(r.Body).Decode(&Req)
			if err != nil {
				http.Error(w, "Oppps somthing went wrong", 500)
				return
			}
			defer r.Body.Close()
			logger.Info("Http запрос",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("exp", Req.Expression))
			next.ServeHTTP(w, r)
		})
	}
}

func checkExpression(exp string) bool {
	re := regexp.MustCompile(`^[0-9+-\/*()\s]+$`)
	return re.MatchString(exp)
}

func SetupLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Error with setting logger %v\n", err)
	}
	return logger
}

func StartServer() {
	logger := SetupLogger()
	r := mux.NewRouter()
	r.Use(MidlwareLogging(logger))
	r.HandleFunc("/api/v1/calculate", CheckMidlware(ExpressionCalcHandler))
	logger.Info("Server start", zap.Int("port", 8080))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

/*
curl http://127.0.0.1:8080/api/v1/calculate --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}'
*/
