package urls

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"trafficlightAPI/internal/handlers"
	"trafficlightAPI/internal/middleware/logger"
	"trafficlightAPI/internal/models"
)

func TestTrafficLightHandler(t *testing.T) {
	logger.InitLogger("../../logs/", "dev")

	tests := []struct {
		name       string
		method     string
		query      string
		body       interface{}
		wantStatus int
		wantErrMsg string
	}{
		// Валидные случаи
		{
			name:       "valid regular traffic light",
			method:     "GET",
			query:      "?type=1&data={\"uuid\":\"test1\",\"current_state\":1,\"current_time\":10}",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid right arrow boundary",
			method:     "GET",
			query:      "?type=2&data={\"uuid\":\"test2\",\"current_state\":7,\"current_time\":19}",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid pedestrian zero time",
			method:     "GET",
			query:      "?type=3&data={\"uuid\":\"test3\",\"current_state\":1,\"current_time\":0}",
			wantStatus: http.StatusOK,
		},
		// Некорректные параметры (400)
		{
			name:       "invalid current state",
			method:     "GET",
			query:      "?type=1&data={\"uuid\":\"test4\",\"current_state\":4,\"current_time\":10}",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrNotValidData.Error(),
		},
		{
			name:       "negative time",
			method:     "GET",
			query:      "?type=2&data={\"uuid\":\"test5\",\"current_state\":1,\"current_time\":-1}",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrNotValidData.Error(),
		},
		// Отсутствие параметров (400)
		{
			name:       "missing type",
			method:     "GET",
			query:      "?data={\"uuid\":\"test6\",\"current_state\":1,\"current_time\":10}",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrNoType.Error(),
		},
		{
			name:       "missing data",
			method:     "GET",
			query:      "?type=1",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrUnmarshalingFromBody.Error(),
		},
		{
			name:       "missing uuid",
			method:     "GET",
			query:      "?type=1&data={\"current_state\":1,\"current_time\":10}",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrNotValidData.Error(),
		},
		{
			name:       "typo in param",
			method:     "GET",
			query:      "?type=1&data={\"uuid\":\"test7\",\"currnt_state\":1,\"current_time\":10}",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: handlers.ErrNoCurrentState.Error(),
		},
		// Тест с JSON в теле
		{
			name:       "valid json body",
			method:     "GET",
			query:      "?type=1",
			body:       models.TrafficRequest{UUID: "test8", CurrentState: 1, CurrentTime: intPtr(10)},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/trafficlight"+tt.query, bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, "/trafficlight"+tt.query, nil)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.ServeTrafficRoute)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantErrMsg != "" {
				var resp models.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Errorf("failed to unmarshal error response: %v", err)
				}
				if !strings.Contains(resp.Error, tt.wantErrMsg) {
					t.Errorf("handler returned unexpected error: got %v want %v", resp.Error, tt.wantErrMsg)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
