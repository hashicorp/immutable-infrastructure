package handler

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"text/template"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/sabinlehaci/go-web-app/db"
	"github.com/sabinlehaci/go-web-app/tmdbApi"
	"golang.org/x/sync/errgroup"
)

//go:embed assets
var static embed.FS

type MovieGetter interface {
	GetTrendingMovies(ctx context.Context) (*tmdbApi.GetTrendingMoviesResult, error)
	GetTopRatedTVShows(ctx context.Context) (*tmdbApi.ListTVResult, error)
	GetTVDetails(ctx context.Context, tvID int) (*tmdbApi.GetTVDetailsResult, error)
	GetTVSeason(ctx context.Context, tvID int, seasonNumber int) (*tmdbApi.GetTVSeasonResult, error)
}

type Handlers struct {
	MovieGetter MovieGetter
	DB          db.Querier
}

var indexHTMLTemplate = template.Must(template.ParseFS(static, "assets/index.html"))

type templateVariables struct {
	Favorites []db.FavoriteMovie
	Movie     tmdbApi.Movie
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/top_seasons":
		h.ServeListTVSeasons(w, r)
		return
	default:
		h.ServeMovieNightFrontPage(w, r)
	}
}

func (h *Handlers) ServeMovieNightFrontPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		tmdbID := r.FormValue("movie_id")
		title := r.FormValue("title")
		_, err := h.DB.AddMovie(r.Context(), db.AddMovieParams{
			TmdbID: tmdbID,
			Title:  title,
		})
		if err != nil {
			var e *pgconn.PgError
			if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
				http.Error(w, "that movie is already your favorite", http.StatusBadRequest)
				return
			}
			http.Error(w, fmt.Sprintf("failed to add movie to favorites: %v", err), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		response, err := h.MovieGetter.GetTrendingMovies(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get trending movies: %v", err), http.StatusInternalServerError)
			return
		}

		movies, err := h.DB.ListMovies(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get favorite movies: %v", err), http.StatusInternalServerError)
			return
		}

		//logic for randomizing movies on every refresh
		randomIndex := rand.Intn(len(response.Movies))
		vars := templateVariables{
			Movie:     response.Movies[randomIndex],
			Favorites: movies,
		}
		err = indexHTMLTemplate.Execute(w, vars)
		if err != nil {
			// This is kinda hopeless
			http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handlers) ServeListTVSeasons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "you must use GET to list TV seasons", http.StatusMethodNotAllowed)
		return
	}
	tvShows, err := h.MovieGetter.GetTopRatedTVShows(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get top rated TV shows: %v", err), http.StatusInternalServerError)
		return
	}
	tvEg, ctx := errgroup.WithContext(r.Context())
	sNamesSlice := make(chan []string, len(tvShows.Results))
	for _, tvShow := range tvShows.Results {
		tvShow := tvShow
		tvEg.Go(func() error {
			var seasonNames []string
			tvDetails, err := h.MovieGetter.GetTVDetails(ctx, tvShow.ID)
			if err != nil {
				return fmt.Errorf("failed to get TV show details: %w", err)
			}
			seasonEg, ctx := errgroup.WithContext(ctx)
			sNames := make(chan string, len(tvDetails.Seasons))
			for _, season := range tvDetails.Seasons {
				season := season
				seasonEg.Go(func() error {
					tvSeasonDetails, err := h.MovieGetter.GetTVSeason(ctx, tvShow.ID, season.SeasonNumber)
					if err != nil {
						return fmt.Errorf("failed to get TV season details: %w", err)
					}
					sNames <- tvSeasonDetails.Name
					return nil
				})
			}
			err = seasonEg.Wait()
			if err != nil {
				return err
			}
			close(sNames)
			for seasonName := range sNames {
				seasonNames = append(seasonNames, seasonName)
			}
			return nil
		})
	}
	err = tvEg.Wait()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	close(sNamesSlice)
	var seasonNames []string
	for seasonNameSlice := range sNamesSlice {
		for _, seasonName := range seasonNameSlice {
			seasonNames = append(seasonNames, seasonName)
		}
	}
	w.Write([]byte(strings.Join(seasonNames, ",\n")))
}
