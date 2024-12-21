package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// type Request struct {
// 	Expression string `json:"expression"`
// }

// type Response struct {
// 	Result float64 `json:"result,omitempty"`
// 	Error  string  `json:"error,omitempty"`
// }

// func calculateHandler(w http.ResponseWriter, r *http.Request) {
// 	var req Request
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
// 		return
// 	}

// 	result, err := calc.Calc(req.Expression)
// 	if err != nil {
// 		if err.Error() == "Invalid expression 1" || err.Error() == "Invalid expression 2" || err.Error() == "Invalid expression 3" {
// 			w.WriteHeader(http.StatusUnprocessableEntity)
// 			json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
// 		} else {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
// 		}
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(Response{Result: result}) // Используем число, а не строку
// }

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult string
		expectedCode   int
	}{
		{"Valid Expression", "2+2*2", `{"result":6}`, http.StatusOK},
		{"Invalid Expression", "2++2", `{"error":"Expression is not valid"}`, http.StatusUnprocessableEntity},
		{"Empty Expression", "", `{"error":"Expression is not valid"}`, http.StatusUnprocessableEntity},
		{"Invalid Characters", "2+2a", `{"error":"Expression is not valid"}`, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос с телом
			reqBody := map[string]string{
				"expression": tt.expression,
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(body))
			w := httptest.NewRecorder()

			// Вызываем обработчик
			calculateHandler(w, req)

			// Проверяем статус-код
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}

			// Декодируем ответ
			var got Response
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Декодируем ожидаемый результат
			var expected Response
			if err := json.Unmarshal([]byte(tt.expectedResult), &expected); err != nil {
				t.Fatalf("Failed to unmarshal expected result: %v", err)
			}

			// Сравниваем структуры, чтобы избежать проблем с форматированием
			if got != expected {
				t.Errorf("Expected response %+v, but got %+v", expected, got)
			}
		})
	}
}
