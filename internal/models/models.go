package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"trafficlightAPI/internal/image_generator"

	"github.com/bytedance/sonic"
)

type TrafficRequest struct {
	UUID         string `json:"uuid"`
	CurrentState int    `json:"current_state"`
	CurrentTime  *int   `json:"current_time"` // Указатель для проверки существования
	NeedImage    bool   `json:"need_image,omitempty"`
}

type TrafficResponse struct {
	UUID              string `json:"uuid"`
	NextState         string `json:"next_state"`
	NextCountdownTime string `json:"next_countdown_time,omitempty"`
	Image             string `json:"image,omitempty"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type ErrorResponse struct {
	Error   string        `json:"error"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type TrafficLight interface {
	GetNextState(TrafficRequest) (TrafficResponse, error)
}

type RegularTrafficLight struct {
	Data TrafficRequest
}

func (r *RegularTrafficLight) GetNextState(tr TrafficRequest) (TrafficResponse, error) {
	response := TrafficResponse{UUID: tr.UUID}

	if *tr.CurrentTime >= 19 {
		response.NextState = strconv.Itoa(tr.CurrentState%3 + 1)
	} else {
		response.NextState = strconv.Itoa(tr.CurrentState)
	}

	if tr.NeedImage {
		image, err := image_generator.TrafficLight1Image(tr.CurrentState)
		if err != nil {
			return TrafficResponse{}, fmt.Errorf("ошибка при создании изображения: %w", err)
		}
		response.Image = image
	}

	return response, nil
}

type TrafficLightWithRightArrow struct {
	Durations [7]int
}

func (r *TrafficLightWithRightArrow) GetNextState(tr TrafficRequest) (TrafficResponse, error) {
	response := TrafficResponse{UUID: tr.UUID}
	idx := (tr.CurrentState - 1) % len(r.Durations)

	currentDuration := r.Durations[idx]
	if *tr.CurrentTime < currentDuration-1 {
		response.NextState = strconv.Itoa(tr.CurrentState)
	} else {
		response.NextState = strconv.Itoa(tr.CurrentState%7 + 1)
	}

	if tr.NeedImage {
		image, err := image_generator.TrafficLight2Image(tr.CurrentState)
		if err != nil {
			return TrafficResponse{}, fmt.Errorf("ошибка при создании изображения: %w", err)
		}
		response.Image = image
	}

	return response, nil
}

type PedestrianTrafficLight struct {
	Durations [2]int
}

func (p *PedestrianTrafficLight) GetNextState(tr TrafficRequest) (TrafficResponse, error) {
	response := TrafficResponse{UUID: tr.UUID}
	idx := tr.CurrentState - 1
	currentDuration := p.Durations[idx]

	if *tr.CurrentTime < currentDuration-1 {
		response.NextState = strconv.Itoa(tr.CurrentState)
		response.NextCountdownTime = strconv.Itoa(currentDuration - *tr.CurrentTime)
	} else {
		nextState := tr.CurrentState%2 + 1
		response.NextState = strconv.Itoa(nextState)
		response.NextCountdownTime = strconv.Itoa(p.Durations[nextState-1])
	}

	return response, nil
}

var trafficLights = [3]TrafficLight{
	&RegularTrafficLight{},
	&TrafficLightWithRightArrow{
		Durations: [7]int{20, 20, 5, 10, 2, 20, 2},
	},
	&PedestrianTrafficLight{
		Durations: [2]int{20, 10},
	},
}

func ManageLights(data TrafficRequest, trafficType int) (json.RawMessage, error) {
	light := trafficLights[trafficType-1]
	nextState, err := light.GetNextState(data)
	if err != nil {
		return nil, err
	}

	response, err := sonic.Marshal(nextState)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании JSON ответа: %w", err)
	}

	return response, nil
}
