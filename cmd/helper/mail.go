// filepath: /home/kartik7120/booking_auth_service/cmd/helper/mail.go
package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// If this mail implementation is not working for us then we can also host mailhog on our local machine

func SendMail(to string, name string, text string, html string, category string, subject string) error {

	validate = validator.New()

	// validating if the email is valid
	log.Info("Starting to send email")

	err := validate.Var(to, "email")

	if err != nil {
		log.Error("Invalid email", err)
		return err
	}

	if len(text) == 0 && len(html) == 0 {
		log.Error("text or html is required")
		return fmt.Errorf("text or html is required")
	}

	// check if text or html does not contain any malicious code

	if len(text) > 0 {
		err = validate.Var(text, "printascii")

		if err != nil {
			log.Error("Invalid text", err)
			return err
		}
	}

	if len(html) > 0 {

		err = validate.Var(html, "printascii")

		if err != nil {
			log.Error("Invalid html")
			return err
		}

		err = validate.Var(html, "html")

		if err != nil {
			log.Error("Invalid html", err)
			return err
		}

	}

	if len(subject) > 0 {
		err = validate.Var(subject, "printascii")

		if err != nil {
			log.Error("Invalid subject")
			return err
		}
	}

	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	payloadString := ""

	if len(html) > 0 {
		payloadString = fmt.Sprintf(`{
            "from": {
                "email": "hello@demomailtrap.co",
                "name": %s
            },
            "to": [
                {
                    "email": %s
                }
            ],
            "subject": %s,
            "html": %s,
            "category": %s
        }`, jsonEscape(name), jsonEscape(to), jsonEscape(subject), jsonEscape(html), jsonEscape(category))
	} else {
		payloadString = fmt.Sprintf(`{
                "from": {
                    "email": "hello@demomailtrap.co",
                    "name": %s
                },
                "to": [
                    {
                        "email": %s
                    }
                ],
                "subject": %s,
                "text": %s,
                "category": %s
            }`, jsonEscape(name), jsonEscape(to), jsonEscape(subject), jsonEscape(text), jsonEscape(category))
	}

	log.Info(payloadString)

	payload := strings.NewReader(payloadString)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Error("Error creating request")
		return err
	}

	tokenString := fmt.Sprintf("Bearer %s", os.Getenv("MAILTRAP_API_KEY"))

	req.Header.Add("Authorization", tokenString)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		log.Error("Error sending email")
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		log.Error("Error reading response")
		return err
	}

	log.Info(string(body))

	// check if body is valid to unmarshal

	if len(body) == 0 {
		log.Error("Empty response", err)
		return errors.New("empty response")
	}

	if len(body) > 0 {
		if body[0] != byte('{') {
			log.Error("Invalid response", err)
			return errors.New("invalid response")
		}
	}

	// convert body to json

	jsonBody := make(map[string]any)

	err = json.Unmarshal(body, &jsonBody)

	if err != nil {
		log.Error("Error unmarshalling response", err)
		return err
	}

	if jsonBody["success"] == false {
		log.Error("Error sending email", jsonBody["errors"])
		return errors.New("error sending email")
	}
	// check if the email was sent successfully

	log.Info("Email sent successfully")
	return nil
}

func jsonEscape(value string) string {
	escapedValue, _ := json.Marshal(value)
	return string(escapedValue)
}
