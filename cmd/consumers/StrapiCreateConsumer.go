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
	"time"

	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
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

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	Total     int `json:"total"`
}

type UpdateCastStrapiBody struct {
	Data Cast `json:"data"`
}

func ptrInt(i int) *int {
	return &i
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

	if !exists || reflect.TypeOf(movieID) == reflect.TypeOf(float64(0)) {
		return errors.New("movie_id not found in cast and crew data")
	}

	typeOfRole, exists := cast_add_event.Data["type"]

	if !exists || reflect.TypeOf(typeOfRole) != reflect.TypeOf("") {
		return errors.New("invalid type for cast and crew")
	}

	castCrew, err := helper.ConvertAnyDataIntoCastAndCrewType(cast_add_event)

	if err != nil {
		return err
	}

	if castCrew.MovieID == 0 {
		fmt.Println("movie id cannot be zero")
		return errors.New("movie id cannot be zero")
	}

	castStrapiUID, exists := cast_add_event.Data["starpi_cast_uid"]

	if !exists || reflect.TypeOf(castStrapiUID) != reflect.TypeOf("") {
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

	resp, err := httpClient.Do(&http.Request{
		Method: "PUT",
		URL: &url.URL{
			Path: fmt.Sprintf("%s/api/casts?filters[starpi_cast_uid][$eq]=%s", strapiURL, castStrapiUID),
		},
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
			"Content-Type":  []string{"application/json"},
		},
	})

	if err != nil {
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

	var updateCastBody UpdateCastStrapiBody

	updateCastBody.Data = Cast{
		CastID:   ptrInt(int(castCrew.ID)),
		IsSynced: true,
	}

	jsonRequestBody, err := json.Marshal(updateCastBody)

	if err != nil {
		tx.Rollback()
		return err
	}

	resp, err = httpClient.Do(&http.Request{
		Method: "PUT",
		URL: &url.URL{
			Path: fmt.Sprintf("%s/api/casts/%s", strapiURL, strapiDocumentID),
		},
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", strapiToken)},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(jsonRequestBody)),
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *Consumer) Strapi_Create() error {

	q, err := c.conn.QueueDeclare(
		"strapi_create",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// 2. Declare the SAME exchange used by producer
	err = c.conn.ExchangeDeclare(
		"strapi_create_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("error declaring exchange: %v", err)
	}

	// 3. BIND queue → exchange → routing key
	err = c.conn.QueueBind(
		q.Name,                   // queue = "strapi_create"
		q.Name,                   // routing key = "strapi_create"
		"strapi_create_exchange", // exchange name
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue bind error: %v", err)
	}

	// 4. Start consuming
	msgs, err := c.conn.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()

		// used for checking if a particular message should be send to DLQ
		m := make(map[string]int)

		for d := range msgs {
			var msg EventStrapiCreate

			fmt.Printf("message received : %#v\n", d.Body)

			if err := json.Unmarshal(d.Body, &msg); err != nil {
				m[d.CorrelationId]++
				fmt.Printf("Unmarshal error: %v\n", err)
				d.Nack(false, true)

				if m[d.CorrelationId] >= 3 {
					fmt.Printf("Message moved to DLQ: %v\n", d.CorrelationId)
					// Logic to move message to DLQ can be implemented here
					d.Ack(false) // Acknowledge the message to remove it from the main queue
				}
				continue
			}

			fmt.Printf("Received message for Strapi creation: %v\n", msg)

			strapiCreateData, err := ConvertAnyToExtendedCastAndCrew(msg.Data)

			fmt.Println("CorrelationId : ", d.CorrelationId)

			if err != nil {
				m[d.CorrelationId]++
				fmt.Printf("Data conversion error: %v\n", err)
				d.Nack(false, true)

				if m[d.CorrelationId] >= 3 {
					fmt.Printf("Message moved to DLQ: %v\n", d.CorrelationId)
					// Logic to move message to DLQ can be implemented here
					d.Ack(false) // Acknowledge the message to remove it from the main queue
				}
				continue
			}

			fmt.Printf("Converted data: %+v\n", strapiCreateData)

			// Here you would add the logic to create the entry in Strapi using its API.
			// For demonstration, we'll just print the message.

			// Create entry in the database
			// If entry is successful then update the strapi model that data is synced with backend
			// Add the primary key of the model from the database to the strapi model so that when deleting the model from strapi, it can also be deleted from the database as well

			if (msg.Action != "create") || (msg.Model != "cast-and-crew") {
				m[d.CorrelationId]++
				fmt.Printf("Invalid action or model: %v, %v\n", msg.Action, msg.Model)
				d.Nack(false, true)

				if m[d.CorrelationId] >= 3 {
					fmt.Printf("Message moved to DLQ: %v\n", d.CorrelationId)
					// Logic to move message to DLQ can be implemented here
					d.Ack(false) // Acknowledge the message to remove it from the main queue
				}
				continue
			}

			err = c.AddCastStrapiEntry(msg)

			if err != nil {
				m[d.CorrelationId]++
				fmt.Printf("Error adding Strapi entry: %v\n", err)
				d.Nack(false, true)

				if m[d.CorrelationId] >= 3 {
					fmt.Printf("Message moved to DLQ: %v\n", d.CorrelationId)
					// Logic to move message to DLQ can be implemented here
					d.Ack(false) // Acknowledge the message to remove it from the main queue
				}
				continue
			}
			// Acknowledge the message after processing
			d.Ack(false)
		}
	}()

	return nil
}
