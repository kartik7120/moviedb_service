package consumers

import (
	"encoding/json"
	"fmt"

	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp091.Channel
}

func NewConsumer(c *amqp091.Channel) Consumer {
	return Consumer{
		conn: c,
	}
}

func (c *Consumer) Send_Mail_Consumer() error {
	q, err := c.conn.QueueDeclare(
		"send_mail_queue2",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	dlq, err := c.conn.QueueDeclare(
		"dead_letter_queue",
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
			var msg helper.SendMailStruct
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("Unmarshal error: %v\n", err)
				d.Nack(false, true)
				continue
			}

			// Ensure headers map is initialized
			if d.Headers == nil {
				d.Headers = amqp091.Table{}
			}

			// Safely extract retry count
			retryCount := 0
			if val, ok := d.Headers["x-retry-count"]; ok {
				switch v := val.(type) {
				case int32:
					retryCount = int(v)
				case int:
					retryCount = v
				case int64:
					retryCount = int(v)
				default:
					fmt.Printf("Unexpected header type for x-retry-count: %T\n", v)
				}
			}

			fmt.Printf(string(d.Body))

			// Attempt to send mail
			if err := helper.SendMail(msg); err != nil {
				retryCount++
				if retryCount >= 3 {
					// Push to DLQ
					err = c.conn.Publish(
						"",
						dlq.Name,
						false,
						false,
						amqp091.Publishing{
							Headers:     d.Headers,
							ContentType: "application/json",
							Body:        d.Body,
						},
					)
					if err != nil {
						fmt.Printf("Failed to publish to DLQ: %v\n", err)
						d.Nack(false, false) // discard if DLQ fails
					} else {
						d.Ack(false)
					}
				} else {
					// Republish with updated retry count
					d.Headers["x-retry-count"] = int32(retryCount)
					err = c.conn.Publish(
						"",
						q.Name,
						false,
						false,
						amqp091.Publishing{
							Headers:     d.Headers,
							ContentType: "application/json",
							Body:        d.Body,
						},
					)
					if err != nil {
						fmt.Printf("Failed to republish: %v\n", err)
						d.Nack(false, false) // discard if republish fails
					} else {
						d.Ack(false)
					}
				}
			}

			// Success
			d.Ack(false)
		}
	}()
	return nil
}

// func (c *Consumer) Lock_Seat_Consumer() error {

// 	q, err := c.conn.QueueDeclare(
// 		"lock_seats_queue",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	dlq, err := c.conn.QueueDeclare(
// 		"dead_letter_queue",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	msgs, err := c.conn.Consume(
// 		q.Name,
// 		"",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				fmt.Printf("Recovered from panic: %v\n", r)
// 			}
// 		}()

// 	}()

// }
