package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	logger "trafficlightAPI/internal/middleware/logger"
	prometheus "trafficlightAPI/internal/middleware/prometheus"
	"trafficlightAPI/internal/models"

	"github.com/bytedance/sonic"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return sonic.ConfigDefault.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, userErr error, errs ...error) {
	response := models.ErrorResponse{
		Error: userErr.Error(),
	}

	for _, err := range errs {
		if err != nil {
			response.Details = append(response.Details, models.ErrorDetail{
				Message: err.Error(),
			})
		}
	}
	logger.LogError(status, userErr, errs...)

	if status >= 400 && status < 600 {
		errorType := strconv.Itoa(status / 100 * 100)
		prometheus.ErrorsAmount.WithLabelValues(errorType).Inc()
	}

	WriteJSON(w, status, response)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("отсутствует тело запроса")
	}
	return sonic.ConfigDefault.NewDecoder(r.Body).Decode(v)
}
