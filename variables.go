package zoho

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type VariableGroup struct {
	APIName string `json:"api_name"`
	ID      string `json:"id"`
}

type Variable struct {
	ReadOnly      bool          `json:"read_only"`
	APIName       string        `json:"api_name"`
	Name          string        `json:"name"`
	Description   *string       `json:"description"` // Using a pointer to handle null values
	ID            string        `json:"id"`
	Source        string        `json:"source"`
	Type          string        `json:"type"`
	VariableGroup VariableGroup `json:"variable_group"`
	Value         string        `json:"value"`
}

type Response struct {
	Variables []Variable `json:"variables"`
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetLastAutoZoho() (string, error) {

	url := "https://www.zohoapis.eu/crm/v7/settings/variables/477339000035617041"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Variables[0].Value, nil
}

func SetLastAutoZoho(date string) error {
	url := "https://www.zohoapis.eu/crm/v7/settings/variables/477339000035617041"

	json := fmt.Sprintf(`{
		"variables": [
		   {
				"id": "477339000035617029",
				"value": "%s"
			}]
	}`, date)

	req, err := http.NewRequest("PUT", url, strings.NewReader(json))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Println(resp.StatusCode)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}
