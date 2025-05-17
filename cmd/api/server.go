package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	moviedb "github.com/kartik7120/booking_moviedb_service/cmd/grpcServer"
	"github.com/kartik7120/booking_moviedb_service/cmd/models"
	log "github.com/sirupsen/logrus"
)

// var validate *validator.Validate

type MoviedbService struct {
	moviedb.UnimplementedMovieDBServiceServer
	MovieDB *MovieDB
}

func NewMoviedbService() *MoviedbService {
	// validate = validator.New()
	return &MoviedbService{}
}

func (m *MoviedbService) AddMovie(ctx context.Context, in *moviedb.Movie) (*moviedb.MovieResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.MovieResponse{
	// 		Status:  408,
	// 		Message: "Context was cancelled",
	// 		Error:   "",
	// 	}, ctx.Err()
	// }

	castAndCrew := make([]models.CastAndCrew, 0)

	for _, cc := range in.CastCrew {
		c := models.CastAndCrew{
			Type:      string(cc.Type),
			Name:      cc.Name,
			Character: cc.CharacterName,
			PhotoURL:  cc.Photourl,
		}
		castAndCrew = append(castAndCrew, c)
	}

	venues := make([]models.Venue, 0)

	releaseDate, err := time.Parse("2006-01-02", in.ReleaseDate)

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "error parsing release date",
			Error:   err.Error(),
		}, err
	}

	movieTimeSlots := make([]models.MovieTimeSlot, 0)
	seats := make([]models.SeatMatrix, 0)

	for _, v := range in.Venues {

		for _, slot := range v.MovieTimeSlots {
			timestring, err := time.Parse("2006-01-02", slot.Date)

			if err != nil {
				return &moviedb.MovieResponse{
					Status:  400,
					Message: "error parsing slot date",
					Error:   err.Error(),
				}, err
			}

			st, err := time.Parse(time.RFC3339, slot.StartTime)

			if err != nil {
				return &moviedb.MovieResponse{
					Status:  400,
					Message: "error parsing slot start time",
					Error:   err.Error(),
				}, err
			}

			ed, err := time.Parse(time.RFC3339, slot.EndTime)

			if err != nil {
				return &moviedb.MovieResponse{
					Status:  400,
					Message: "error parsing slot start time",
					Error:   err.Error(),
				}, err
			}

			timeSlot := models.MovieTimeSlot{
				StartTime:   st,
				EndTime:     ed,
				Duration:    int(slot.Duration),
				Date:        timestring,
				MovieFormat: slot.MovieFormat.String(),
			}

			movieTimeSlots = append(movieTimeSlots, timeSlot)
		}

		for _, seat := range v.Seats {

			seat := models.SeatMatrix{
				SeatNumber: seat.SeatNumber,
				Type:       seat.Type.String(),
				Price:      int(seat.Price),
				Row:        int(seat.Row),
				Column:     int(seat.Column),
			}

			seats = append(seats, seat)
		}

		venue := models.Venue{
			Name:         v.Name,
			Type:         string(v.Type),
			Address:      v.Address,
			Rows:         int(v.Rows),
			Columns:      int(v.Columns),
			ScreenNumber: int(v.ScreenNumber),
			Longitude:    float64(v.Longitude),
			Latitude:     float64(v.Latitude),
		}
		venues = append(venues, venue)
	}

	movie := models.Movie{
		Title:           in.Title,
		Description:     in.Description,
		Duration:        int(in.Duration),
		Language:        in.Language,
		Type:            in.Type,
		CastCrew:        castAndCrew,
		PosterURL:       in.PosterUrl,
		TrailerURL:      in.TrailerUrl,
		ReleaseDate:     releaseDate,
		MovieResolution: in.MovieResolution,
		Venues:          venues,
	}

	_, status, err := m.MovieDB.AddMovie(movie, movieTimeSlots, seats)

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  int32(status),
			Message: "error adding movie",
			Movie:   in,
			Error:   "",
		}, err
	}

	return &moviedb.MovieResponse{
		Status:  200,
		Message: "Movie added successfully",
		Movie:   in,
		Error:   "",
	}, nil
}

