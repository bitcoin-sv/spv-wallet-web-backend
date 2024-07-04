package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Service is a service for fetching and caching BSV exchange rates.
type Service struct {
	exchangeRate *float64

	mutex     sync.Mutex
	lastFetch time.Time
}

// ExchangeRate is a struct that contains exchange rate data.
type ExchangeRate struct {
	Rate float64
}

// NewRatesService creates a new RatesService instance.
func NewRatesService(log *zerolog.Logger) *Service {
	s := &Service{
		exchangeRate: nil,
	}

	err := s.loadExchangeRate()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return s
}

// GetExchangeRate returns the current exchange rate.
func (s *Service) GetExchangeRate() (*float64, error) {
	err := s.loadExchangeRate()
	if err != nil {
		return nil, err
	}

	return s.exchangeRate, nil
}

func (s *Service) loadExchangeRate() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.useCachedValue() {
		return nil
	}

	exchangeRate, err := s.fetchExchangeRate()
	if err != nil {
		return err
	}

	s.lastFetch = time.Now()
	s.exchangeRate = exchangeRate

	return nil
}

func (s *Service) fetchExchangeRate() (*float64, error) {
	exchangeRateURL := viper.GetString(config.EnvEndpointsExchangeRate)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, exchangeRateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error during creating exchange rate request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during getting exchange rate: %w", err)
	}
	defer res.Body.Close() //nolint: all

	var exchangeRate *ExchangeRate
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error during reading response body: %w", err)
	}

	err = json.Unmarshal(bodyBytes, &exchangeRate)
	if err != nil {
		return nil, fmt.Errorf("error during unmarshalling response body: %w", err)
	}
	return &exchangeRate.Rate, nil
}

func (s *Service) useCachedValue() bool {
	if s.exchangeRate != nil && time.Since(s.lastFetch) < viper.GetDuration(config.EnvCacheSettingsTTL) {
		return true
	}
	return false
}
