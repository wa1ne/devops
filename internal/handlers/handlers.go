package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
	_ "trafficlightAPI/docs"
	"trafficlightAPI/internal/config"
	"trafficlightAPI/internal/models"

	"github.com/pkg/errors"

	"github.com/bytedance/sonic"

	prometheus "trafficlightAPI/internal/middleware/prometheus"

	"github.com/go-chi/chi/v5"
	promm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	ErrUnmarshalingFromBody    = errors.New("ошибка при разборе JSON из тела")
	ErrUnmarshalingFromQuery   = errors.New("ошибка при разборе JSON из параметра")
	ErrNoType                  = errors.New("отсутствует параметр type")
	ErrInvalidTrafficlightType = errors.New("некорректный номер светофора")
)

// @Summary     Processing of traffic light control request
// @Tags        Trafficlight
// @Accept      json
// @Produce     json
// @Param       type query    int                    true "Type of the trafficlight" Enums(1, 2, 3)
// @Param       body body     models.TrafficRequest  true "Json request"
// @Success     200  {object} models.TrafficResponse      "Json response"
// @Failure     400  {object} models.ErrorResponse        "Invalid request data"
// @Failure     500  {object} models.ErrorResponse        "Server error"
// @Router      /trafficlight [post]
func ServeTrafficRoute(w http.ResponseWriter, r *http.Request) {
	var request models.TrafficRequest
	start := time.Now()

	if r.URL.Query().Has("data") {
		dataJSON := r.URL.Query().Get("data")
		if err := sonic.Unmarshal([]byte(dataJSON), &request); err != nil {
			WriteError(w, http.StatusBadRequest, ErrUnmarshalingFromQuery, err)
			return
		}
	} else {
		if err := ParseJSON(r, &request); err != nil {
			WriteError(w, http.StatusBadRequest, ErrUnmarshalingFromBody, err)
			return
		}
	}
	defer r.Body.Close()

	trafficTypeStr := r.URL.Query().Get("type")
	if trafficTypeStr == "" {
		WriteError(w, http.StatusBadRequest, ErrNoType)
		return
	}

	var trafficType int
	switch trafficTypeStr {
	case "1":
		trafficType = 1
	case "2":
		trafficType = 2
	case "3":
		trafficType = 3
	default:
		WriteError(w, http.StatusBadRequest, ErrInvalidTrafficlightType, fmt.Errorf("номер в запросе: %s", trafficTypeStr))
		return
	}

	if err := ValidateRequest(request, trafficType); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	prometheus.RequestedTypes.WithLabelValues(fmt.Sprintf("trafficlight%d", trafficType)).Inc()
	prometheus.RequestedTotal.Inc()

	imageRequested := "false"
	if request.NeedImage {
		imageRequested = "true"
	}

	if imageRequested == "true" {
		prometheus.ImageRequest.WithLabelValues("image_requested").Inc()
	} else {
		prometheus.ImageRequest.WithLabelValues("image_not_requested").Inc()
	}

	responseData, err := models.ManageLights(request, trafficType)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	duration := float64(time.Since(start).Seconds())
	prometheus.RequestDuration.WithLabelValues(
		trafficTypeStr,
		imageRequested,
	).Observe(duration)

	if err := WriteJSON(w, http.StatusOK, responseData); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при отправке JSON-ответа: %w. %+v", err, request))
		return
	}
}

func Run(cfg *config.Config, logger *slog.Logger) {
	router := chi.NewRouter()

	router.Use(prometheus.ResponseTimeMiddleware)

	router.Get("/trafficlight", ServeTrafficRoute)

	router.Get("/metrics", promhttp.InstrumentHandlerCounter(
		prometheus.RequestedTypes.MustCurryWith(promm.Labels{"type": "metrics"}),
		promhttp.Handler(),
	).ServeHTTP)

	router.Get("/docs/*", httpSwagger.WrapHandler.ServeHTTP)

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Info(
		"запуск HTTP сервера",
		slog.String("port", cfg.Server.Address),
		slog.String("env", cfg.Env),
	)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(
			"ошибка при запуске HTTP сервера",
			slog.Any("err", err),
		)
	}
}