func (m *MoviedbService) GetMovie(ctx context.Context, in *moviedb.MovieRequest) (*moviedb.MovieResponse, error) {

	log.Info("Starting gRPC GetMovie function call")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.MovieResponse{
	// 		Status:  408,
	// 		Message: "Context was cancelled",
	// 		Error:   "",
	// 	}, ctx.Err()
	// }

	movieID, err := strconv.ParseUint(in.Movieid, 10, 32)

	if err != nil {
		log.Error("Invalid movie id", err)
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "Invalid movie ID",
			Error:   err.Error(),
		}, nil
	}

	movie, status, err := m.MovieDB.GetMovieDetails(uint(movieID))

	if err != nil {
		log.Info("error calling get movie details function", err)
		return &moviedb.MovieResponse{
			Status:  int32(status),
			Message: "Movie not found",
			Error:   err.Error(),
		}, nil
	}

	castCrew := make([]*moviedb.CastAndCrew, 0)

	for _, cc := range movie.CastCrew {
		c := &moviedb.CastAndCrew{
			Type:          moviedb.CastAndCrewType(moviedb.CastAndCrewType_value[cc.Type]),
			Name:          cc.Name,
			CharacterName: cc.Character,
			Photourl:      cc.PhotoURL,
		}
		castCrew = append(castCrew, c)
	}

	return &moviedb.MovieResponse{
		Status:  200,
		Message: "Sucess",
		Movie: &moviedb.Movie{
			Title:           movie.Title,
			Description:     movie.Description,
			Duration:        int32(movie.Duration),
			Language:        movie.Language,
			Type:            movie.Type,
			CastCrew:        castCrew,
			PosterUrl:       movie.PosterURL,
			ReleaseDate:     movie.ReleaseDate.GoString(),
			TrailerUrl:      movie.TrailerURL,
			MovieResolution: movie.MovieResolution,
			Id:              int32(movie.ID),
			Votes:           int64(movie.Votes),
		},
	}, nil
}

func (m *MoviedbService) UpdateMovie(ctx context.Context, in *moviedb.Movie) (*moviedb.MovieResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.MovieResponse{
	// 		Status:  408,
	// 		Message: "Context was cancelled",
	// 		Error:   "",
	// 	}, ctx.Err()
	// }

	castAndCrew := make([]models.CastAndCrew, 0)

	for _, cc := range in.CastCrew {
		c := models.CastAndCrew{
			Type:      string(cc.Type),
			Name:      cc.Name,
			Character: cc.CharacterName,
			PhotoURL:  cc.Photourl,
		}
		castAndCrew = append(castAndCrew, c)
	}

	releaseDate, err := time.Parse("2006-01-02", in.ReleaseDate)

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "error parsing release date",
			Error:   err.Error(),
		}, err
	}

	movie := models.Movie{
		Title:           in.Title,
		Description:     in.Description,
		Duration:        int(in.Duration),
		Language:        in.Language,
		Type:            in.Type,
		CastCrew:        castAndCrew,
		PosterURL:       in.PosterUrl,
		TrailerURL:      in.TrailerUrl,
		ReleaseDate:     releaseDate,
		MovieResolution: in.MovieResolution,
	}

	movieID := in.Id

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "Invalid movie ID",
			Error:   err.Error(),
		}, nil
	}

	m.MovieDB.UpdateMovie(uint(movieID), movie)

	return &moviedb.MovieResponse{
		Status:  200,
		Message: "Movie updated successfully",
		Movie:   in,
		Error:   "",
	}, nil

}

func (m *MoviedbService) DeleteMovie(ctx context.Context, in *moviedb.MovieRequest) (*moviedb.MovieResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.MovieResponse{
	// 		Status:  408,
	// 		Message: "Context was cancelled",
	// 		Error:   "",
	// 	}, ctx.Err()
	// }

	movieID, err := strconv.ParseUint(in.Movieid, 10, 32)

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "Invalid movie ID",
			Error:   err.Error(),
		}, nil
	}

	status, err := m.MovieDB.DeleteMovie(uint(movieID))

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  int32(status),
			Message: "Movie not found",
			Error:   "",
		}, nil
	}

	return &moviedb.MovieResponse{
		Status:  200,
		Message: "Movie deleted successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) DeleteVenue(ctx context.Context, in *moviedb.MovieRequest) (*moviedb.MovieResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.MovieResponse{
	// 		Status:  408,
	// 		Message: "Context was cancelled",
	// 		Error:   ctx.Err().Error(),
	// 	}, ctx.Err()
	// }

	Venueid, err := strconv.ParseUint(in.Venueid, 10, 32)

	if err != nil {
		return &moviedb.MovieResponse{
			Status:  400,
			Message: "Invalid venue ID",
			Error:   err.Error(),
		}, err
	}

	status, err := m.MovieDB.DeleteVenue(uint(Venueid))

	if status != 200 || err != nil {
		return &moviedb.MovieResponse{
			Status:  int32(status),
			Message: "error deleting venue",
			Error:   err.Error(),
		}, err
	}

	return &moviedb.MovieResponse{
		Status:  int32(status),
		Message: "Venue deleted",
	}, nil
}

