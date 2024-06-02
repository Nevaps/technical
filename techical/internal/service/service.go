package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"strings"
	"sync"
	"techical/internal/config"
	"techical/internal/domain"
	"techical/internal/service/structs"
	"time"
)

const (
	FETCH_USD_PRICES    = "https://api.fastforex.io/fetch-all"
	FETCH_CRYPTO_PRICES = "https://api.fastforex.io/crypto/fetch-prices"
)

type CurrencyService struct {
	config       *config.Config
	log          *zerolog.Logger
	client       *resty.Client
	mutex        sync.RWMutex
	repository   CurrencyRepository
	fiatPrices   map[string]float64
	cryptoPrices map[string]float64
}

type CurrencyRepository interface {
	GetAllAvailable(context.Context) ([]domain.Currency, error)
}

func NewCurrencyService(config *config.Config, logger *zerolog.Logger, repository CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		config:     config,
		log:        logger,
		repository: repository,
		client:     resty.New(),
	}
}

func (s *CurrencyService) Run(ctx context.Context) {
	go s.UpdateRatesPeriodically(ctx)
}

func (s *CurrencyService) UpdateRatesPeriodically(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fiatPrices := s.getFiatPrices()
			cryptoPrices := s.getCryptoPrices()

			s.mutex.Lock()
			if fiatPrices != nil {
				s.fiatPrices = fiatPrices
			}
			if cryptoPrices != nil {
				s.cryptoPrices = cryptoPrices
			}
			s.mutex.Unlock()

			s.log.Info().Msgf("Prices have been updated fiat prices count: %d,  crypto prices count: %d", len(fiatPrices), len(cryptoPrices))

			ticker.Reset(time.Second * 30)
		}
	}
}

func (s *CurrencyService) getFiatPrices() map[string]float64 {
	resp, err := s.client.R().SetQueryParam("api_key", s.config.FastForexAPIKey).Get(FETCH_USD_PRICES)

	if err != nil {
		s.log.Error().Err(err).Msg("failed to get fiat prices")
		return nil
	}

	fiatResponse := structs.FiatPrices{}
	if err = json.Unmarshal(resp.Body(), &fiatResponse); err != nil {
		s.log.Error().Err(err).Msg("failed unmarshalling fiat response")
		return nil
	}

	return fiatResponse.Prices
}

func (s *CurrencyService) getCryptoPrices() map[string]float64 {
	crypto, err := s.repository.GetAllAvailable(context.Background())
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to get crypto pairs from db: %s", err)
		return nil
	}

	cryptoPairs := make([]string, len(crypto))

	for _, ticker := range crypto {
		cryptoPairs = append(cryptoPairs, ticker.Ticker+"/USD")
	}

	fmt.Println(cryptoPairs)

	cryptoRates := make(map[string]float64, len(crypto))

	for i := 0; i < (len(cryptoPairs)-1)/10+1; i++ {
		maxRight := (i + 1) * 10
		if ((i + 1) * 10) > len(cryptoPairs) {
			maxRight = len(cryptoPairs)
		}

		result := strings.Join(cryptoPairs[i*10:maxRight], ",")

		resp, err := s.client.R().SetQueryParams(map[string]string{
			"api_key": s.config.FastForexAPIKey, "pairs": result,
		}).Get(FETCH_CRYPTO_PRICES)

		if err != nil {
			s.log.Error().Err(err).Msg("failed to get fiat prices")
			return nil
		}

		cryptoResponse := structs.CryptoPrices{}
		if err = json.Unmarshal(resp.Body(), &cryptoResponse); err != nil {
			s.log.Error().Err(err).Msg("error unmarshalling fiat response")
			return nil
		}

		for key, value := range cryptoResponse.Prices {
			cryptoRates[strings.Split(key, "/")[0]] = value
		}

	}

	return cryptoRates
}

func (s *CurrencyService) Convert(from string, to string, amount float64) (float64, error) {
	defer s.mutex.RUnlock()

	s.mutex.RLock()

	var cryptoPrice, fiatPrice float64
	var okCrypto, okFiat bool

	if cryptoPrice, okCrypto = s.cryptoPrices[from]; okCrypto {
		fiatPrice, okFiat = s.fiatPrices[to]
		if !okFiat {
			return 0, fmt.Errorf("unsupported currency or C2C, F2F convert")
		}
	} else {
		cryptoPrice, okCrypto = s.cryptoPrices[to]
		fiatPrice, okFiat = s.fiatPrices[from]

		if !okCrypto || !okFiat {
			return 0, fmt.Errorf("unsupported currency or C2C, F2F convert")
		}
	}

	return amount * cryptoPrice * fiatPrice, nil
}
