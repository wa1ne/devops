package models_test

import (
	"testing"

	. "trafficlightAPI/internal/models"
)

func TestRegularTrafficLight(t *testing.T) {
	light := &RegularTrafficLight{}
	tests := []struct {
		req  TrafficRequest
		want TrafficResponse
		err  error
	}{
		{req: TrafficRequest{UUID: "1", CurrentState: 1, CurrentTime: intPtr(19)}, want: TrafficResponse{UUID: "1", NextState: "2"}},
		{req: TrafficRequest{UUID: "2", CurrentState: 3, CurrentTime: intPtr(18)}, want: TrafficResponse{UUID: "2", NextState: "3"}},
	}

	for _, tt := range tests {
		got, err := light.GetNextState(tt.req)
		if tt.err != nil && err == nil {
			t.Errorf("expected error %v, got nil", tt.err)
		}
		if tt.err == nil && got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}

func TestTrafficLightWithRightArrow(t *testing.T) {
	light := &TrafficLightWithRightArrow{Durations: [7]int{20, 20, 5, 10, 2, 20, 2}}
	tests := []struct {
		req  TrafficRequest
		want TrafficResponse
		err  error
	}{
		{req: TrafficRequest{UUID: "1", CurrentState: 1, CurrentTime: intPtr(19)}, want: TrafficResponse{UUID: "1", NextState: "2"}},
		{req: TrafficRequest{UUID: "2", CurrentState: 2, CurrentTime: intPtr(0)}, want: TrafficResponse{UUID: "2", NextState: "2"}},
		{req: TrafficRequest{UUID: "3", CurrentState: 7, CurrentTime: intPtr(19)}, want: TrafficResponse{UUID: "3", NextState: "1"}},
	}

	for _, tt := range tests {
		got, err := light.GetNextState(tt.req)
		if tt.err != nil && err == nil {
			t.Errorf("expected error %v, got nil", tt.err)
		}
		if tt.err == nil && got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}

func TestPedestrianTrafficLight(t *testing.T) {
	light := &PedestrianTrafficLight{Durations: [2]int{20, 10}}
	tests := []struct {
		req  TrafficRequest
		want TrafficResponse
		err  error
	}{
		{req: TrafficRequest{UUID: "1", CurrentState: 1, CurrentTime: intPtr(19)}, want: TrafficResponse{UUID: "1", NextState: "2", NextCountdownTime: "10"}},
		{req: TrafficRequest{UUID: "2", CurrentState: 2, CurrentTime: intPtr(5)}, want: TrafficResponse{UUID: "2", NextState: "2", NextCountdownTime: "5"}},
		{req: TrafficRequest{UUID: "3", CurrentState: 1, CurrentTime: intPtr(0)}, want: TrafficResponse{UUID: "3", NextState: "1", NextCountdownTime: "20"}},
	}

	for _, tt := range tests {
		got, err := light.GetNextState(tt.req)
		if tt.err != nil && err == nil {
			t.Errorf("expected error %v, got nil", tt.err)
		}
		if tt.err == nil && got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}

func TestManageLights(t *testing.T) {
	tests := []struct {
		req         TrafficRequest
		trafficType int
		wantErr     bool
	}{
		{req: TrafficRequest{UUID: "1", CurrentState: 1, CurrentTime: intPtr(20)}, trafficType: 1, wantErr: false},
		{req: TrafficRequest{UUID: "3", CurrentState: 1, CurrentTime: intPtr(5)}, trafficType: 3, wantErr: false},
	}

	for _, tt := range tests {
		_, err := ManageLights(tt.req, tt.trafficType)
		if tt.wantErr && err == nil {
			t.Errorf("expected error, got nil")
		}
		if !tt.wantErr && err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

// Helper function
func intPtr(i int) *int { return &i }