func (m *MoviedbService) UpdateVenue(ctx context.Context, in *moviedb.Venue) (*moviedb.VenueResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.VenueResponse{
	// 		Status:  500,
	// 		Message: "context is already cancelled",
	// 	}, ctx.Err()
	// }

	v := models.Venue{
		Name:         in.Name,
		Type:         in.Type.String(),
		Address:      in.Address,
		Rows:         int(in.Rows),
		Columns:      int(in.Columns),
		ScreenNumber: int(in.ScreenNumber),
		Longitude:    float64(in.Longitude),
		Latitude:     float64(in.Latitude),
	}

	movieFormatSupported := make([]string, 0)

	for _, val := range v.MovieFormatSupported {
		movieFormatSupported = append(movieFormatSupported, val)
	}

	languageSupported := make([]string, 0)

	for _, val := range v.LanguagesSupported {
		languageSupported = append(languageSupported, val)
	}

	if len(movieFormatSupported) > 0 {
		v.MovieFormatSupported = movieFormatSupported
	}

	if len(languageSupported) > 0 {
		v.LanguagesSupported = languageSupported
	}

	_, status, err := m.MovieDB.UpdateVenue(uint(in.Id), v)

	if status != 200 || err != nil {
		return &moviedb.VenueResponse{
			Status:  int32(status),
			Message: "error updating venue",
			Error:   err.Error(),
		}, err
	}

	return &moviedb.VenueResponse{
		Status:  200,
		Message: "Updated venue",
	}, nil
}

func (m *MoviedbService) AddVenue(ctx context.Context, in *moviedb.Venue) (*moviedb.VenueResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.VenueResponse{
	// 		Status:  500,
	// 		Message: "context is already cancelled",
	// 	}, nil
	// }

	venue := models.Venue{
		Name:         in.Name,
		Type:         in.Type.String(),
		Address:      in.Address,
		Rows:         int(in.Rows),
		Columns:      int(in.Columns),
		ScreenNumber: int(in.ScreenNumber),
		Longitude:    float64(in.Longitude),
		Latitude:     float64(in.Latitude),
	}

	movieFormatSupported := make([]string, 0)

	movieFormatSupported = append(movieFormatSupported, in.MovieFormatSupported...)

	if len(movieFormatSupported) > 0 {
		venue.MovieFormatSupported = movieFormatSupported
	}

	languageSupported := make([]string, 0)

	languageSupported = append(languageSupported, in.LanguageSupported...)

	if len(languageSupported) > 0 {
		venue.LanguagesSupported = languageSupported
	}

	seats := make([]models.SeatMatrix, 0)

	for _, val := range in.Seats {
		seat := models.SeatMatrix{
			SeatNumber: val.SeatNumber,
			Type:       val.Type.String(),
			Price:      int(val.Price),
			Row:        int(val.Row),
			Column:     int(val.Column),
		}

		seats = append(seats, seat)
	}

	if len(seats) > 0 {
		venue.Seats = seats
	}

	_, status, err := m.MovieDB.AddVenue(venue)

	if status != 200 || err != nil {
		return &moviedb.VenueResponse{
			Status:  int32(status),
			Message: "error adding a new venue",
			Error:   err.Error(),
		}, err
	}

	return &moviedb.VenueResponse{
		Status:  int32(status),
		Message: "added a new venue",
	}, nil
}

func (m *MoviedbService) GetVenue(ctx context.Context, in *moviedb.MovieRequest) (*moviedb.VenueResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// if ctx.Done() != nil {
	// 	return &moviedb.VenueResponse{
	// 		Status:  500,
	// 		Message: "context is already cancelled",
	// 	}, ctx.Err()
	// }

	Venueid, err := strconv.ParseUint(in.Venueid, 10, 32)

	if err != nil {
		return &moviedb.VenueResponse{
			Status: 500,
			Error:  err.Error(),
		}, err
	}

	venue, status, err := m.MovieDB.GetVenue(uint(Venueid))

	if status != 200 || err != nil {
		return &moviedb.VenueResponse{
			Status:  int32(status),
			Message: "error getting venue",
		}, err
	}

	return &moviedb.VenueResponse{
		Status:  200,
		Message: "success",
		Venue: &moviedb.Venue{
			Name:    venue.Name,
			Address: venue.Address,
			Rows:    int32(venue.Rows),
			Columns: int32(venue.Columns),
		},
	}, nil
}

func (m *MoviedbService) GetUpcomingMovies(ctx context.Context, in *moviedb.GetUpcomingMovieRequest) (*moviedb.GetUpcomingMovieResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	movies, status, err := m.MovieDB.GetUpcomingMovies(in.Date)

	if status != 200 {
		return &moviedb.GetUpcomingMovieResponse{
			Status:    int32(status),
			Message:   "error getting upcoming movies",
			MovieList: nil,
			Error:     err.Error(),
		}, nil
	}

	if err != nil {
		return &moviedb.GetUpcomingMovieResponse{
			Status:    int32(status),
			Message:   "error getting upcoming movies",
			MovieList: nil,
			Error:     err.Error(),
		}, nil
	}

	movielist := make([]*moviedb.Movie, 0)

	for _, v := range movies {
		cast_and_crew_arr := make([]*moviedb.CastAndCrew, 0)
		for _, cc := range v.CastCrew {

			cast_and_crew := &moviedb.CastAndCrew{
				Type:          moviedb.CastAndCrewType(moviedb.CastAndCrewType_value[cc.Type]),
				Name:          cc.Name,
				CharacterName: cc.Character,
				Photourl:      cc.PhotoURL,
			}

			cast_and_crew_arr = append(cast_and_crew_arr, cast_and_crew)
		}
		log.Info("movie id:", v.ID)
		movielist = append(movielist, &moviedb.Movie{
			Title:           v.Title,
			Description:     v.Description,
			Duration:        int32(v.Duration),
			Language:        v.Language,
			Type:            v.Type,
			PosterUrl:       v.PosterURL,
			TrailerUrl:      v.TrailerURL,
			ReleaseDate:     v.ReleaseDate.Local().String(),
			MovieResolution: v.MovieResolution,
			Votes:           int64(v.Votes),
			Ranking:         int32(v.Ranking),
			CastCrew:        cast_and_crew_arr,
			Id:              int32(v.ID),
		})
	}

	return &moviedb.GetUpcomingMovieResponse{
		Status:    200,
		Message:   "Success",
		MovieList: movielist,
		Error:     "",
	}, nil

}

func (m *MoviedbService) GetNowPlayingMovies(ctx context.Context, in *moviedb.GetNowPlayingMovieRequest) (*moviedb.GetUpcomingMovieResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	movies, status, err := m.MovieDB.GetNowPlayingMovies(
		int32(in.Longitude),
		int32(in.Latitude),
	)

	if status != 200 {
		return &moviedb.GetUpcomingMovieResponse{
			Status:    int32(status),
			Message:   "error getting now playing movies",
			MovieList: nil,
			Error:     err.Error(),
		}, nil
	}

	if err != nil {
		return &moviedb.GetUpcomingMovieResponse{
			Status:    int32(status),
			Message:   "error getting now playing movies",
			MovieList: nil,
			Error:     err.Error(),
		}, nil
	}

	movielist := make([]*moviedb.Movie, 0)

	for _, v := range movies {
		movielist = append(movielist, &moviedb.Movie{
			Title:           v.Title,
			Description:     v.Description,
			Duration:        int32(v.Duration),
			Language:        v.Language,
			Type:            v.Type,
			PosterUrl:       v.PosterURL,
			TrailerUrl:      v.TrailerURL,
			ReleaseDate:     v.ReleaseDate.Local().String(),
			MovieResolution: v.MovieResolution,
			Votes:           int64(v.Votes),
			Ranking:         int32(v.Ranking),
			Id:              int32(v.ID),
		})
	}

	return &moviedb.GetUpcomingMovieResponse{
		Status:    200,
		Message:   "Success",
		MovieList: movielist,
		Error:     "",
	}, nil
}

func (m *MoviedbService) AddReview(ctx context.Context, in *moviedb.Review) (*moviedb.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	review := models.Review{
		Rating:  int(in.Rating),
		Comment: in.Comment,
		UserID:  uint(in.UserID),
		MovieID: uint(in.MovieID),
		Title:   in.Title,
	}

	_, status, err := m.MovieDB.AddReview(review)

	if status != 200 || err != nil {
		return &moviedb.ReviewResponse{
			Status:  int32(status),
			Message: "error adding review",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.ReviewResponse{
		Status:  int32(status),
		Message: "review added successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) GetReview(ctx context.Context, in *moviedb.ReviewRequest) (*moviedb.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	review, status, err := m.MovieDB.GetReview(uint(in.UserID), uint(in.MovieID), uint(in.ReviewID))

	if err != nil {
		return &moviedb.ReviewResponse{
			Status:  int32(status),
			Error:   err.Error(),
			Message: "error getting review",
		}, err
	}

	r := &moviedb.Review{
		Rating:  int32(review.Rating),
		Comment: review.Comment,
		MovieID: int32(review.MovieID),
		Title:   review.Title,
		UserID:  int32(review.UserID),
	}

	return &moviedb.ReviewResponse{
		Status:  int32(status),
		Message: "success",
		Review:  r,
		Error:   "",
	}, nil

}

func (m *MoviedbService) UpdateReview(ctx context.Context, in *moviedb.ReviewUpdateRequest) (*moviedb.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	review, status, err := m.MovieDB.UpdateReview(in.Title, in.Comment, int(in.Rating), uint(in.UserID), uint(in.MovieID), uint(in.ReviewID))

	if status != 200 || err != nil {
		return &moviedb.ReviewResponse{
			Status:  int32(status),
			Message: "error updating review",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.ReviewResponse{
		Status:  int32(status),
		Message: "review updated successfully",
		Error:   "",
		Review: &moviedb.Review{
			Rating:  int32(review.Rating),
			Comment: review.Comment,
			MovieID: int32(review.MovieID),
			Title:   review.Title,
			UserID:  int32(review.UserID),
		},
	}, nil
}

func (m *MoviedbService) DeleteReview(ctx context.Context, in *moviedb.ReviewRequest) (*moviedb.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := m.MovieDB.DeleteReview(uint(in.UserID), uint(in.MovieID), uint(in.ReviewID))

	if status != 200 || err != nil {
		return &moviedb.ReviewResponse{
			Status:  int32(status),
			Message: "error deleting review",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.ReviewResponse{
		Status:  int32(status),
		Message: "review deleted successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) GetAllMovieReviews(ctx context.Context, in *moviedb.GetAllMovieReviewsRequest) (*moviedb.ReviewListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set default values for sortBy and filterBy if not provided
	sortBy := in.SortBy.String()
	if sortBy == "" {
		sortBy = "ASCENDING" // Default to ASCENDING
	}

	filterBy := in.FilterBy.String()
	if filterBy == "" {
		filterBy = "RATING" // Default to RATING
	}

	reviews, status, err := m.MovieDB.GetAllMovieReviews(uint(in.MovieID), int(in.Limit), int(in.Offset), sortBy, filterBy)

	if status != 200 || err != nil {
		return &moviedb.ReviewListResponse{
			Status:  int32(status),
			Message: "error getting reviews",
			Error:   err.Error(),
		}, nil
	}

	reviewList := make([]*moviedb.Review, 0)

	for _, review := range reviews.Reviews {
		reviewList = append(reviewList, &moviedb.Review{
			Rating:       int32(review.Rating),
			Comment:      review.Comment,
			MovieID:      int32(review.MovieID),
			Title:        review.Title,
			UserID:       int32(review.UserID),
			ReviewID:     int32(review.ID),
			CreatedAt:    int32(review.CreatedAt.UnixMilli()),
			ReviewerName: review.Username,
		})
	}

	return &moviedb.ReviewListResponse{
		Status:  200,
		Message: "success",
		Error:   "",
		ReviewList: &moviedb.ReviewList{
			Reviews: reviewList,
		},
		TotalReviewCount: int32(reviews.TotalReviews),
		TotalVotes:       int32(reviews.TotalVotes),
	}, nil
}

func (m *MoviedbService) GetMovieTimeSlots(ctx context.Context, in *moviedb.GetMovieTimeSlotRequest) (*moviedb.GetMovieTimeSlotResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	movieID, err := strconv.ParseUint(in.Movieid, 10, 32)

	if err != nil {
		return &moviedb.GetMovieTimeSlotResponse{
			Status:  400,
			Message: "Invalid movie ID",
			Error:   err.Error(),
		}, nil
	}

	if in.Latitude == 0 && in.Longitude == 0 {
		return &moviedb.GetMovieTimeSlotResponse{
			Status:  400,
			Message: "No latitude or longitude provided",
		}, nil
	}

	venues, timeSlots, status, err := m.MovieDB.GetMovieTimeSlots(in.StartDate, in.EndDate, uint(movieID), in.Latitude, in.Longitude)

	if status != 200 || err != nil {
		return &moviedb.GetMovieTimeSlotResponse{
			Status:  int32(status),
			Message: "error getting movie time slots",
			Error:   err.Error(),
		}, nil
	}

	timeSlotList := make([]*moviedb.MovieTimeSlot, 0)

	for _, v := range timeSlots {
		movieFormat := moviedb.SeatType_TWO_D

		if v.MovieFormat == "IMAX" {
			movieFormat = moviedb.SeatType_VIP
		} else if v.MovieFormat == "FOUR_D" {
			movieFormat = moviedb.SeatType_FOUR_D
		} else if v.MovieFormat == "THREE_D" {
			movieFormat = moviedb.SeatType_THREE_D
		} else if v.MovieFormat == "TWO_D" {
			movieFormat = moviedb.SeatType_TWO_D
		}

		st := fmt.Sprintf("%02d:%02d", v.StartTime.Hour(), v.StartTime.Minute())
		ed := fmt.Sprintf("%02d:%02d", v.EndTime.Hour(), v.EndTime.Minute())

		year, month, day := v.StartTime.Date()

		y := fmt.Sprintf("%02d", year)
		m := fmt.Sprintf("%02d", int(month))
		d := fmt.Sprintf("%02d", day)

		dt := fmt.Sprintf("%s", y+"-"+m+"-"+d)

		timeSlotList = append(timeSlotList, &moviedb.MovieTimeSlot{
			StartTime:   st,
			EndTime:     ed,
			Date:        dt,
			Duration:    int32(v.Duration),
			MovieFormat: movieFormat,
			Venueid:     int32(v.VenueID),
			Movieid:     int32(v.MovieID),
		})
	}

	venueArr := make([]*moviedb.Venue, 0)

	for _, v := range venues {
		venueArr = append(venueArr, &moviedb.Venue{
			Name:         v.Name,
			Address:      v.Address,
			Type:         moviedb.VenueType(moviedb.VenueType_value[v.Type]),
			Rows:         int32(v.Rows),
			Columns:      int32(v.Columns),
			Longitude:    float32(v.Longitude),
			Latitude:     float32(v.Latitude),
			ScreenNumber: int32(v.ScreenNumber),
			// Seats:        v.Seats,
			Id:                   int32(v.ID),
			MovieFormatSupported: v.MovieFormatSupported,
			LanguageSupported:    v.LanguagesSupported,
		})
	}

	if len(venueArr) == 0 && in.Latitude != 0 && in.Longitude != 0 {
		return &moviedb.GetMovieTimeSlotResponse{
			Status:  200,
			Message: "No venues found within 40km radius",
		}, nil
	}

	if len(venueArr) == 0 {
		return &moviedb.GetMovieTimeSlotResponse{
			Status:  200,
			Message: "No venues found",
		}, nil
	}

	return &moviedb.GetMovieTimeSlotResponse{
		Status:         200,
		Message:        "success",
		MovieTimeSlots: timeSlotList,
		Error:          "",
		Venues:         venueArr,
	}, nil
}

func (m *MoviedbService) AddMovieTimeSlot(ctx context.Context, in *moviedb.MovieTimeSlot) (*moviedb.MovieTimeSlotResponse, error) {
	fmt.Println("AddMovieTimeSlot function is called in server.go")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Printf("Movie time slot: %#v", in)

	d, err := time.Parse("2006-01-02", in.Date)

	if err != nil {
		return &moviedb.MovieTimeSlotResponse{
			Status:  400,
			Message: "error parsing date",
			Error:   err.Error(),
		}, err
	}

	st, err := time.Parse(time.RFC3339, in.StartTime)

	if err != nil {
		return &moviedb.MovieTimeSlotResponse{
			Status:  400,
			Message: "error parsing start time",
			Error:   err.Error(),
		}, err
	}

	ed, err := time.Parse(time.RFC3339, in.EndTime)

	if err != nil {
		return &moviedb.MovieTimeSlotResponse{
			Status:  400,
			Message: "error parsing end time",
			Error:   err.Error(),
		}, err
	}

	movieTimeSlot := models.MovieTimeSlot{
		StartTime:   st,
		EndTime:     ed,
		Date:        d,
		Duration:    int(in.Duration),
		MovieID:     uint(in.Movieid),
		VenueID:     uint(in.Venueid),
		MovieFormat: in.MovieFormat.String(),
	}

	_, status, err := m.MovieDB.AddMovieTimeSlot(movieTimeSlot)

	if status != 200 || err != nil {
		return &moviedb.MovieTimeSlotResponse{
			Status:  int32(status),
			Message: "error adding movie time slot",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.MovieTimeSlotResponse{
		Status:  200,
		Message: "movie time slot added successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) DeleteMovieTimeSlot(ctx context.Context, in *moviedb.MovieTimeSlotDelete) (*moviedb.MovieTimeSlotResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := m.MovieDB.DeleteMovieTimeSlot(uint(in.MovieTimeSlotId))

	if err != nil {
		return &moviedb.MovieTimeSlotResponse{
			Status:  int32(status),
			Message: "error deleting movie time slot",
			Error:   err.Error(),
		}, err
	}

	return &moviedb.MovieTimeSlotResponse{
		Status:  200,
		Message: "movie time slot deleted successfully",
		Error:   "",
	}, nil

}

func (m *MoviedbService) UpdateMovieTimeSlot(ctx context.Context, in *moviedb.MovieTimeSlotUpdate) (*moviedb.MovieTimeSlotUpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	d, err := time.Parse("2006-01-02", in.Date)

	if err != nil {
		return &moviedb.MovieTimeSlotUpdateResponse{
			Status:  400,
			Message: "error parsing date",
			Error:   err.Error(),
		}, err
	}

	st, err := time.Parse(time.RFC3339, in.StartTime)

	if err != nil {
		return &moviedb.MovieTimeSlotUpdateResponse{
			Status:  400,
			Message: "",
			Error:   err.Error(),
		}, err
	}

	ed, err := time.Parse(time.RFC3339, in.EndTime)

	if err != nil {
		return &moviedb.MovieTimeSlotUpdateResponse{
			Status:  400,
			Message: "",
			Error:   err.Error(),
		}, err
	}

	movieTimeSlot := models.MovieTimeSlot{
		StartTime:   st,
		EndTime:     ed,
		Date:        d,
		Duration:    int(in.Duration),
		MovieID:     uint(in.Movieid),
		VenueID:     uint(in.Venueid),
		MovieFormat: in.MovieFormat.String(),
	}

	_, status, err := m.MovieDB.UpdateMovieTimeSlot(uint(in.MovieTimeSlotId), movieTimeSlot)

	if status != 200 || err != nil {
		return &moviedb.MovieTimeSlotUpdateResponse{
			Status:  int32(status),
			Message: "error updating movie time slot",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.MovieTimeSlotUpdateResponse{
		Status:  200,
		Message: "movie time slot updated successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) GetSeatMatrix(ctx context.Context, in *moviedb.GetSeatMatrixRequest) (*moviedb.GetSeatMatrixResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	seatMatrix, status, err := m.MovieDB.GetSeatMatrix(int(in.Venueid))

	if status != 200 || err != nil {
		return &moviedb.GetSeatMatrixResponse{
			Status:  int32(status),
			Message: "error getting seat matrix",
			Error:   err.Error(),
		}, nil
	}

	var seats []*moviedb.SeatMatrix

	for _, v := range seatMatrix {
		seat := &moviedb.SeatMatrix{
			SeatNumber: v.SeatNumber,
			Type:       moviedb.SeatType(moviedb.SeatType_value[v.Type]),
			Price:      int32(v.Price),
			Row:        int32(v.Row),
			Column:     int32(v.Column),
			Id:         int32(v.ID),
		}

		seats = append(seats, seat)
	}

	return &moviedb.GetSeatMatrixResponse{
		Status:  200,
		Message: "success",
		Seats:   seats,
		Error:   "",
	}, nil
}

func (m *MoviedbService) AddSeatMatrix(ctx context.Context, in *moviedb.AddSeatMatrixInput) (*moviedb.AddSeatMatrixResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var seats []models.SeatMatrix

	for _, v := range in.Seats {

		seat := models.SeatMatrix{
			SeatNumber: v.SeatNumber,
			Type:       v.Type.String(),
			Price:      int(v.Price),
			Row:        int(v.Row),
			Column:     int(v.Column),
			VenueID:    uint(in.Venueid),
		}

		seats = append(seats, seat)
	}

	status, err := m.MovieDB.AddSeatMatrix(int(in.Venueid), seats)

	if status != 200 || err != nil {
		return &moviedb.AddSeatMatrixResponse{
			Status:  int32(status),
			Message: "error adding seat matrix",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.AddSeatMatrixResponse{
		Status:  200,
		Message: "seat matrix added successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) UpdateSeatMatrix(ctx context.Context, in *moviedb.UpdateSeatMatrixRequest) (*moviedb.UpdateSeatMatrixResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var seats []models.SeatMatrix

	for _, v := range in.Seats {
		seat := models.SeatMatrix{
			SeatNumber: v.SeatNumber,
			Type:       v.Type.String(),
			Price:      int(v.Price),
			Row:        int(v.Row),
			Column:     int(v.Column),
			VenueID:    uint(in.Venueid),
		}

		seats = append(seats, seat)
	}

	for _, v := range in.Seats {
		seat := models.SeatMatrix{
			SeatNumber: v.SeatNumber,
			Type:       v.Type.String(),
			Price:      int(v.Price),
			Row:        int(v.Row),
			Column:     int(v.Column),
			VenueID:    uint(in.Venueid),
		}

		_, status, err := m.MovieDB.UpdateSeatMatrix(uint(v.Id), seat)

		if status != 200 || err != nil {
			return &moviedb.UpdateSeatMatrixResponse{
				Status:  int32(status),
				Message: "error updating seat matrix",
				Error:   err.Error(),
			}, nil
		}
	}

	return &moviedb.UpdateSeatMatrixResponse{
		Status:  200,
		Message: "seat matrix updated successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) DeleteSeatMatrix(ctx context.Context, in *moviedb.DeleteSeatMatrixRequest) (*moviedb.DeleteSeatMatrixResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := m.MovieDB.DeleteSeatMatrix(uint(in.SeatMatrixId))

	if status != 200 || err != nil {
		return &moviedb.DeleteSeatMatrixResponse{
			Status:  int32(status),
			Message: "error deleting seat matrix",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.DeleteSeatMatrixResponse{
		Status:  200,
		Message: "seat matrix deleted successfully",
		Error:   "",
	}, nil
}

func (m *MoviedbService) DeleteEntireSeatMatrix(ctx context.Context, in *moviedb.DeleteEntireSeatMatrixRequest) (*moviedb.DeleteEntireSeatMatrixResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := m.MovieDB.DeleteEntireSeatMatrix(uint(in.Venueid))

	if status != 200 || err != nil {
		return &moviedb.DeleteEntireSeatMatrixResponse{
			Status:  int32(status),
			Message: "error deleting entire seat matrix",
			Error:   err.Error(),
		}, nil
	}

	return &moviedb.DeleteEntireSeatMatrixResponse{
		Status:  200,
		Message: "entire seat matrix deleted successfully",
		Error:   "",
	}, nil
}
