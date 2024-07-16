package appointment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

func GetAvailableAppointments(appointmentUrl string, aPosID int32, lastName, licenseNumber, bearer, userAgent string) ([]Appointment, error) {
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

	req, err := http.NewRequest("POST", appointmentUrl, bytes.NewBuffer(jsonData))
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
		return nil, fmt.Errorf("invalid authentication credentials")
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authentication has expired")
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

func FindExamAppointment(appointments []Appointment, startDate, lastDate string) (string, error) {
	// Check the most recent date only for now
	if len(appointments) <= 0 || !isSuitableDate(appointments[0].AppointmentDt.Date, startDate, lastDate) {
		return "", fmt.Errorf("no available appointments between %s and %s", startDate, lastDate)
	}

	return appointments[0].AppointmentDt.Date, nil
}

func isSuitableDate(examDate, startDate, lastDate string) bool {
	return examDate > startDate && examDate < lastDate
}
