package tests

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/kartik7120/booking_moviedb_service/cmd/api"
	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
	"github.com/kartik7120/booking_moviedb_service/cmd/models"
	"github.com/lib/pq"
)

func TestMovieDB(t *testing.T) {
	t.Run("Add movie to database", func(t *testing.T) {

		if testing.Short() {
			t.Skip("Skipping this test in short mode")
		}

		err := godotenv.Load()

		if err != nil {
			t.Errorf("Could not load .env file")
		}

		m := api.NewMovieDB()

		// connect to database

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Errorf("unable to connect to database")
		}

		m.DB.Conn = conn

		_, err = time.Parse("2006-01-02", "2022-03-04")

		if err != nil {
			t.Errorf("error parsing release date")
			return
		}

		// movieTimeSlotDate, err := time.Parse("2006-01-02", "2025-04-05")

		if err != nil {
			t.Errorf("error parsing movie time slot date")
			return
		}

		// add movie to database
		// movie := models.Movie{
		// 	Title:       "The Batman",
		// 	Description: "The Batman is an upcoming American superhero film based on the DC Comics character Batman.",
		// 	ReleaseDate: "2022-03-04",
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 	Duration:    10560,               // should be in seconds,
		// 	Language:    []string{"English"}, // Correctly formatted array of strings
		// 	Type:        []string{"Action", "Crime", "Drama"},
		// 	CastCrew: []models.CastAndCrew{
		// 		{
		// 			Name:      "Robert Pattinson",
		// 			Type:      "Actor",
		// 			Character: "Bruce Wayne / Batman",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 		{
		// 			Name:      "Zoë Kravitz",
		// 			Type:      "Actor",
		// 			Character: "Selina Kyle / Catwoman",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 		{
		// 			Name:      "Paul Dano",
		// 			Type:      "Actor",
		// 			Character: "Edward Nashton / The Riddler",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 		{
		// 			Name:      "Jeffrey Wright",
		// 			Type:      "Actor",
		// 			Character: "James Gordon",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 		{
		// 			Name:      "Andy Serkis",
		// 			Type:      "Actor",
		// 			Character: "Alfred Pennyworth",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 		{
		// 			Name:      "Colin Farrell",
		// 			Type:      "Actor",
		// 			Character: "Oswald Cobblepot / The Penguin",
		// 			PhotoURL:  "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		// 		},
		// 	},
		// 	TrailerURL:      "https://www.youtube.com/watch?v=IhVf_3TeTQE",
		// 	MovieResolution: []string{"4K", "2K", "HD"},
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:      "PVR Cinemas",
		// 			Type:      "Multiplex",
		// 			Address:   "PVR Plaza, Connaught Place, New Delhi, Delhi 110001",
		// 			Latitude:  28.6315,
		// 			Longitude: 77.2167,
		// 			Rows:      10,
		// 			Columns:   10,
		// 			Seats: []models.SeatMatrix{
		// 				{
		// 					Row:        1,
		// 					Column:     1,
		// 					Price:      200,
		// 					SeatNumber: "A1",
		// 					IsBooked:   false,
		// 					Type:       "Regular",
		// 				},
		// 			},
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:           "The Lord of the Rings: The Fellowship of the Ring",
		// 	Description:     "A young hobbit, Frodo Baggins, embarks on a journey to destroy the One Ring and defeat the Dark Lord Sauron.",
		// 	ReleaseDate:     releaseDate,
		// 	PosterURL:       "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/6oom5QYQ2yQTMJIbnvbkBL9cHo6.jpg",
		// 	Duration:        178, // 2 hours 58 minutes
		// 	Language:        pq.StringArray([]string{"English", "Elvish", "Dwarvish"}),
		// 	Type:            pq.StringArray([]string{"Fantasy", "Adventure", "Drama"}),
		// 	MovieResolution: pq.StringArray([]string{"4K", "1080p", "720p"}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Elijah Wood", Character: "Frodo Baggins", PhotoURL: "https://example.com/elijah_wood.jpg"},
		// 		{Type: "Cast", Name: "Ian McKellen", Character: "Gandalf", PhotoURL: "https://example.com/ian_mckellen.jpg"},
		// 		{Type: "Cast", Name: "Viggo Mortensen", Character: "Aragorn", PhotoURL: "https://example.com/viggo_mortensen.jpg"},
		// 		{Type: "Cast", Name: "Sean Astin", Character: "Samwise Gamgee", PhotoURL: "https://example.com/sean_astin.jpg"},
		// 		{Type: "Crew", Name: "Peter Jackson", Character: "Director", PhotoURL: "https://example.com/peter_jackson.jpg"},
		// 	},
		// 	Ranking: 23,
		// 	Votes:   1002,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Rivendell Grand Theater",
		// 			Type:                 "IMAX",
		// 			Address:              "123 Elven Road, Middle-earth",
		// 			Latitude:             40.7128,
		// 			Longitude:            -74.0060,
		// 			Rows:                 20,
		// 			Columns:              30,
		// 			ScreenNumber:         1,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX", "3D", "2D"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English", "Elvish"}),

		// 			// Seats: []models.SeatMatrix{
		// 			// 	{Row: 1, Column: 1, Price: 1500, SeatNumber: "A1", IsBooked: false, Type: "Platinum"},
		// 			// 	{Row: 1, Column: 2, Price: 1500, SeatNumber: "A2", IsBooked: true, Type: "Platinum"},
		// 			// 	{Row: 2, Column: 1, Price: 1200, SeatNumber: "B1", IsBooked: false, Type: "Gold"},
		// 			// 	{Row: 2, Column: 2, Price: 1200, SeatNumber: "B2", IsBooked: true, Type: "Gold"},
		// 			// },

		// 			// MovieTimeSlots: []models.MovieTimeSlot{
		// 			// 	{
		// 			// 		StartTime:   "16:00", // 4:00 PM
		// 			// 		EndTime:     "19:00", // 7:00 PM
		// 			// 		Duration:    180,     // 3 hours
		// 			// 		Date:        time.Date(2025, 3, 22, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "IMAX",
		// 			// 	},
		// 			// 	{
		// 			// 		StartTime:   "20:00", // 8:00 PM
		// 			// 		EndTime:     "23:00", // 11:00 PM
		// 			// 		Duration:    180,     // 3 hours
		// 			// 		Date:        time.Date(2025, 3, 22, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "3D",
		// 			// 	},
		// 			// },
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:           "The Batman",
		// 	Description:     "Batman ventures into Gotham City's underworld when a sadistic killer leaves behind a trail of cryptic clues.",
		// 	ReleaseDate:     time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC),
		// 	PosterURL:       "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/74xTEgt7R36Fpooo50r9T25onhq.jpg",
		// 	Duration:        176, // 2 hours 56 minutes
		// 	Language:        pq.StringArray([]string{"English"}),
		// 	Type:            pq.StringArray([]string{"Action", "Crime", "Drama"}),
		// 	MovieResolution: pq.StringArray([]string{"4K", "1080p", "720p"}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Robert Pattinson", Character: "Bruce Wayne / Batman", PhotoURL: "https://example.com/robert_pattinson.jpg"},
		// 		{Type: "Cast", Name: "Zoë Kravitz", Character: "Selina Kyle / Catwoman", PhotoURL: "https://example.com/zoe_kravitz.jpg"},
		// 		{Type: "Cast", Name: "Paul Dano", Character: "Edward Nashton / Riddler", PhotoURL: "https://example.com/paul_dano.jpg"},
		// 		{Type: "Crew", Name: "Matt Reeves", Character: "Director", PhotoURL: "https://example.com/matt_reeves.jpg"},
		// 	},
		// 	Ranking: 45,
		// 	Votes:   1200,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Gotham Cineplex",
		// 			Type:                 "IMAX",
		// 			Address:              "200 Dark Alley, Gotham City",
		// 			Latitude:             40.7128,
		// 			Longitude:            -74.0060,
		// 			Rows:                 25,
		// 			Columns:              35,
		// 			ScreenNumber:         2,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX", "3D", "2D"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots
		// 			MovieTimeSlots: []models.MovieTimeSlot{
		// 				{
		// 					StartTime:   "18:00", // 6:00 PM
		// 					EndTime:     "21:00", // 9:00 PM
		// 					Duration:    180,     // 3 hours
		// 					Date:        time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "IMAX",
		// 				},
		// 				{
		// 					StartTime:   "22:00", // 10:00 PM
		// 					EndTime:     "01:00", // 1:00 AM
		// 					Duration:    180,     // 3 hours
		// 					Date:        time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "3D",
		// 				},
		// 			},
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:           "Blade Runner 2049",
		// 	Description:     "A young blade runner's discovery of a long-buried secret leads him to track down former blade runner Rick Deckard, who's been missing for thirty years.",
		// 	ReleaseDate:     time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		// 	PosterURL:       "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/gajva2L0rPYkEWjzgFlBXCAVBE5.jpg",
		// 	Duration:        164, // 2 hours 44 minutes
		// 	Language:        pq.StringArray([]string{"English"}),
		// 	Type:            pq.StringArray([]string{"Science Fiction", "Drama", "Mystery"}),
		// 	MovieResolution: pq.StringArray([]string{"4K", "1080p", "720p"}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Ryan Gosling", Character: "K", PhotoURL: "https://example.com/ryan_gosling.jpg"},
		// 		{Type: "Cast", Name: "Harrison Ford", Character: "Rick Deckard", PhotoURL: "https://example.com/harrison_ford.jpg"},
		// 		{Type: "Cast", Name: "Ana de Armas", Character: "Joi", PhotoURL: "https://example.com/ana_de_armas.jpg"},
		// 		{Type: "Crew", Name: "Denis Villeneuve", Character: "Director", PhotoURL: "https://example.com/denis_villeneuve.jpg"},
		// 	},
		// 	Ranking: 38,
		// 	Votes:   950,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Neo Los Angeles Cinema",
		// 			Type:                 "Dolby Atmos",
		// 			Address:              "2049 Cyber Street, Los Angeles",
		// 			Latitude:             34.0522,
		// 			Longitude:            -118.2437,
		// 			Rows:                 18,
		// 			Columns:              28,
		// 			ScreenNumber:         3,
		// 			MovieFormatSupported: pq.StringArray([]string{"Dolby Atmos", "IMAX", "2D"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots
		// 			MovieTimeSlots: []models.MovieTimeSlot{
		// 				{
		// 					StartTime:   "17:00", // 5:00 PM
		// 					EndTime:     "20:00", // 8:00 PM
		// 					Duration:    180,     // 3 hours
		// 					Date:        time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "Dolby Atmos",
		// 				},
		// 				{
		// 					StartTime:   "21:00", // 9:00 PM
		// 					EndTime:     "00:00", // Midnight
		// 					Duration:    180,     // 3 hours
		// 					Date:        time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "IMAX",
		// 				},
		// 			},
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:           "Spider-Man 2",
		// 	Description:     "Peter Parker struggles to balance his personal life and his responsibilities as Spider-Man while facing the formidable Doctor Octopus.",
		// 	ReleaseDate:     time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC),
		// 	PosterURL:       "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/olxpyq9kJAZ2NU1siLshhhXEPR7.jpg",
		// 	Duration:        127, // 2 hours 7 minutes
		// 	Language:        pq.StringArray([]string{"English"}),
		// 	Type:            pq.StringArray([]string{"Action", "Adventure", "Sci-Fi"}),
		// 	MovieResolution: pq.StringArray([]string{"4K", "1080p", "720p"}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Tobey Maguire", Character: "Peter Parker / Spider-Man", PhotoURL: "https://example.com/tobey_maguire.jpg"},
		// 		{Type: "Cast", Name: "Kirsten Dunst", Character: "Mary Jane Watson", PhotoURL: "https://example.com/kirsten_dunst.jpg"},
		// 		{Type: "Cast", Name: "Alfred Molina", Character: "Dr. Otto Octavius / Doc Ock", PhotoURL: "https://example.com/alfred_molina.jpg"},
		// 		{Type: "Crew", Name: "Sam Raimi", Character: "Director", PhotoURL: "https://example.com/sam_raimi.jpg"},
		// 	},
		// 	Ranking: 32,
		// 	Votes:   1100,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Queens Cineplex",
		// 			Type:                 "IMAX",
		// 			Address:              "20 Parker Street, New York City",
		// 			Latitude:             40.7306,
		// 			Longitude:            -73.9352,
		// 			Rows:                 22,
		// 			Columns:              32,
		// 			ScreenNumber:         4,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX", "3D", "2D"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots
		// 			MovieTimeSlots: []models.MovieTimeSlot{
		// 				{
		// 					StartTime:   "15:00", // 3:00 PM
		// 					EndTime:     "17:30", // 5:30 PM
		// 					Duration:    150,     // 2 hours 30 minutes
		// 					Date:        time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "IMAX",
		// 				},
		// 				{
		// 					StartTime:   "19:00", // 7:00 PM
		// 					EndTime:     "21:30", // 9:30 PM
		// 					Duration:    150,     // 2 hours 30 minutes
		// 					Date:        time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC),
		// 					MovieFormat: "3D",
		// 				},
		// 			},
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:       "Ford v Ferrari",
		// 	Description: "American car designer Carroll Shelby and fearless driver Ken Miles battle corporate interference and the laws of physics to build a revolutionary race car for Ford to challenge Ferrari at the 24 Hours of Le Mans in 1966.",
		// 	ReleaseDate: time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC), // Today's date
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/6ApDtO7xaWAfPqfi2IARXIzj8QS.jpg",
		// 	Duration:    152, // 2 hours 32 minutes
		// 	Language:    pq.StringArray([]string{"English", "Italian"}),
		// 	Type:        pq.StringArray([]string{"Action", "Biography", "Drama"}),
		// 	MovieResolution: pq.StringArray([]string{
		// 		"1080p", "720p", "480p",
		// 	}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Matt Damon", Character: "Carroll Shelby", PhotoURL: "https://example.com/matt_damon.jpg"},
		// 		{Type: "Cast", Name: "Christian Bale", Character: "Ken Miles", PhotoURL: "https://example.com/christian_bale.jpg"},
		// 		{Type: "Cast", Name: "Jon Bernthal", Character: "Lee Iacocca", PhotoURL: "https://example.com/jon_bernthal.jpg"},
		// 		{Type: "Crew", Name: "James Mangold", Character: "Director", PhotoURL: "https://example.com/james_mangold.jpg"},
		// 	},
		// 	Ranking: 8,
		// 	Votes:   920,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Classic Drive-In Theater",
		// 			Type:                 "2D",
		// 			Address:              "66 Le Mans Blvd, Los Angeles",
		// 			Latitude:             34.0522,
		// 			Longitude:            -118.2437,
		// 			Rows:                 20,
		// 			Columns:              30,
		// 			ScreenNumber:         2,
		// 			MovieFormatSupported: pq.StringArray([]string{"2D", "Dolby Digital"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots: []models.MovieTimeSlot{
		// 			// 	{
		// 			// 		StartTime:   "16:00",
		// 			// 		EndTime:     "18:32",
		// 			// 		Duration:    152,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "2D",
		// 			// 	},
		// 			// 	{
		// 			// 		StartTime:   "20:00",
		// 			// 		EndTime:     "22:32",
		// 			// 		Duration:    152,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "Dolby Digital",
		// 			// 	},
		// 			// },
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:       "Drive",
		// 	Description: "A mysterious Hollywood stuntman and mechanic moonlights as a getaway driver and finds himself in trouble when he helps out his neighbor.",
		// 	ReleaseDate: time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC), // Today's date
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/602vevIURmpDfzbnv5Ubi6wIkQm.jpg",
		// 	Duration:    100, // 1 hour 40 minutes
		// 	Language:    pq.StringArray([]string{"English"}),
		// 	Type:        pq.StringArray([]string{"Crime", "Drama", "Thriller"}),
		// 	MovieResolution: pq.StringArray([]string{
		// 		"1080p", "720p",
		// 	}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Ryan Gosling", Character: "Driver", PhotoURL: "https://example.com/ryan_gosling.jpg"},
		// 		{Type: "Cast", Name: "Carey Mulligan", Character: "Irene", PhotoURL: "https://example.com/carey_mulligan.jpg"},
		// 		{Type: "Cast", Name: "Bryan Cranston", Character: "Shannon", PhotoURL: "https://example.com/bryan_cranston.jpg"},
		// 		{Type: "Crew", Name: "Nicolas Winding Refn", Character: "Director", PhotoURL: "https://example.com/nicolas_refn.jpg"},
		// 	},
		// 	Ranking: 9,
		// 	Votes:   1500,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Neo Noir Cinema",
		// 			Type:                 "2D",
		// 			Address:              "42 Night Drive, Los Angeles, CA",
		// 			Latitude:             34.0522,
		// 			Longitude:            -118.2437,
		// 			Rows:                 15,
		// 			Columns:              25,
		// 			ScreenNumber:         3,
		// 			MovieFormatSupported: pq.StringArray([]string{"2D", "Dolby Digital"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots: []models.MovieTimeSlot{
		// 			// 	{
		// 			// 		StartTime:   "19:00",
		// 			// 		EndTime:     "20:40",
		// 			// 		Duration:    100,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "2D",
		// 			// 	},
		// 			// 	{
		// 			// 		StartTime:   "21:30",
		// 			// 		EndTime:     "23:10",
		// 			// 		Duration:    100,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "Dolby Digital",
		// 			// 	},
		// 			// },
		// 		},
		// 	},
		// }

		// movie := models.Movie{
		// 	Title:       "Tron: Legacy",
		// 	Description: "The son of a virtual world designer goes looking for his father and ends up inside the digital world that his father designed. He meets his father's corrupted creation and a unique ally who was born inside the digital world.",
		// 	ReleaseDate: time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC), // Today's date
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/9xkGLvAxu7f5PawQ6qJ4fF1wR0i.jpg",
		// 	Duration:    125, // 2 hours 5 minutes
		// 	Language:    pq.StringArray([]string{"English"}),
		// 	Type:        pq.StringArray([]string{"Action", "Science Fiction", "Adventure"}),
		// 	MovieResolution: pq.StringArray([]string{
		// 		"4K", "1080p", "720p",
		// 	}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Garrett Hedlund", Character: "Sam Flynn", PhotoURL: "https://example.com/garrett_hedlund.jpg"},
		// 		{Type: "Cast", Name: "Jeff Bridges", Character: "Kevin Flynn / Clu", PhotoURL: "https://example.com/jeff_bridges.jpg"},
		// 		{Type: "Cast", Name: "Olivia Wilde", Character: "Quorra", PhotoURL: "https://example.com/olivia_wilde.jpg"},
		// 		{Type: "Crew", Name: "Joseph Kosinski", Character: "Director", PhotoURL: "https://example.com/joseph_kosinski.jpg"},
		// 	},
		// 	Ranking: 11,
		// 	Votes:   1350,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Grid Central IMAX",
		// 			Type:                 "IMAX 3D",
		// 			Address:              "88 Lightcycle Ave, San Francisco, CA",
		// 			Latitude:             37.7749,
		// 			Longitude:            -122.4194,
		// 			Rows:                 30,
		// 			Columns:              40,
		// 			ScreenNumber:         7,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX 3D", "2D", "Dolby Atmos"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),

		// 			// MovieTimeSlots: []models.MovieTimeSlot{
		// 			// 	{
		// 			// 		StartTime:   "17:00",
		// 			// 		EndTime:     "19:05",
		// 			// 		Duration:    125,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "IMAX 3D",
		// 			// 	},
		// 			// 	{
		// 			// 		StartTime:   "20:00",
		// 			// 		EndTime:     "22:05",
		// 			// 		Duration:    125,
		// 			// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 			// 		MovieFormat: "Dolby Atmos",
		// 			// 	},
		// 			// },
		// 		},
		// 	},
		// }

		// st := time.Date(2025, 4, 5, 16, 0, 0, 0, time.UTC)
		// et := time.Date(2025, 4, 5, 19, 0, 0, 0, time.UTC)

		// st2 := time.Date(2025, 4, 6, 20, 0, 0, 0, time.UTC)
		// et2 := time.Date(2025, 4, 6, 23, 0, 0, 0, time.UTC)

		// _, status, err := m.AddMovie(movie, []models.MovieTimeSlot{
		// 	{
		// 		StartTime:   st,  // 4:00 PM
		// 		EndTime:     et,  // 7:00 PM
		// 		Duration:    180, // 3 hours
		// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 		MovieFormat: "IMAX",
		// 	},
		// 	{
		// 		StartTime:   st2, // 8:00 PM
		// 		EndTime:     et2, // 11:00 PM
		// 		Duration:    180, // 3 hours
		// 		Date:        time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC),
		// 		MovieFormat: "3D",
		// 	},
		// }, []models.SeatMatrix{
		// 	{Row: 1, Column: 1, Price: 1500, SeatNumber: "A1", Type: "Platinum"},
		// 	{Row: 1, Column: 2, Price: 1500, SeatNumber: "A2", Type: "Platinum"},
		// 	{Row: 2, Column: 1, Price: 1200, SeatNumber: "B1", Type: "Gold"},
		// 	{Row: 2, Column: 2, Price: 1200, SeatNumber: "B2", Type: "Gold"},
		// })

		// releaseDate := time.Now().AddDate(0, 2, 0) // 2 months from today

		// movie := models.Movie{
		// 	Title:       "Dune: Part Two",
		// 	Description: "Paul Atreides unites with Chani and the Fremen while seeking revenge against those who destroyed his family, and faces a choice between the love of his life and the fate of the universe.",
		// 	ReleaseDate: releaseDate,
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/8b8R8l88Qje9dn9OE8PY05Nxl1X.jpg",
		// 	Duration:    166, // 2 hours 46 minutes
		// 	Language:    pq.StringArray([]string{"English"}),
		// 	Type:        pq.StringArray([]string{"Science Fiction", "Adventure", "Drama"}),
		// 	MovieResolution: pq.StringArray([]string{
		// 		"4K", "1080p", "720p",
		// 	}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Timothée Chalamet", Character: "Paul Atreides", PhotoURL: "https://example.com/timothee_chalamet.jpg"},
		// 		{Type: "Cast", Name: "Zendaya", Character: "Chani", PhotoURL: "https://example.com/zendaya.jpg"},
		// 		{Type: "Cast", Name: "Rebecca Ferguson", Character: "Lady Jessica", PhotoURL: "https://example.com/rebecca_ferguson.jpg"},
		// 		{Type: "Crew", Name: "Denis Villeneuve", Character: "Director", PhotoURL: "https://example.com/denis_villeneuve.jpg"},
		// 	},
		// 	Ranking: 5,
		// 	Votes:   2450,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Arrakis Grand IMAX",
		// 			Type:                 "IMAX",
		// 			Address:              "101 Desert Way, Los Angeles, CA",
		// 			Latitude:             34.0522,
		// 			Longitude:            -118.2437,
		// 			Rows:                 25,
		// 			Columns:              35,
		// 			ScreenNumber:         1,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX", "3D", "Dolby Atmos"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),
		// 		},
		// 	},
		// }

		// // Movie time slots based on release date
		// st := time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 16, 0, 0, 0, time.UTC)
		// et := st.Add(3 * time.Hour)

		// st2 := st.AddDate(0, 0, 1).Add(4 * time.Hour) // next day, 8 PM
		// et2 := st2.Add(3 * time.Hour)

		// _, status, err := m.AddMovie(movie, []models.MovieTimeSlot{
		// 	{
		// 		StartTime:   st,
		// 		EndTime:     et,
		// 		Duration:    180,
		// 		Date:        releaseDate,
		// 		MovieFormat: "IMAX",
		// 	},
		// 	{
		// 		StartTime:   st2,
		// 		EndTime:     et2,
		// 		Duration:    180,
		// 		Date:        releaseDate.AddDate(0, 0, 1),
		// 		MovieFormat: "Dolby Atmos",
		// 	},
		// }, []models.SeatMatrix{})

		// releaseDate := time.Now().AddDate(0, 2, 0) // 2 months from today

		// movie := models.Movie{
		// 	Title:       "Avatar: The Way of Water",
		// 	Description: "Jake Sully lives with his newfound family on the planet of Pandora. Once a familiar threat returns, Jake must work with Neytiri and the army of the Na'vi to protect their home.",
		// 	ReleaseDate: releaseDate,
		// 	PosterURL:   "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/t6HIqrRAclMCA60NsSmeqe9RmNV.jpg",
		// 	Duration:    192, // 3 hours 12 minutes
		// 	Language:    pq.StringArray([]string{"English"}),
		// 	Type:        pq.StringArray([]string{"Action", "Adventure", "Science Fiction"}),
		// 	MovieResolution: pq.StringArray([]string{
		// 		"4K", "1080p", "720p",
		// 	}),
		// 	CastCrew: []models.CastAndCrew{
		// 		{Type: "Cast", Name: "Sam Worthington", Character: "Jake Sully", PhotoURL: "https://example.com/sam_worthington.jpg"},
		// 		{Type: "Cast", Name: "Zoe Saldaña", Character: "Neytiri", PhotoURL: "https://example.com/zoe_saldana.jpg"},
		// 		{Type: "Cast", Name: "Sigourney Weaver", Character: "Kiri", PhotoURL: "https://example.com/sigourney_weaver.jpg"},
		// 		{Type: "Crew", Name: "James Cameron", Character: "Director", PhotoURL: "https://example.com/james_cameron.jpg"},
		// 	},
		// 	Ranking: 3,
		// 	Votes:   3200,
		// 	Venues: []models.Venue{
		// 		{
		// 			Name:                 "Pandora Ocean Dome",
		// 			Type:                 "IMAX 3D",
		// 			Address:              "500 Reef Street, Honolulu, HI",
		// 			Latitude:             21.3069,
		// 			Longitude:            -157.8583,
		// 			Rows:                 28,
		// 			Columns:              38,
		// 			ScreenNumber:         2,
		// 			MovieFormatSupported: pq.StringArray([]string{"IMAX 3D", "Dolby Atmos", "3D"}),
		// 			LanguagesSupported:   pq.StringArray([]string{"English"}),
		// 		},
		// 	},
		// }

		// // Movie time slots based on release date
		// st := time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 17, 0, 0, 0, time.UTC)
		// et := st.Add(3 * time.Hour)

		// st2 := st.AddDate(0, 0, 1).Add(3 * time.Hour) // next day, 8 PM
		// et2 := st2.Add(3 * time.Hour)

		// _, status, err := m.AddMovie(movie, []models.MovieTimeSlot{
		// 	{
		// 		StartTime:   st,
		// 		EndTime:     et,
		// 		Duration:    180,
		// 		Date:        releaseDate,
		// 		MovieFormat: "IMAX 3D",
		// 	},
		// 	{
		// 		StartTime:   st2,
		// 		EndTime:     et2,
		// 		Duration:    180,
		// 		Date:        releaseDate.AddDate(0, 0, 1),
		// 		MovieFormat: "Dolby Atmos",
		// 	},
		// }, []models.SeatMatrix{
		// 	{Row: 1, Column: 1, Price: 2000, SeatNumber: "A1", Type: "Platinum"},
		// 	{Row: 1, Column: 2, Price: 2000, SeatNumber: "A2", Type: "Platinum"},
		// 	{Row: 2, Column: 1, Price: 1500, SeatNumber: "B1", Type: "Gold"},
		// 	{Row: 2, Column: 2, Price: 1500, SeatNumber: "B2", Type: "Gold"},
		// })

		releaseDate := time.Now().AddDate(0, 0, -3) // Three days ago

		movie := models.Movie{
			Title:           "Freakier Friday",
			Description:     "A fantasy comedy sequel where the mother and daughter swap bodies again under magical—and unpredictable—circumstances.",
			ReleaseDate:     releaseDate,
			PosterURL:       "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/your_poster_path.jpg", // replace with actual URL
			Duration:        111,                                                                        // minutes
			Language:        pq.StringArray([]string{"English"}),
			Type:            pq.StringArray([]string{"Fantasy", "Comedy", "Family"}),
			MovieResolution: pq.StringArray([]string{"4K", "1080p", "720p"}),
			CastCrew: []models.CastAndCrew{
				{Type: "Cast", Name: "Jamie Lee Curtis", Character: "Mom", PhotoURL: "https://example.com/jamie_lee_curtis.jpg"},
				{Type: "Cast", Name: "Lindsay Lohan", Character: "Daughter", PhotoURL: "https://example.com/lindsay_lohan.jpg"},
				{Type: "Crew", Name: "Nisha Ganatra", Character: "Director", PhotoURL: "https://example.com/nisha_ganatra.jpg"},
			},
			Ranking: 8,
			Votes:   2100,
			Venues: []models.Venue{
				{
					Name:                 "Magic Mirror Cinema",
					Type:                 "Digital",
					Address:              "123 Elm Street, Los Angeles, CA",
					Latitude:             34.0522,
					Longitude:            -118.2437,
					Rows:                 20,
					Columns:              30,
					ScreenNumber:         5,
					MovieFormatSupported: pq.StringArray([]string{"Digital", "3D", "Dolby Atmos"}),
					LanguagesSupported:   pq.StringArray([]string{"English"}),
				},
			},
		}

		st := time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 18, 30, 0, 0, time.UTC)
		et := st.Add(2*time.Hour - 30*time.Minute) // 2 hours duration

		_, status, err := m.AddMovie(movie, []models.MovieTimeSlot{
			{
				StartTime:   st,
				EndTime:     et,
				Duration:    int(et.Sub(st).Minutes()),
				Date:        releaseDate,
				MovieFormat: "Digital",
			},
		}, []models.SeatMatrix{
			{Row: 1, Column: 1, Price: 1200, SeatNumber: "A1", Type: "Gold"},
			{Row: 1, Column: 2, Price: 1200, SeatNumber: "A2", Type: "Gold"},
			{Row: 2, Column: 1, Price: 900, SeatNumber: "B1", Type: "Silver"},
			{Row: 2, Column: 2, Price: 900, SeatNumber: "B2", Type: "Silver"},
		})

		if err != nil {
			t.Error(err.Error())
			return
		}

		if status != 200 {
			t.Errorf("Status should be 200 after succesful addition of movies")
			return
		}

	})

	t.Run("Update a movie in database", func(t *testing.T) {

		if testing.Short() {
			t.Skip("Skipping this test in short mode")
		}

		err := godotenv.Load()

		if err != nil {
			t.Errorf("Error loading in .env file")
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Error("Error connecting to the database", err)
			return
		}

		m.DB.Conn = conn

		movieID := 23

		updateMovieObj := models.Movie{
			Title: "Blade Runner 2050",
		}

		_, status, err := m.UpdateMovie(uint(movieID), updateMovieObj)

		if status != 200 {
			t.Errorf("Movie should have been updated")
			return
		}

		if err != nil {
			t.Error("Error updating movies", err)
			return
		}

	})

	t.Run("Delete movie in database", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			t.Error("Failed to load .env file")
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Error("Failed to connect to the database")
			return
		}

		m.DB.Conn = conn

		movieID := 23

		status, err := m.DeleteMovie(uint(movieID))

		if status != 200 {
			t.Error("Movie should have been deleted with status 200")
			return
		}

		if err != nil {
			t.Error("Error delete movie from database", err)
			return
		}
	})

	t.Run("Delete venue in database", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			t.Error("error occured when loading .env file", err)
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Error("error connecting to the database", err)
			return
		}

		m.DB.Conn = conn

		venueID := 21

		m.DB.Conn.AutoMigrate(&models.SeatMatrix{})

		status, err := m.DeleteVenue(uint(venueID))

		if status != 200 {
			t.Error("status should have been", err)
			return
		}

		if err != nil {
			t.Error("error should have been nil", err)
			return
		}
	})

	t.Run("Update a venue in database", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			t.Error("error loading .env file", err)
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Error("error connecting to the database", err)
			return
		}

		m.DB.Conn = conn

		venueID := 4

		venue := models.Venue{
			ScreenNumber: 3,
		}

		_, status, err := m.UpdateVenue(uint(venueID), venue)

		if status != 200 {
			t.Error("status should be 200 when updating venue", err)
			return
		}

		if err != nil {
			t.Error("error should be nil", err)
			return
		}
	})

	t.Run("Add a venue to the database", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			t.Fatal("error loading .env file")
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()
		if err != nil {
			t.Fatalf("Error connecting to the database: %v", err) // Use t.Fatalf for fatal errors
		}

		m.DB.Conn = conn

		venue := models.Venue{
			Name:                 "IMAX Theater",
			Type:                 "Multiplex",
			Address:              "123 Movie Street, City",
			Rows:                 10,
			Columns:              20,
			ScreenNumber:         1,
			Longitude:            12.34,
			Latitude:             56.78,
			MovieFormatSupported: pq.StringArray{"2D", "3D", "IMAX"},
			LanguagesSupported:   pq.StringArray{"English", "Spanish"},
		}

		// Insert into DB
		result := m.DB.Conn.Create(&venue)
		if result.Error != nil {
			t.Errorf("Failed to add venue: %v", result.Error)
			return
		}

		// Verify venue exists
		var savedVenue models.Venue
		if err := m.DB.Conn.First(&savedVenue, venue.ID).Error; err != nil {
			t.Errorf("Venue was not saved in the database: %v", err)
		} else {
			t.Logf("Venue successfully added: %v", savedVenue)
		}
	})

	t.Run("Add venue along side movies in database", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			t.Fatal("error loading .env file", err)
			return
		}

		m := api.NewMovieDB()

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Fatal("error connecting to database", err)
			return
		}

		m.DB.Conn = conn

		// m.DB.Conn.Migrator().DropTable("movie_venues")
		m.DB.Conn.AutoMigrate(&models.Venue{}, &models.Venue{})

		venue := models.Venue{
			Name:                 "IMAX Theater",
			Type:                 "Multiplex",
			Address:              "123 Movie Street, City",
			Rows:                 10,
			Columns:              20,
			ScreenNumber:         1,
			Longitude:            12.34,
			Latitude:             56.78,
			MovieFormatSupported: pq.StringArray{"2D", "3D"},
			LanguagesSupported:   pq.StringArray{"English", "Spanish"},
			Movies: []models.Movie{
				{
					Title:           "Inception",
					Description:     "A mind-bending thriller",
					Duration:        148,
					Language:        pq.StringArray{"English"},
					Type:            pq.StringArray{"Sci-Fi", "Thriller"},
					ReleaseDate:     time.Now(),
					MovieResolution: pq.StringArray{"1080p", "4K"},
				},
			},
		}

		_, status, err := m.AddVenue(venue)

		if status != 200 {
			t.Error("status should be 200", err)
			return
		}

		if err != nil {
			t.Error("error should be nil", err)
			return
		}

	})

	t.Run("Code to migrate the database", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping this test in short mode")
		}

		err := godotenv.Load()

		if err != nil {
			t.Errorf("Could not load .env file")
		}

		m := api.NewMovieDB()

		// connect to database

		conn, err := helper.ConnectToDB()

		if err != nil {
			t.Errorf("unable to connect to database")
		}

		m.DB.Conn = conn

		// m.DB.Conn.AutoMigrate(&models.BookedSeats{}, &models.MovieTimeSlot{}, &models.CastAndCrew{}, &models.Review{}, &models.Venue{}, &models.SeatMatrix{}, &models.Movie{}, &models.User{}, &models.Idempotent{})

		m.DB.Conn.AutoMigrate(&models.Movie{})
	})
}
