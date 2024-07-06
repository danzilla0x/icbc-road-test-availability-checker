package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Request the road test availability in Port Coquitlam
// POST request
//  URL: https://onlinebusiness.icbc.com/deas-api/v1/web/getAvailableAppointments
// 	Data: in example
// Headers:
//	Authorization: Bearer
//  User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36

const (
	URL          = "https://onlinebusiness.icbc.com/deas-api/v1/web/getAvailableAppointments"
	BEARER_TOKEN = ""
	USER_AGENT   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

type RequestPayload struct {
	APosID            int32  `json:"aPosID"`
	ExamType          string `json:"examType"`
	ExamData          string `json:"examDate"`
	IgnoreReserveTime bool   `json:"ignoreReserveTime"`
	PrfDaysOfWeek     string `json:"prfDaysOfWeek"`
	PrfPartsOfDay     string `json:"prfPartsOfDay"`
	LastName          string `json:"lastName"`
	LicenseNumber     string `json:"licenseNumber"`
}

type AppointmentDate struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"dayOfWeek"`
}

type DetailExam struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Appointment struct {
	AppointmentDt AppointmentDate `json:"appointmentDt"`
	DlExam        DetailExam      `json:"dlExam"`
	LemgMsgId     int32           `json:"lemgMsgId"`
	PosId         int32           `json:"posId"`
	ResourceId    int32           `json:"resourceId"`
	Signature     string          `json:"signature"`
	StartTm       string          `json:"startTm"`
	EndTm         string          `json:"endTm"`
}

func PrettyStruct(data any) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func getAvailableAppointments(bearer string) ([]Appointment, error) {

	payload := RequestPayload{
		APosID:            0,
		ExamType:          "5-R-1",
		ExamData:          "2024-07-06", // TODO update date to tommorow
		IgnoreReserveTime: false,
		PrfDaysOfWeek:     "[0,1,2,3,4,5,6]",
		PrfPartsOfDay:     "[0,1]",
		LastName:          "",
		LicenseNumber:     "",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	// add authorization header to the req
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("User-Agent", USER_AGENT)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid credentials")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var appointments []Appointment
	err = json.NewDecoder(response.Body).Decode(&appointments)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful JSON decode: %w", err)
	}

	return appointments, nil
}

func main() {

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + BEARER_TOKEN // TODO migh be provided in the runtime

	appointments, err := getAvailableAppointments(bearer)
	if err != nil {
		fmt.Println("Internal error: " + err.Error())
	}

	res, err := PrettyStruct(appointments)
	if err != nil {
		fmt.Println("PrettyStruct didn't work: " + err.Error())
	}

	fmt.Println(res)
}
