package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	LOGIN_URL        = "https://onlinebusiness.icbc.com/deas-api/v1/webLogin/webLogin"
	APPOINTMENTS_URL = "https://onlinebusiness.icbc.com/deas-api/v1/web/getAvailableAppointments"
)

type LoginPayload struct {
	DrvrLastName  string `json:"drvrLastName"`
	LicenceNumber string `json:"licenceNumber"`
	Keyword       string `json:"keyword"`
}

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

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	userAgent := os.Getenv("APP_USER_AGENT")
	aPosIDStr := os.Getenv("APP_APPOINTMENT_POSITION_ID")
	lastName := os.Getenv("APP_LAST_NAME")
	licenseNumber := os.Getenv("APP_LICENSE_NUMBER")
	keyword := os.Getenv("APP_KEYWORD")

	aPosID, err := strconv.ParseInt(aPosIDStr, 10, 32)
	if err != nil {
		panic("aPosID is incorecrly set: " + err.Error())
	}

	bearerToken, err := getBearerToken(lastName, licenseNumber, keyword, userAgent)
	if err != nil {
		fmt.Println("Login error: " + err.Error())
	}

	appointments, err := getAvailableAppointments(int32(aPosID), lastName, licenseNumber, bearerToken, userAgent)
	if err != nil {
		fmt.Println("Get appointments error: " + err.Error())
	}

	res, err := PrettyStruct(appointments)
	if err != nil {
		fmt.Println("PrettyStruct didn't work: " + err.Error())
	}

	fmt.Println(res)
}

func getBearerToken(lastName, licenceNumber, keyword, userAgent string) (string, error) {
	payload := LoginPayload{
		DrvrLastName:  lastName,
		LicenceNumber: licenceNumber,
		Keyword:       keyword,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", LOGIN_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", "https://onlinebusiness.icbc.com/webdeas-ui/login;type=driver")
	req.Header.Add("Cache-Control", "no-cache, no-store")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unsuccessful request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("invalid credentials")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s, status code: %d", response.Status, response.StatusCode)
	}

	return response.Header.Get("Authorization"), nil
}

func getAvailableAppointments(aPosID int32, lastName, licenseNumber, bearer, userAgent string) ([]Appointment, error) {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	payload := RequestPayload{
		APosID:            aPosID,
		ExamType:          "5-R-1",
		ExamData:          tomorrow,
		IgnoreReserveTime: false,
		PrfDaysOfWeek:     "[0,1,2,3,4,5,6]",
		PrfPartsOfDay:     "[0,1]",
		LastName:          lastName,
		LicenseNumber:     licenseNumber,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", APPOINTMENTS_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("User-Agent", userAgent)

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

func PrettyStruct(data any) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
