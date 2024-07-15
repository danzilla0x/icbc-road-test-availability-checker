package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"icbc-checker/internal/appointment"
	"icbc-checker/internal/auth"
)

const (
	LOGIN_URL        = "https://onlinebusiness.icbc.com/deas-api/v1/webLogin/webLogin"
	APPOINTMENTS_URL = "https://onlinebusiness.icbc.com/deas-api/v1/web/getAvailableAppointments"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	userAgent := os.Getenv("APP_USER_AGENT")
	examLastDate := os.Getenv("APP_EXAM_LAST_DATE")
	lastName := os.Getenv("APP_LAST_NAME")
	if lastName == "" {
		panic("lastName is empty.")
	}
	licenseNumber := os.Getenv("APP_LICENSE_NUMBER")
	if licenseNumber == "" {
		panic("licenseNumber is empty.")
	}
	keyword := os.Getenv("APP_KEYWORD")
	if keyword == "" {
		panic("keyword is empty.")
	}
	aPosIDStr := os.Getenv("APP_APPOINTMENT_POSITION_ID")
	aPosID, err := strconv.ParseInt(aPosIDStr, 10, 32)
	if err != nil {
		panic("aPosID is incorecrly set: " + err.Error())
	}

	bearerToken, err := auth.GetBearerToken(LOGIN_URL, lastName, licenseNumber, keyword, userAgent)
	if err != nil {
		fmt.Println("Login error: " + err.Error())
		return
	}

	appointments, err := appointment.GetAvailableAppointments(APPOINTMENTS_URL, int32(aPosID), lastName, licenseNumber, bearerToken, userAgent)
	if err != nil {
		fmt.Println("Get appointments error: " + err.Error())
		return
	}

	availableExamDate, err := appointment.FindExamAppointment(appointments, examLastDate)
	if err != nil {
		fmt.Println("Upsss: " + err.Error())
		return
	}

	fmt.Println("Available date: " + availableExamDate)
}
