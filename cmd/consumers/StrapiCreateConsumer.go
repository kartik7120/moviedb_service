package consumers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/kartik7120/booking_moviedb_service/cmd/models"
	"gorm.io/gorm/clause"
)

type EventStrapiCreate struct {
	Data   map[string]interface{} `json:"data"`
	Action string                 `json:"action"`
	Model  string                 `json:"model"`
}

type ExtendedCastAndCrew struct {
	models.CastAndCrew
	StarpiCastUid string `json:"strapi_cast_uid"`
}
type CastUpdateBody struct {
	Data struct {
		CastID          string `json:"cast_id,omitempty"`
		IsSynced        bool   `json:"is_synced"`
		Starpi_cast_uid string `json:"starpi_cast_uid,omitempty"`
	} `json:"data"`
}

type UpdateCastStrapiBody struct {
	Data CastUpdateBody `json:"data"`
}

func ConvertAnyDataIntoCastAndCrewType(event EventStrapiCreate) (models.CastAndCrew, error) {
	var castCrew models.CastAndCrew

	// Use ONLY event.Data, not entire event
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return castCrew, err
	}

	// Fix MovieID mismatch by mapping before unmarshaling
	// because JSON key is "MovieID" but your struct expects "movie_id"
	var temp map[string]interface{}
	json.Unmarshal(dataBytes, &temp)

	if val, ok := temp["MovieID"]; ok {
		temp["movie_id"] = val // map to expected field
		delete(temp, "MovieID")
	}

	if val, ok := temp["cast_id"]; ok {
		// convert to float64 first
		switch v := val.(type) {
		case float64:
			temp["cast_id"] = int(v)
		case int:
			// already int, do nothing
		default:
			return castCrew, errors.New("invalid type for cast_id")
		}
	}

	temp["id"] = temp["cast_id"]

	// marshal updated map
	dataBytes, _ = json.Marshal(temp)

	// final unmarshal
	err = json.Unmarshal(dataBytes, &castCrew)
	if err != nil {
		return castCrew, err
	}

	return castCrew, nil
}

func ConvertAnyToExtendedCastAndCrew(data any) (ExtendedCastAndCrew, error) {
	var extendedCast ExtendedCastAndCrew

	// First, marshal the 'data' back to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return extendedCast, err
	}

	// Then, unmarshal the JSON into the ExtendedCastAndCrew struct
	err = json.Unmarshal(jsonData, &extendedCast)
	if err != nil {
		return extendedCast, err
	}

	return extendedCast, nil
}

// TODO: Need to correctly map the fields from Strapi to our model
func ConvertAnyToMovieTimeSlotStrapi(data any) (MovieTimeSlotStrapi, error) {
	var mts MovieTimeSlotStrapi

	raw, err := json.Marshal(data)
	if err != nil {
		return mts, err
	}

	if err := json.Unmarshal(raw, &mts); err != nil {
		return mts, err
	}

	// Validation
	if mts.MovieID == 0 {
		return mts, errors.New("movie_id missing in movietimeslot")
	}
	if mts.VenueID == 0 {
		return mts, errors.New("venue_id missing in movietimeslot")
	}
	if mts.StarpiMovieTimeSlotUID == "" {
		return mts, errors.New("strapi movie time slot uid missing")
	}

	return mts, nil
}

func MapMovieTimeSlotStrapiToModel(
	s MovieTimeSlotStrapi,
) (models.MovieTimeSlot, error) {

	startTime, err := time.Parse(time.RFC3339, s.StartTime)
	if err != nil {
		return models.MovieTimeSlot{}, err
	}

	endTime, err := time.Parse(time.RFC3339, s.EndTime)
	if err != nil {
		return models.MovieTimeSlot{}, err
	}

	date, err := time.Parse(time.RFC3339, s.Date)
	if err != nil {
		return models.MovieTimeSlot{}, err
	}

	return models.MovieTimeSlot{
		MovieID:     uint(s.MovieID),
		VenueID:     uint(s.VenueID),
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    int(s.Duration),
		MovieFormat: s.MovieFormat,
	}, nil
}

type CastResponse struct {
	Data []Cast `json:"data"`
	Meta Meta   `json:"meta"`
}

