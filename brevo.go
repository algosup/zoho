package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type Attributes struct {
	LastName   string `json:"LASTNAME"`
	FirstName  string `json:"FIRSTNAME"`
	Pipeline   string `json:"PIPELINE"`
	Type       string `json:"TYPE"`
	LeadSource string `json:"LEADSOURCE"`
	MaxScore   int    `json:"MAXSCORE"`
	SMS        string `json:"SMS"`
	WhatsApp   string `json:"WHATSAPP"`
}

// ContactData struct representing the entire JSON structure
type BrevoContact struct {
	Email         string     `json:"email"`
	UpdateEnabled bool       `json:"updateEnabled"`
	Attributes    Attributes `json:"attributes"`
	ListIds       []int      `json:"listIds"`
}

func CreateBrevoContact(data BrevoContact) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/contacts", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("api-key", os.Getenv("BREVO_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		log.Println(resp.StatusCode)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}
