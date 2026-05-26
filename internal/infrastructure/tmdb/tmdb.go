package tmdb

import (
	"MovieTrackerBack/internal/domain"
	"MovieTrackerBack/internal/infrastructure/redis"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type TMDBClient struct {
	baseURL string
	apiKey  string
	client  *resty.Client
	redis   *redis.CacheRepository
}

func NewTMDBClient(apiKey string, baseURL string, redis *redis.CacheRepository) *TMDBClient {
	return &TMDBClient{apiKey: apiKey, baseURL: baseURL, client: resty.New(), redis: redis}
}

type tmdbSearchResponse struct {
	Results []domain.MediaResult `json:"results"`
}

func (c *TMDBClient) Search(query string) ([]domain.MediaResult, error) {
	// Normalizar la query para evitar duplicados en caché por mayúsculas o espacios
	cleanQuery := strings.ToLower(strings.TrimSpace(query))
	cacheKey := "tmdb_search_" + cleanQuery
	var searchResp tmdbSearchResponse

	// Intentar obtener de la caché primero
	cachedData, err := c.redis.Get(cacheKey)
	if err == nil && cachedData != nil {
		if jsonErr := json.Unmarshal(cachedData, &searchResp.Results); jsonErr == nil {
			return searchResp.Results, nil
		}
	}

	// Si no está en caché o hay un error, consultar a la API de TMDB
	resp, err := c.client.R().
		SetQueryParam("query", cleanQuery).
		SetQueryParam("include_adult", "true").
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&searchResp).
		Get(c.baseURL + "/search/multi")

	if err != nil {
		return nil, fmt.Errorf("error calling tmdb: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("tmdb api error: %s", resp.Status())
	}

	// Guardar el resultado en caché para futuras búsquedas (ej. 24 horas)
	if bytes, jsonErr := json.Marshal(searchResp.Results); jsonErr == nil {
		c.redis.Set(cacheKey, bytes, 24*time.Hour)
	}

	return searchResp.Results, nil
}

func (c *TMDBClient) Recomendations() ([]domain.MediaResult, error) {
	var searchResp tmdbSearchResponse

	// Recomendaciones
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&searchResp).
		Get(c.baseURL + "/person/popular/")

	if err != nil {
		return nil, fmt.Errorf("error calling tmdb: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("tmdb api error: %s", resp.Status())
	}
	return searchResp.Results, nil
}