type Cast struct {
	ID            int       `json:"id"`
	DocumentID    string    `json:"documentId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PublishedAt   time.Time `json:"publishedAt"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	CharacterName string    `json:"character_name"`
	PhotoURL      *string   `json:"photo_url"` // null → pointer
	IsSynced      bool      `json:"is_synced"`
	StrapiCastUID string    `json:"starpi_cast_uid"`
	CastID        *int      `json:"cast_id"` // null → pointer
}

type StrapiMovie struct {
	MovieResolution  string    `json:"movieResolution"`
	CreatedAt        time.Time `json:"createdAt"`
	Description      *string   `json:"description,omitempty"`
	DocumentID       string    `json:"documentId"`
	Duration         int64     `json:"duration"`
	ID               int       `json:"id"`
	IsSynced         bool      `json:"is_synced"`
	Language         *string   `json:"language,omitempty"`
	Languages        string    `json:"languages"`
	LogoImageURL     *string   `json:"logoImageURL,omitempty"`
	MovieID          *string   `json:"movieid,omitempty"`
	PosterURL        string    `json:"posterURL"`
	PublishedAt      time.Time `json:"publishedAt"`
	Ranking          int       `json:"ranking"`
	ReleaseDate      time.Time `json:"releaseDate"`
	ScreenWidePoster *string   `json:"screenWidePoster,omitempty"`
	Title            string    `json:"title"`
	TrailerURL       *string   `json:"trailerURL,omitempty"`
	Type             string    `json:"type"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Votes            int       `json:"votes"`
}

type StrapiVenue struct {
	CreatedAt   time.Time `json:"createdAt"`
	DocumentID  string    `json:"documentId"`
	ID          int       `json:"id"`
	IsSynced    bool      `json:"is_synced"`
	PublishedAt time.Time `json:"publishedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	VenueID     int       `json:"venueID"`
}

type MovieTimeSlotStrapi struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Date      string `json:"date"`

	Duration    int64  `json:"duration"`
	MovieFormat string `json:"movie_format"`

	MovieID int `json:"movie_id"`
	VenueID int `json:"venue_id"`

	IsSynced bool `json:"is_synced"`

	StarpiMovieTimeSlotUID string `json:"strapi_movie_uid"`
}

type MovieTimeSlotStrapiDataResponse struct {
	ID         int    `json:"id"`
	DocumentID string `json:"documentId"`

	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	PublishedAt string `json:"publishedAt"`

	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Date      string `json:"date"`

	Duration    int64  `json:"duration"`
	MovieFormat string `json:"movieformat"`

	IsSynced bool `json:"is_synced"`

	StrapiUID string `json:"starpi_movie_time_slot_uid"`

	MovieTimeSlotID *int `json:"movie_time_slot_id"`
}

