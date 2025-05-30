package handlers_test

import (
	"testing"
	"trafficlightAPI/internal/handlers"
	"trafficlightAPI/internal/models"

	"github.com/pkg/errors"
)

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name        string
		req         models.TrafficRequest
		trafficType int
		wantErr     error
	}{
		{
			name:        "valid request",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 1, CurrentTime: intPtr(10)},
			trafficType: 1,
			wantErr:     nil,
		},
		{
			name:        "nil CurrentTime",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 1, CurrentTime: nil},
			trafficType: 1,
			wantErr:     handlers.ErrNoCurrentTime,
		},
		{
			name:        "zero CurrentState",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 0, CurrentTime: intPtr(10)},
			trafficType: 1,
			wantErr:     handlers.ErrNoCurrentState,
		},
		{
			name:        "empty UUID",
			req:         models.TrafficRequest{UUID: "", CurrentState: 1, CurrentTime: intPtr(10)},
			trafficType: 1,
			wantErr:     handlers.ErrNotValidData,
		},
		{
			name:        "CurrentState too high",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 4, CurrentTime: intPtr(10)},
			trafficType: 1,
			wantErr:     handlers.ErrNotValidData,
		},
		{
			name:        "negative CurrentTime",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 1, CurrentTime: intPtr(-1)},
			trafficType: 1,
			wantErr:     handlers.ErrNotValidData,
		},
		{
			name:        "CurrentTime too high",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 1, CurrentTime: intPtr(20)},
			trafficType: 1,
			wantErr:     handlers.ErrNotValidData,
		},
		{
			name:        "valid for pedestrian",
			req:         models.TrafficRequest{UUID: "123", CurrentState: 2, CurrentTime: intPtr(5)},
			trafficType: 3,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handlers.ValidateRequest(tt.req, tt.trafficType)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateRequest() error = '%v', wantErr '%v'", err, tt.wantErr)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
