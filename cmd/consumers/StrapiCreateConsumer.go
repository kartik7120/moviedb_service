package consumers

import (
	"encoding/json"
	"fmt"
)

type EventStrapiCreate struct {
	Data   map[string]interface{} `json:"data"`
	Action string                 `json:"action"`
	Model  string                 `json:"model"`
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

	msgs, err := c.conn.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()

		for d := range msgs {
			var msg EventStrapiCreate
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("Unmarshal error: %v\n", err)
				d.Nack(false, true)
				continue
			}

			fmt.Printf("Received message for Strapi creation: %v\n", msg)

			// Here you would add the logic to create the entry in Strapi using its API.
			// For demonstration, we'll just print the message.

			// Create entry in the database
			// If entry is successful then update the strapi model that data is synced with backend
			// Add the primary key of the model from the database to the strapi model so that when deleting the model from strapi, it can also be deleted from the database as well

			// Acknowledge the message after processing
			d.Ack(false)
		}
	}()

	return nil
}
