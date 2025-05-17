package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
	"github.com/kartik7120/booking_moviedb_service/cmd/models"
	"gorm.io/gorm"
)

type MovieDB struct {
	DB helper.DBConfig
}

var validate *validator.Validate

func NewMovieDB() *MovieDB {
	validate = validator.New()
	return &MovieDB{}
}

func (m *MovieDB) GetCurrentMovies(
	latitude float64,
	longitude float64,
) ([]models.Movie, int, error) {
	var movies []models.Movie
	var venues []models.Venue

	// Fetch all venues from the database
	result := m.DB.Conn.Find(&venues)

	if result.Error != nil {
		return movies, 500, result.Error
	}

	// Filter venues within 30km radius
	var nearbyVenues []models.Venue
	for _, venue := range venues {
		distance := helper.Haversine(latitude, longitude, venue.Latitude, venue.Longitude)
		if distance <= 30 {
			nearbyVenues = append(nearbyVenues, venue)
		}
	}

	// Fetch movies for the nearby venues
	for _, venue := range nearbyVenues {
		var venueMovies []models.Movie
		result := m.DB.Conn.Model(&venue).Association("Movies").Find(&venueMovies)

		if result.Error() != "" {
			return movies, 500, errors.New("error fetching movies for venue")
		}

		movies = append(movies, venueMovies...)
	}

	return movies, 200, nil
}

func (m *MovieDB) GetMovieDetails(movieID uint) (models.Movie, int, error) {
	var movie models.Movie
	result := m.DB.Conn.Preload("CastCrew").First(&movie, movieID)

	if result.Error != nil {
		return movie, 500, result.Error
	}
	return movie, 200, nil
}

func (m *MovieDB) GetMovieShowtimes(movieID uint, venueID uint, movie_format string, date string) ([]models.MovieTimeSlot, int, error) {
	var movie_time_slots []models.MovieTimeSlot

	result := m.DB.Conn.Preload("Venue").Preload("MovieTimeSlots").Preload("MovieTimeSlots.SeatLayout").Preload("MovieTimeSlots.SeatLayout.Seats").Find(&models.Movie{}, movieID)

	if result.Error != nil {
		return movie_time_slots, 500, result.Error
	}

	return movie_time_slots, 200, nil
}

func (m *MovieDB) GetMovieSeatLayout(movieID uint, venueID uint, movie_format string, date string, start_time string) (models.Venue, int, error) {
	var venue models.Venue
	result := m.DB.Conn.Where("id = ?", venueID).Find(&venue)

	if result.Error != nil {
		return venue, 500, result.Error
	}

	return venue, 200, nil
}

func (m *MovieDB) AddVenue(venue models.Venue) (models.Venue, int, error) {
	err := validate.Struct(venue)
	if err != nil {
		return venue, 400, err
	}

	result := m.DB.Conn.Create(&venue)
	if result.Error != nil {
		return venue, 500, result.Error
	}

	return venue, 200, nil
}

func (m *MovieDB) AddMovie(movie models.Movie, movieTimeSlots []models.MovieTimeSlot, seats []models.SeatMatrix) (models.Movie, int, error) {
	// Validate movie struct
	if err := validate.Struct(movie); err != nil {
		return movie, 400, err
	}

	// Start transaction
	tx := m.DB.Conn.Begin()
	if tx.Error != nil {
		return movie, 500, tx.Error
	}

	// Ensure rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Step 1: Insert Movie
	result := tx.Create(&movie)
	if result.Error != nil {
		tx.Rollback()
		return movie, 500, result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return movie, 500, errors.New("failed to insert movie, no rows affected")
	}

	// Step 2: Insert Venues
	for i := range movie.Venues {
		venue := &movie.Venues[i]

		// Create Venue

		fmt.Println("Inserted Venue ID:", venue.ID) // Debugging

		// Step 3: Insert MovieTimeSlots (from function parameter)
		for j := range movieTimeSlots {
			movieTimeSlots[j].MovieID = movie.ID
			movieTimeSlots[j].VenueID = venue.ID
		}

		if len(movieTimeSlots) > 0 {
			if err := tx.Create(&movieTimeSlots).Error; err != nil {
				tx.Rollback()
				return movie, 500, fmt.Errorf("error inserting time slots: %v", err)
			}
		}

		// Step 4: Insert Seat Matrices (from function parameter)
		for k := range seats {
			seats[k].VenueID = venue.ID
		}

		if len(seats) > 0 {
			if err := tx.Create(&seats).Error; err != nil {
				tx.Rollback()
				return movie, 500, fmt.Errorf("error inserting seat matrix: %v", err)
			}
		}
	}

	// Step 5: Commit transaction
	if err := tx.Commit().Error; err != nil {
		return movie, 500, fmt.Errorf("commit error: %v", err)
	}

	return movie, 200, nil
}

