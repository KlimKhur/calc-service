package main

import (
	"encoding/json"
	"net/http"

	calc "github.com/KlimKhur/calc-service/calc"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		return
	}

	result, err := calc.Calc(req.Expression)
	if err != nil {
		if err.Error() == "Invalid expression 1" || err.Error() == "Invalid expression 2" || err.Error() == "Invalid expression 3" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: result}) // Используем число, а не строку
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.ListenAndServe(":8080", nil)
}
