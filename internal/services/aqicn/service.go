package aqicn

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sync"
)

const (
	ApiUrl   = "https://api.waqi.info/feed/limassol/"
	statusOk = "ok"
)

type Service struct {
	data        Response
	dataRWMutex sync.RWMutex
	logger      *zap.Logger
	token       string
}

func (s *Service) Data() Response {
	return s.data
}

func (s *Service) Update() {
	s.logger.Debug("updating Air Quality data")

	apiUrl, err := url.Parse(ApiUrl)
	if err != nil {
		s.logger.Error("failed to parse url", zap.Error(err))
		return
	}
	query := apiUrl.Query()
	query.Add("token", s.token)
	apiUrl.RawQuery = query.Encode()

	res, err := http.Get(apiUrl.String())
	if err != nil {
		s.logger.Error("failed to fetch air quality data", zap.Error(err))
		return
	}

	defer res.Body.Close()
	s.dataRWMutex.RLock()
	defer s.dataRWMutex.RUnlock()
	var data Response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		s.logger.Error("failed to unmarshal air quality data ", zap.Error(err))
		return
	}
	if data.Status != statusOk {
		s.logger.Error("failed to fetch air quality data ", zap.Error(errors.New(data.Message)))
		return
	}
	s.data = data
	s.logger.Debug("Air Quality data updated")
}

func New(
	logger *zap.Logger,
	config *Config,
) *Service {
	return &Service{
		logger:      logger,
		dataRWMutex: sync.RWMutex{},
		token:       config.Token,
	}
}