func (m *MovieDB) UpdateMovie(movieID uint, movie models.Movie) (models.Movie, int, error) {
	err := validate.Struct(movie)
	if err != nil {
		return movie, 400, err
	}

	var existingMovie models.Movie
	result := m.DB.Conn.First(&existingMovie, movieID)

	if result.Error != nil {
		return movie, 500, result.Error
	}

	result = m.DB.Conn.Model(&existingMovie).Updates(&movie)

	if result.Error != nil {
		return movie, 500, result.Error
	}

	return movie, 200, nil
}

func (m *MovieDB) DeleteMovie(movieID uint) (int, error) {

	result := m.DB.Conn.Unscoped().Delete(&models.Movie{}, movieID)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

func (m *MovieDB) DeleteVenue(venueId uint) (int, error) {
	var venue models.Venue
	var seats models.SeatMatrix

	result := m.DB.Conn.Unscoped().Where("venue_id = ?", venueId).Delete(&seats)

	if result.Error != nil {
		return 500, result.Error
	}

	result = m.DB.Conn.Unscoped().Delete(&venue, venueId)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

func (m *MovieDB) UpdateVenue(venueId uint, venue models.Venue) (models.Venue, int, error) {
	err := validate.Struct(venue)
	if err != nil {
		return venue, 400, err
	}

	var existingVenue models.Venue
	result := m.DB.Conn.First(&existingVenue, venueId)

	if result.Error != nil {
		return venue, 500, result.Error
	}

	result = m.DB.Conn.Model(&existingVenue).Updates(&venue)

	if result.Error != nil {
		return venue, 500, result.Error
	}

	return venue, 200, nil
}

func (m *MovieDB) GetVenue(venueId uint) (models.Venue, int, error) {
	var venue models.Venue
	result := m.DB.Conn.First(&venue, venueId)

	if result.Error != nil {
		return venue, 500, result.Error
	}

	return venue, 200, nil
}

// Used to fetch upcoming movies based on the range date given by user,starting from date + 2 weeks to date + 2 weeks + 1 month
func (m *MovieDB) GetUpcomingMovies(date string) ([]models.Movie, int, error) {
	// Parse the input date
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, 400, err
	}

	// Calculate start and end dates
	startDate := d.AddDate(0, 0, 14)
	endDate := d.AddDate(0, 1, 14)

	// Query the database
	var movies []models.Movie
	result := m.DB.Conn.Table("movies").Where("release_date BETWEEN ? AND ?", startDate, endDate).Preload("CastCrew").Find(&movies)

	if result.Error != nil {
		return nil, 500, result.Error
	}

	// Return the movies
	return movies, 200, nil
}

func (m *MovieDB) GetNowPlayingMovies(longitude, latitude int32) ([]models.Movie, int, error) {
	today := time.Now().Truncate(24 * time.Hour)

	var movies []models.Movie

	if longitude == 0 && latitude == 0 {
		// If no coordinates are provided, fetch all movies released today or earlier
		err := m.DB.Conn.
			Joins("JOIN movie_time_slots mts ON mts.movie_id = movies.id").
			Where("movies.release_date <= ?", today).
			Where("DATE(mts.date) = ?", today).
			Group("movies.id").
			Find(&movies).Error

		if err != nil {
			return nil, 500, err
		}
	}

	// If coordinates are provided, fetch movies released today or earlier and within 30km of the coordinates

	err := m.DB.Conn.
		Joins("JOIN movie_time_slots mts ON mts.movie_id = movies.id").
		Joins("JOIN venues venue ON mts.venue_id = venue.id"). // Add this JOIN
		Where("movies.release_date <= ?", today).
		Where("DATE(mts.date) = ?", today).
		Where("ST_DistanceSphere(ST_MakePoint(?, ?), ST_MakePoint(venue.longitude, venue.latitude)) <= ?", longitude, latitude, 30000).
		Group("movies.id").
		Find(&movies).Error

	if err != nil {
		return nil, 500, err
	}
	return movies, 200, nil
}

func (m *MovieDB) AddReview(review models.Review) (models.Review, int, error) {
	err := validate.Struct(review)

	if err != nil {
		return review, 400, err
	}

	result := m.DB.Conn.Create(&review)
	if result.Error != nil {
		return review, 500, result.Error
	}

	return review, 200, nil
}

func (m *MovieDB) GetReview(userID, movieID, reviewID uint) (models.Review, int, error) {
	var review models.Review
	result := m.DB.Conn.Table("reviews").Where("user_id = ? AND movie_id = ? AND id = ?", userID, movieID, reviewID).First(&review)

	if result.Error != nil {
		return review, 500, result.Error
	}

	return review, 200, nil
}

func (m *MovieDB) UpdateReview(title string, comment string, rating int, userID, movieID, reviewID uint) (models.Review, int, error) {
	var review models.Review
	result := m.DB.Conn.Table("reviews").Where("user_id = ? AND movie_id = ? AND id = ?", userID, movieID, reviewID).First(&review)

	if result.Error != nil {
		return review, 500, result.Error
	}

	review.Title = title
	review.Comment = comment
	review.Rating = rating

	result = m.DB.Conn.Save(&review)

	if result.Error != nil {
		return review, 500, result.Error
	}

	return review, 200, nil
}

func (m *MovieDB) DeleteReview(userID, movieID, reviewID uint) (int, error) {
	var review models.Review
	result := m.DB.Conn.Table("reviews").Where("user_id = ? AND movie_id = ? AND id = ?", userID, movieID, reviewID).First(&review)

	if result.Error != nil {
		return 500, result.Error
	}

	result = m.DB.Conn.Unscoped().Delete(&review)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

type ReviewListResponse struct {
	Reviews []struct {
		gorm.Model
		MovieID uint `json:"movie_id"`
		// Also need to add movie id in the review table
		Rating   int    `json:"rating" gorm:"not null"` // rating out of 5
		Comment  string `json:"comment"`
		Title    string `json:"title"`
		UserID   uint   `json:"user_id"` // user who wrote the review
		Username string `json:"username"`
	} `json:"reviews"`
	TotalReviews int64 `json:"total_reviews"`
	TotalVotes   int64 `json:"total_votes"`
}

func (m *MovieDB) GetAllMovieReviews(movieID uint, limit int, offset int, sortBy string, filter string) (ReviewListResponse, int, error) {
	var reviews []models.Review
	var totalReviews int64
	var movie models.Movie
	// var user models.User

	err := m.DB.Conn.Raw("SELECT COUNT(*) FROM reviews WHERE movie_id = ?", movieID).Scan(&totalReviews).Error
	if err != nil {
		return ReviewListResponse{}, 500, err
	}

	// Load movie

	err = m.DB.Conn.First(&movie, movieID).Error
	if err != nil {
		return ReviewListResponse{}, 500, err
	}

	query := m.DB.Conn.Where("movie_id = ?", movieID)

	if filter == "RATING" {
		query = query.Order("rating DESC")
	} else if filter == "DATE" {
		query = query.Order("created_at DESC")
	}

	if sortBy == "ASCENDING" {
		query = query.Order("created_at ASC")
	} else if sortBy == "DESCENDING" {
		query = query.Order("created_at DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&reviews).Error; err != nil {
		return ReviewListResponse{}, 500, err
	}

	ReviewResults := make([]struct {
		gorm.Model
		MovieID  uint   `json:"movie_id"`
		Rating   int    `json:"rating" gorm:"not null"` // rating out of 5
		Comment  string `json:"comment"`
		Title    string `json:"title"`
		UserID   uint   `json:"user_id"` // user who wrote the review
		Username string `json:"username"`
	}, 0)

	for i := range reviews {
		// Fetch user details for each review
		var user models.User
		err := m.DB.Conn.Table("users").Where("id = ?", reviews[i].UserID).First(&user).Error
		if err != nil {
			return ReviewListResponse{}, 500, err
		}

		var reviewResult struct {
			gorm.Model
			MovieID uint `json:"movie_id"`
			// Also need to add movie id in the review table
			Rating   int    `json:"rating" gorm:"not null"` // rating out of 5
			Comment  string `json:"comment"`
			Title    string `json:"title"`
			UserID   uint   `json:"user_id"` // user who wrote the review
			Username string `json:"username"`
		}

		reviewResult.Model = reviews[i].Model
		reviewResult.MovieID = reviews[i].MovieID
		reviewResult.Rating = reviews[i].Rating
		reviewResult.Comment = reviews[i].Comment
		reviewResult.Title = reviews[i].Title
		reviewResult.UserID = reviews[i].UserID
		reviewResult.Username = user.Username
		reviewResult.CreatedAt = reviews[i].CreatedAt
		reviewResult.ID = reviews[i].ID

		ReviewResults = append(ReviewResults, reviewResult)
	}

	return ReviewListResponse{Reviews: ReviewResults, TotalReviews: totalReviews, TotalVotes: int64(movie.Votes)}, 200, nil
}

/*
GetMovieTimeSlots fetches movie time slots based on the given date range, movie ID, and venue ID

	startTime: The start date in "YYYY-MM-DD:HH:MM:SS" format
	endTime: The end date in "YYYY-MM-DD:HH:MM:SS" format
	movieID: The ID of the movie
	venueID: The ID of the venue
*/
func (m *MovieDB) GetMovieTimeSlots(startDate string, endDate string, movieID uint, latitude float32, longitude float32) ([]models.Venue, []models.MovieTimeSlot, int, error) {
	var movieTimeSlots []models.MovieTimeSlot
	var venues []models.Venue
	// Parse the input dates

	start, err := time.Parse(time.DateOnly, startDate)

	fmt.Println("Start date:", start)

	if err != nil {
		return nil, nil, 400, err
	}

	end, err := time.Parse(time.DateOnly, endDate)

	fmt.Println("End date:", end)

	end = helper.EndOfDay(end)

	if err != nil {
		return nil, nil, 400, err
	}

	// Query the database
	result := m.DB.Conn.Debug().
		Preload("MovieTimeSlots", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("movie_id = ? AND start_time >= ? AND end_time <= ?", movieID, start.UTC(), end.UTC())
		}).
		Joins("JOIN movie_time_slots mts ON mts.venue_id = venues.id").
		Where("mts.movie_id = ? AND mts.start_time >= ? AND mts.end_time <= ?", movieID, start.UTC(), end.UTC()).
		// Adding Haversine formula in ORDER BY to calculate distance and sort results
		Order(fmt.Sprintf(`
        6371 * acos(
            cos(radians(%f)) * cos(radians(venues.latitude)) *
            cos(radians(venues.longitude) - radians(%f)) +
            sin(radians(%f)) * sin(radians(venues.latitude))
        ) ASC`, latitude, longitude, latitude)).
		Group("venues.id").
		Find(&venues)

	for _, v := range venues {
		movieTimeSlots = append(movieTimeSlots, v.MovieTimeSlots...)
	}

	if result.Error != nil {
		return nil, nil, 500, result.Error
	}

	return venues, movieTimeSlots, 200, nil
}

func (m *MovieDB) UpdateMovieTimeSlot(movieTimeSlotID uint, updatedMovieTimeSlot models.MovieTimeSlot) (models.MovieTimeSlot, int, error) {
	err := validate.Struct(updatedMovieTimeSlot)
	if err != nil {
		return updatedMovieTimeSlot, 400, err
	}

	result := m.DB.Conn.First(&models.MovieTimeSlot{}, movieTimeSlotID)

	if result.Error != nil {
		return updatedMovieTimeSlot, 500, result.Error
	}

	result = m.DB.Conn.Model(&models.MovieTimeSlot{}).Where("id = ?", movieTimeSlotID).Updates(&updatedMovieTimeSlot)

	if result.Error != nil {
		return updatedMovieTimeSlot, 500, result.Error
	}

	return updatedMovieTimeSlot, 200, nil
}

func (m *MovieDB) DeleteMovieTimeSlot(movieTimeSlotID uint) (int, error) {

	var movieTimeSlot models.MovieTimeSlot

	result := m.DB.Conn.Unscoped().Where("id = ?", movieTimeSlotID).First(&movieTimeSlot)

	if result.Error != nil {
		return 500, result.Error
	}

	result = m.DB.Conn.Unscoped().Delete(&models.MovieTimeSlot{}, movieTimeSlotID)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

func (m *MovieDB) AddMovieTimeSlot(movieTimeSlot models.MovieTimeSlot) (models.MovieTimeSlot, int, error) {
	fmt.Println("calling AddMovieTimeSlot in moviedb.go file")

	err := validate.Struct(movieTimeSlot)
	if err != nil {
		return movieTimeSlot, 400, err
	}

	result := m.DB.Conn.Create(&movieTimeSlot)

	if result.Error != nil {
		return movieTimeSlot, 500, result.Error
	}

	return movieTimeSlot, 200, nil
}

func (m *MovieDB) AddSeatMatrix(venueID int, seatMatrix []models.SeatMatrix) (int, error) {

	for i := range seatMatrix {
		err := validate.Struct(seatMatrix[i])

		if err != nil {
			return 400, err
		}
	}

	for i := range seatMatrix {
		seatMatrix[i].VenueID = uint(venueID)
	}

	for _, v := range seatMatrix {
		var existingSeatMatrix models.SeatMatrix
		m.DB.Conn.Where("row = ? AND column = ? AND venue_id = ?", v.SeatNumber, v.Row, v.Column, v.VenueID).First(&existingSeatMatrix)

		if existingSeatMatrix.ID != 0 {
			return 400, errors.New("seat with same row and column already exists")
		}
	}

	result := m.DB.Conn.Create(&seatMatrix)

	if result.Error != nil && result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"idx_unique_seat\" (SQLSTATE 23505)" {
		return 400, errors.New("duplicate seat number found")
	}

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

func (m *MovieDB) GetSeatMatrix(venueID int) ([]models.SeatMatrix, int, error) {
	var seatMatrix []models.SeatMatrix
	result := m.DB.Conn.Where("venue_id = ?", venueID).Find(&seatMatrix)

	if result.Error != nil {
		return nil, 500, result.Error
	}

	return seatMatrix, 200, nil
}

func (m *MovieDB) UpdateSeatMatrix(seatMatrixID uint, updatedSeatMatrix models.SeatMatrix) (models.SeatMatrix, int, error) {

	// Use a transaction to ensure atomicity

	tx := m.DB.Conn.Begin()

	if tx.Error != nil {
		return updatedSeatMatrix, 500, tx.Error
	}

	// Ensure rollback on panic

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.First(&models.SeatMatrix{}, seatMatrixID)

	if result.Error != nil {
		tx.Rollback()
		return updatedSeatMatrix, 500, result.Error
	}

	result = tx.Model(&models.SeatMatrix{}).Where("id = ?", seatMatrixID).Updates(&updatedSeatMatrix)

	if result.Error != nil {
		tx.Rollback()
		return updatedSeatMatrix, 500, result.Error
	}

	// Commit the transaction

	if err := tx.Commit().Error; err != nil {
		return updatedSeatMatrix, 500, fmt.Errorf("commit error: %v", err)
	}

	return updatedSeatMatrix, 200, nil
}

func (m *MovieDB) DeleteSeatMatrix(seatMatrixID uint) (int, error) {
	var seatMatrix models.SeatMatrix

	result := m.DB.Conn.Unscoped().Where("id = ?", seatMatrixID).First(&seatMatrix)

	if result.Error != nil {
		return 500, result.Error
	}

	result = m.DB.Conn.Unscoped().Delete(&models.SeatMatrix{}, seatMatrixID)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}

func (m *MovieDB) DeleteEntireSeatMatrix(venueID uint) (int, error) {
	var seatMatrix []models.SeatMatrix

	result := m.DB.Conn.Unscoped().Where("venue_id = ?", venueID).Delete(&seatMatrix)

	if result.Error != nil {
		return 500, result.Error
	}

	return 200, nil
}
