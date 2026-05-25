package domain

type MediaResult struct {
	adult             bool     `json:"adult"`
	backdrop_path     string   `json:"backdrop_path"`
	id                int      `json:"id"`
	original_name     string   `json:"original_name"`
	name              string   `json:"name"`
	overview          string   `json:"overview"`
	poster_path       string   `json:"poster_path"`
	media_type        string   `json:"media_type"`
	original_languaje string   `json:"original_language"`
	genre_ids         []int    `json:"genre_ids"`
	popularity        float64  `json:"popularity"`
	first_air_date    string   `json:"first_air_date"`
	softcore          bool     `json:"softcore"`
	vote_average      float64  `json:"vote_average"`
	vote_count        int      `json:"vote_count"`
	origin_country    []string `json:"origin_country"`
}
