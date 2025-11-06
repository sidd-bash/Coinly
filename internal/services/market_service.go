package services

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

type MarketService struct {
    cache      map[string]float64
    lastUpdate time.Time
    mu         sync.Mutex
}

var market = &MarketService{
    cache: make(map[string]float64),
}

// Fetch from CoinGecko
func fetchPrices() (map[string]float64, error) {
    url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=inr"
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var data map[string]map[string]float64
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }

    prices := map[string]float64{
        "bitcoin":  data["bitcoin"]["inr"],
        "ethereum": data["ethereum"]["inr"],
    }
    return prices, nil
}

// Public getter with caching
func GetPrices() (map[string]float64, error) {
    market.mu.Lock()
    defer market.mu.Unlock()

    if time.Since(market.lastUpdate) < 60*time.Second && len(market.cache) > 0 {
        return market.cache, nil
    }

    prices, err := fetchPrices()
    if err != nil {
        return nil, err
    }

    market.cache = prices
    market.lastUpdate = time.Now()
    return prices, nil
}
