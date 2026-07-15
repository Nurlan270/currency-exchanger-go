package helpers

import (
	http2 "currency-exchanger/internal/http"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
)

func SendError(w http.ResponseWriter, message string, statusCode int) {
	msg := map[string]string{
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Println(err)
	}
}

func Round(x float64) float64 {
	return math.Round(x*100) / 100
}

func ValidateEmpty(values ...string) error {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)

		if trimmed == "" {
			return http2.ErrEmptyValue
		}
	}

	return nil
}
