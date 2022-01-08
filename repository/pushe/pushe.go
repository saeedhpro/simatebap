package pushe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/simateb-project/simateb-backend/repository/env"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"io/ioutil"
	"net/http"
)

type Notification struct {
	ID               int64  `json:"id"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	UserID           string `json:"user_id"`
	UserFName        string `json:"user_fname"`
	UserLName        string `json:"user_lname"`
	StaffID          string `json:"staff_id"`
	StaffFName       string `json:"staff_fname"`
	StaffLName       string `json:"staff_lname"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Ids              string `json:"ids"`
	ActionUrl        string `json:"action_url"`
	CloseOnClick     string `json:"close_on_click"`
	Content          string `json:"content"`
	Type             string `json:"type"`
	CreatedAt        string `json:"created_at"`
}

type SendNotificationRequest struct {
	Title        string `json:"title"`
	Body         string `json:"body"`
	UserID       string `json:"user_id"`
	Ids          string `json:"ids"`
	ActionUrl    string `json:"action_url"`
	CloseOnClick string `json:"close_on_click"`
	Content      string `json:"content"`
	Type         string `json:"type"`
}

func SendNotification(staffUser *auth.UserClaims, request SendNotificationRequest) (bool, error) {
	url := "https://api.pushe.co/v2/messaging/notifications/"
	client := &http.Client{}
	PusheApiKey := env.GetDotEnvVariable("PUSHE_API_KEY")
	AppIds := env.GetDotEnvVariable("PUSH_APP_IDS")
	payload := map[string]interface{}{
		"app_ids" : AppIds,
		"data": map[string]string{
			"title": request.Title,
			"content": request.Content,
		},
	}

	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token " + PusheApiKey)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return true, err
}
