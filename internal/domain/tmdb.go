package domain

type MediaResult struct {
	Adult            bool     `json:"adult"`
	BackdropPath     string   `json:"backdrop_path"`
	ID               int      `json:"id"`
	OriginalName     string   `json:"original_name"`
	Name             string   `json:"name"`
	Overview         string   `json:"overview"`
	PosterPath       string   `json:"poster_path"`
	MediaType        string   `json:"media_type"`
	OriginalLanguage string   `json:"original_language"`
	GenreIDs         []int    `json:"genre_ids"`
	Popularity       float64  `json:"popularity"`
	FirstAirDate     string   `json:"first_air_date"`
	Softcore         bool     `json:"softcore"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	OriginCountry    []string `json:"origin_country"`
}
