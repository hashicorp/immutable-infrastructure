package tmdbApi

//different params ? ? to map to a series of structs
// given the structure of json response

//This will be our response array that contains an array of movies
// we can then make a movie struct to access an individual movie or to map these individual movies

//Deserialization
type Movie struct {
	Adult         bool    `json:"adult"`
	Backdrop_path string  `json:"backdrop_path"`
	Genre_id      []int   `json:"genre_ids"`
	ID            int     `json:"id"`
	Media_type    string  `json:"movie"`
	Title         string  `json:"title"`
	Overview      string  `json:"overview"`
	Popularity    float64 `json:"popularity"`
	Poster        string  `json:"poster_path"`
	Release_date  string  `json:"release_date"`
}

type GetTrendingMoviesResult struct {
	Page         int     `json:"page"`
	Movies       []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}
