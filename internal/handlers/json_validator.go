package handlers

import (
	"trafficlightAPI/internal/models"

	"github.com/pkg/errors"
)

var (
	trafficDurations  = [3]int{3, 7, 2}
	ErrNoCurrentTime  = errors.New("отсутствует поле current_time")
	ErrNoCurrentState = errors.New("отсутствует поле current_state")
	ErrNotValidData   = errors.New("некорректные входные данные")
)

func ValidateRequest(v models.TrafficRequest, trafficType int) error {
	if v.CurrentTime == nil {
		return ErrNoCurrentTime
	}
	if v.CurrentState == 0 {
		return ErrNoCurrentState
	}

	typeDurations := trafficDurations[trafficType-1]
	if v.UUID == "" || v.CurrentState < 1 || v.CurrentState > typeDurations || *v.CurrentTime < 0 || *v.CurrentTime > 19 {
		return errors.Wrapf(ErrNotValidData, "uuid:%s, current_state:%d, current_time:%d", v.UUID, v.CurrentState, *v.CurrentTime)
	}

	return nil
}