type MovieTimeSlotStrapiResponse struct {
	Data []MovieTimeSlotStrapiDataResponse `json:"data"`
	Meta Meta                              `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	Total     int `json:"total"`
}

func (c *Consumer) AddCastStrapiEntry(cast_add_event EventStrapiCreate) error {
	var castCrew models.CastAndCrew

	strapiURL := os.Getenv("STRAPI_URL")
	strapiToken := os.Getenv("STRAPI_API_TOKEN")

	if strapiURL == "" || strapiToken == "" {
		return errors.New("strapi url or token not set in environment variables")
	}

	// convert cast_add_event.data any type to models.CastAndCrew

	if reflect.TypeOf(cast_add_event.Data) != reflect.TypeOf(map[string]interface{}{}) {
		return errors.New("invalid data type for cast and crew")
	}

	name, exists := cast_add_event.Data["name"]
	if !exists || reflect.TypeOf(name) != reflect.TypeOf("") {
		return errors.New("invalid name type for cast and crew")
	}

	movieID, exists := cast_add_event.Data["movie_id"]
	if !exists {
		return errors.New("movie_id not found in cast and crew data")
	}

	movieIDFloat, ok := movieID.(float64)
	if !ok {
		return errors.New("movie_id must be a number")
	}

	typeOfRole, exists := cast_add_event.Data["type"]
	if !exists || reflect.TypeOf(typeOfRole) != reflect.TypeOf("") {
		return errors.New("invalid type for cast and crew")
	}

	// Convert to struct
	castCrew, err := ConvertAnyDataIntoCastAndCrewType(cast_add_event)
	if err != nil {
		return err
	}

	// Apply correct MovieID
	castCrew.MovieID = uint(movieIDFloat)

	if castCrew.MovieID == 0 {
		return errors.New("movie id cannot be zero")
	}

	castStrapiUID, exists := cast_add_event.Data["strapi_cast_uid"]

	if !exists || reflect.ValueOf(castStrapiUID) == reflect.ValueOf("") {
		return errors.New("invalid strapi_cast_uid type for cast and crew")
	}

	tx := c.DB.Conn.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&castCrew)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// add the cast and crew id in the strapi model using the strapi api

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	u, _ := url.Parse(strapiURL)

	u.Path = "/api/casts"
	q := u.Query()
	q.Set("filters[starpi_cast_uid][$eq]", castStrapiUID.(string))
	u.RawQuery = q.Encode()

	resp, err := httpClient.Do(&http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
			"Content-Type":  []string{"application/json"},
		},
	})

	if err != nil {
		fmt.Println("error calling the strapi put endpoint : ", err.Error())
		tx.Rollback()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tx.Rollback()
		return fmt.Errorf("failed to update strapi cast and crew with id %s, status code: %d", castStrapiUID, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		tx.Rollback()
		return err
	}

	fmt.Println(string(bodyBytes))

	var strapiCastResponse CastResponse

	err = json.Unmarshal(bodyBytes, &strapiCastResponse)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit error: %v", err)
	}

	var strapiDocumentID string

	for _, v := range strapiCastResponse.Data {
		strapiDocumentID = v.DocumentID
	}

	updateBody := CastUpdateBody{}
	updateBody.Data.CastID = fmt.Sprintf("%d", castCrew.ID) // convert to string
	updateBody.Data.IsSynced = true

	jsonRequestBody, err := json.Marshal(updateBody)

	if err != nil {
		tx.Rollback()
		return err
	}

	fmt.Println("document id : ", strapiDocumentID)

	u2, _ := url.Parse(strapiURL)
	u2.Path = fmt.Sprintf("/api/casts/%s", strapiDocumentID)

	resp, err = httpClient.Do(&http.Request{
		Method: "PUT",
		URL:    u2,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(jsonRequestBody)),
	})

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Strapi Response:", string(respBody))

	if err != nil {
		fmt.Println("error updating the document : ", err)
		tx.Rollback()
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error updating the document : ", string(respBody))
		tx.Rollback()
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *Consumer) DeleteCastStrapiEntry(cast_delete_event EventStrapiCreate) error {
	// Implementation for deleting a cast entry from Strapi can be added here
	var castCrew models.CastAndCrew

	strapiURL := os.Getenv("STRAPI_URL")
	strapiToken := os.Getenv("STRAPI_API_TOKEN")

	if strapiURL == "" || strapiToken == "" {
		return errors.New("strapi url or token not set in environment variables")
	}

	// convert cast_add_event.data any type to models.CastAndCrew

	if reflect.TypeOf(cast_delete_event.Data) != reflect.TypeOf(map[string]interface{}{}) {
		return errors.New("invalid data type for cast and crew")
	}

	castCrew, err := ConvertAnyDataIntoCastAndCrewType(cast_delete_event)

	if err != nil {
		return err
	}

	if castCrew.MovieID == 0 {
		fmt.Println("movie id cannot be zero")
		return errors.New("movie id cannot be zero")
	}

	castStrapiUID, exists := cast_delete_event.Data["strapi_cast_uid"]

	if !exists || reflect.ValueOf(castStrapiUID) == reflect.ValueOf("") {
		return errors.New("invalid strapi_cast_uid type for cast and crew")
	}

	tx := c.DB.Conn.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := c.DB.Conn.Unscoped().Where("id = ?", castCrew.ID).Delete(&models.CastAndCrew{})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// add the cast and crew id in the strapi model using the strapi api

	// httpClient := http.Client{
	// 	Timeout: time.Second * 10,
	// }

	// u, _ := url.Parse(strapiURL)

	// u.Path = "/api/casts"
	// q := u.Query()
	// q.Set("filters[starpi_cast_uid][$eq]", castStrapiUID.(string))
	// u.RawQuery = q.Encode()

	// resp, err := httpClient.Do(&http.Request{
	// 	Method: "GET",
	// 	URL:    u,
	// 	Header: http.Header{
	// 		"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
	// 		"Content-Type":  []string{"application/json"},
	// 	},
	// })

	// if err != nil {
	// 	fmt.Println("error calling the strapi put endpoint : ", err.Error())
	// 	tx.Rollback()
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	tx.Rollback()
	// 	return fmt.Errorf("failed to update strapi cast and crew with id %s, status code: %d", castStrapiUID, resp.StatusCode)
	// }

	// bodyBytes, err := io.ReadAll(resp.Body)

	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// fmt.Println(string(bodyBytes))

	// var strapiCastResponse CastResponse

	// err = json.Unmarshal(bodyBytes, &strapiCastResponse)

	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// if err := tx.Commit().Error; err != nil {
	// 	return fmt.Errorf("commit error: %v", err)
	// }

	// var strapiDocumentID string

	// for _, v := range strapiCastResponse.Data {
	// 	strapiDocumentID = v.DocumentID
	// }

	// updateBody := CastUpdateBody{}
	// updateBody.Data.CastID = fmt.Sprintf("%d", castCrew.ID) // convert to string
	// updateBody.Data.IsSynced = true
	// updateBody.Data.Starpi_cast_uid = castStrapiUID.(string)

	// jsonRequestBody, err := json.Marshal(updateBody)

	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// fmt.Println("document id : ", strapiDocumentID)

	// u2, _ := url.Parse(strapiURL)
	// u2.Path = fmt.Sprintf("/api/casts/%s", strapiDocumentID)

	// resp, err = httpClient.Do(&http.Request{
	// 	Method: "DELETE",
	// 	URL:    u2,
	// 	Header: http.Header{
	// 		"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
	// 		"Content-Type":  []string{"application/json"},
	// 	},
	// 	Body: io.NopCloser(bytes.NewReader(jsonRequestBody)),
	// })

	// respBody, _ := io.ReadAll(resp.Body)
	// fmt.Println("Strapi Response:", string(respBody))

	// if err != nil {
	// 	fmt.Println("error updating the document : ", err)
	// 	tx.Rollback()
	// 	return err
	// }

	// if resp.StatusCode != http.StatusNoContent {
	// 	fmt.Println("error updating the document : ", string(respBody))
	// 	tx.Rollback()
	// 	return err
	// }

	// defer resp.Body.Close()

	fmt.Println("Successfully deleted cast and crew with ID:", castCrew.ID)

	return nil

}

func (c *Consumer) AddMovieTimeSlot(movie_time_slot_add_event EventStrapiCreate) error {

	// var movietimeslot models.MovieTimeSlot

	strapiURL := os.Getenv("STRAPI_URL")
	strapiToken := os.Getenv("STRAPI_API_TOKEN")

	if strapiURL == "" || strapiToken == "" {
		return errors.New("strapi url or token not set in environment variables")
	}

	if reflect.TypeOf(movie_time_slot_add_event.Data) != reflect.TypeOf(map[string]interface{}{}) {
		return errors.New("invalid data type for movie time slot")
	}

	fmt.Printf("event data received %+v", movie_time_slot_add_event)

	_, isExists := movie_time_slot_add_event.Data["movie_id"].(float64)

	if !isExists {
		fmt.Printf("movie id does not exists")
		return errors.New("movie id cannot be empty")
	}

	_, isExists = movie_time_slot_add_event.Data["venue_id"].(float64)

	if !isExists {
		fmt.Printf("venue id does not exists")
		return errors.New("venue id cannot be empty")
	}

	strapiID, isExists := movie_time_slot_add_event.Data["strapi_movie_uid"].(string)

	if !isExists {
		fmt.Printf("strapi movie timeslot uid does not exists")
		return errors.New("strapi movie timeslot uid cannot be empty")
	}

	tx := c.DB.Conn.Begin()

	if tx.Error != nil {
		fmt.Printf("error creating a transaction : %s", tx.Error.Error())
		return tx.Error
	}

	defer func() {
		fmt.Println("panic occured while in movie time slot consumer")
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Printf("recover error in movie time slot function: %+v", r)
		}
	}()

	movieTimeSlotStrapiType, err := ConvertAnyToMovieTimeSlotStrapi(movie_time_slot_add_event.Data)

	if err != nil {
		fmt.Printf("error converting to movie time slot strapi type : %s", err.Error())
		return err
	}

	fmt.Printf("movie time slot strapi type %+v", movieTimeSlotStrapiType)

	movietimeslot, err := MapMovieTimeSlotStrapiToModel(movieTimeSlotStrapiType)

	if err != nil {
		fmt.Printf("error mapping to model movie time slot : %s", err.Error())
		return err
	}

	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&movietimeslot)

	if result.Error != nil {
		fmt.Printf("error creating movie time slot in db : %s", result.Error.Error())
		tx.Rollback()
		return result.Error
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("error committing the transaction : %s", err.Error())
		return fmt.Errorf("commit error: %v", err)
	}

	u, err := url.Parse(strapiURL)

	if err != nil {
		fmt.Println("error parsing the strapi url : ", err.Error())
		return err
	}

	u.Path = "/api/movietimeslots"
	q := u.Query()
	q.Set("filters[starpi_movie_time_slot_uid][%24eq]", strapiID)
	u.RawQuery = q.Encode()

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Do(
		&http.Request{
			Method: "GET",
			URL:    u,
			Header: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
				"Content-Type":  []string{"application/json"},
			},
		},
	)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error occured when calling the document strapi call : %s", err.Error())
		return fmt.Errorf("failed to fetch strapi movie time slot with id %s, status code: %d", strapiID, resp.StatusCode)
	}

	if err != nil {
		fmt.Printf("error calling the strapi endpoint : %s", err.Error())
		fmt.Println("error calling the strapi get endpoint : ", err.Error())
		return err
	}

	respBody, err := io.ReadAll(resp.Body)

	fmt.Println("Strapi Response:", string(respBody))

	if err != nil {
		fmt.Printf("error reading response body : %s", err.Error())
		fmt.Println("error fetching the document : ", err)
		return err
	}

	var movieTimeSlotStrapiResponse MovieTimeSlotStrapiResponse

	err = json.Unmarshal(respBody, &movieTimeSlotStrapiResponse)

	if err != nil {
		fmt.Printf("error unmarshalling the response body : %s", err.Error())
		return err
	}

	documentID := ""

	for _, v := range movieTimeSlotStrapiResponse.Data {
		documentID = v.DocumentID
	}

	updateBody := struct {
		Data struct {
			MovieTimeSlotID string `json:"movie_time_slot_id,omitempty"`
			IsSynced        bool   `json:"is_synced"`
		} `json:"data"`
	}{}

	updateBody.Data.MovieTimeSlotID = fmt.Sprintf("%d", movietimeslot.ID)
	updateBody.Data.IsSynced = true

	jsonRequestBody, err := json.Marshal(updateBody)

	if err != nil {
		return err
	}

	u2, _ := url.Parse(strapiURL)
	u2.Path = fmt.Sprintf("/api/movietimeslots/%s", documentID)

	resp, err = httpClient.Do(&http.Request{
		Method: "PUT",
		URL:    u2,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(jsonRequestBody)),
	})

	respBody, _ = io.ReadAll(resp.Body)
	fmt.Println("Strapi Response:", string(respBody))

	if err != nil {
		fmt.Println("error calling the strapi put endpoint : ", err.Error())
		fmt.Println("error updating the document : ", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error updating the document : ", string(respBody))
		return err
	}

	fmt.Println("Successfully added movie time slot with ID:", movietimeslot.ID, " Strapi UID: ", strapiID)

	return nil
}

func (c *Consumer) Strapi_Create() error {

	// Declare queue
	q, err := c.conn.QueueDeclare(
		"strapi_create",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	// Declare exchange
	// err = c.conn.ExchangeDeclare(
	// 	"strapi_create_exchange",
	// 	"direct",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("error declaring exchange: %v", err)
	// }

	// // Bind routing key
	// err = c.conn.QueueBind(
	// 	q.Name,
	// 	q.Name,
	// 	"strapi_create_exchange",
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("queue bind error: %v", err)
	// }

	// Start consumer
	msgs, err := c.conn.Consume(
		q.Name,
		"",
		false, // manual ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Retry tracking map (per worker, not per message)
	retryCounter := sync.Map{} // map[string]int

	go func() {
		for d := range msgs {

			deliveryKey := d.MessageId

			// Increase retry count
			val, _ := retryCounter.LoadOrStore(deliveryKey, 0)
			retry := val.(int)

			fmt.Printf("\n---- New Message [%d] ----\n", d.DeliveryTag)
			fmt.Printf("Raw body: %s\n", string(d.Body))
			fmt.Printf("Retry count: %d\n", retry)

			// Parse message
			var msg EventStrapiCreate
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("Unmarshal error: %v\n", err)
				retryCounter.Store(deliveryKey, retry+1)

				// Send to DLQ after 3 failures
				if retry+1 >= 3 {
					fmt.Printf("Sending message to DLQ (unmarshal failed)\n")
					d.Nack(false, false) // requeue = false → DLQ
					retryCounter.Delete(deliveryKey)
					continue
				}

				d.Nack(false, true) // retry
				continue
			}

			// Validate action/model
			if msg.Model != "cast-and-crew" && msg.Model != "movie-time-slot" {
				fmt.Printf("Invalid event type: %v / %v\n", msg.Action, msg.Model)

				retryCounter.Store(deliveryKey, retry+1)

				if retry+1 >= 3 {
					fmt.Printf("Sending message to DLQ (invalid event)\n")
					d.Nack(false, false)
					retryCounter.Delete(deliveryKey)
					continue
				}

				d.Nack(false, true)
				continue
			}

			if msg.Model == "cast-and-crew" {

				// Convert message → struct
				castData, err := ConvertAnyToExtendedCastAndCrew(msg.Data)

				if err != nil {
					fmt.Printf("Data conversion error: %v\n", err)
					retryCounter.Store(deliveryKey, retry+1)

					if retry+1 >= 3 {
						fmt.Printf("Sending message to DLQ (convert error)\n")
						d.Nack(false, false)
						retryCounter.Delete(deliveryKey)
						continue
					}

					d.Nack(false, true)
					continue
				}

				fmt.Printf("Converted data: %+v\n", castData)

				// Process

				switch msg.Action {
				case "create":
					fmt.Printf("Processing cast creation for: %s\n", castData.Name)
					err = c.AddCastStrapiEntry(msg)
				case "delete":
					fmt.Printf("Processing cast deletion for: %s\n", castData.Name)
					err = c.DeleteCastStrapiEntry(msg)
				default:
					fmt.Printf("Unknown action: %s\n", msg.Action)
					err = fmt.Errorf("unknown action: %s", msg.Action)
				}
			} else if msg.Model == "movie-time-slot" {

				mtsStrapi, err := ConvertAnyToMovieTimeSlotStrapi(msg.Data)

				if err != nil {
					fmt.Printf("MovieTimeSlot conversion error: %v\n", err)
					retryCounter.Store(deliveryKey, retry+1)

					if retry+1 >= 3 {
						d.Nack(false, false)
						retryCounter.Delete(deliveryKey)
						continue
					}

					d.Nack(false, true)
					continue
				}

				fmt.Printf("Converted MovieTimeSlot Strapi: %+v\n", mtsStrapi)

				switch msg.Action {
				case "create":
					err = c.AddMovieTimeSlot(msg)
				case "delete":
					// future: c.DeleteMovieTimeSlot(msg)
				default:
					err = fmt.Errorf("unknown action: %s", msg.Action)
				}

			}

			if err != nil {
				fmt.Printf("Error deleting Strapi entry: %v\n", err)
				retryCounter.Store(deliveryKey, retry+1)

				if retry+1 >= 3 {
					fmt.Printf("Sending message to DLQ (business logic failed)\n")
					d.Nack(false, false)
					retryCounter.Delete(deliveryKey)
					continue
				}

				d.Nack(false, true)
				continue
			}

			// Success
			fmt.Printf("Successfully processed message.\n")
			d.Ack(false)
			retryCounter.Delete(deliveryKey)
		}
	}()

	return nil
}
