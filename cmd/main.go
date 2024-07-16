package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"icbc-checker/internal/appointment"
	"icbc-checker/internal/auth"
	"icbc-checker/internal/notification"
)

const (
	LOGIN_URL         = "https://onlinebusiness.icbc.com/deas-api/v1/webLogin/webLogin"
	APPOINTMENTS_URL  = "https://onlinebusiness.icbc.com/deas-api/v1/web/getAvailableAppointments"
	MIN_WAIT_INTERVAL = 20
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	aPosIDStr := os.Getenv("APP_APPOINTMENT_POSITION_ID")
	aPosID, err := strconv.ParseInt(aPosIDStr, 10, 32)
	if err != nil {
		panic("aPosID is incorecrly set: " + err.Error())
	}
	userAgent := os.Getenv("APP_USER_AGENT")
	examStartDate := os.Getenv("APP_EXAM_START_DATE")
	if examStartDate == "" {
		examStartDate = time.Now().Format("2006-01-02")
	}
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
	pushoverUserKey := os.Getenv("APP_PUSHOVER_USER_KEY")
	if pushoverUserKey == "" {
		fmt.Println("Pushover key is not provided.")
	}
	pushoverApiToken := os.Getenv("APP_PUSHOVER_API_TOKEN")
	if pushoverApiToken == "" {
		fmt.Println("Pushover key is not provided.")
	}

	for {
		bearerToken, err := auth.GetBearerToken(LOGIN_URL, lastName, licenseNumber, keyword, userAgent)
		if err != nil {
			fmt.Println("Login error: " + err.Error())
			return
		}

		for {
			appointments, err := appointment.GetAvailableAppointments(APPOINTMENTS_URL, int32(aPosID), lastName, licenseNumber, bearerToken, userAgent)
			if err != nil {
				fmt.Println("Get appointments error: " + err.Error())
				break
			}

			availableExamDate, err := appointment.FindExamAppointment(appointments, examStartDate, examLastDate)
			if err != nil {
				waitInterval := rand.Intn(20) + MIN_WAIT_INTERVAL
				fmt.Printf("Upsss: %s. Sleep for %d seconds.\n", err.Error(), waitInterval)
				time.Sleep(time.Duration(waitInterval) * time.Second)
				continue
			}

			err = notification.SendNotification(pushoverUserKey, pushoverApiToken, availableExamDate)
			if err != nil {
				fmt.Println("Failed to send notification: " + err.Error())
			}

			return
		}
	}
}
