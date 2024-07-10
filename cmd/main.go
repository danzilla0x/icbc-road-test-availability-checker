package main

import (
	"encoding/json"
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
	aPosIDStr := os.Getenv("APP_APPOINTMENT_POSITION_ID")
	lastName := os.Getenv("APP_LAST_NAME")
	licenseNumber := os.Getenv("APP_LICENSE_NUMBER")
	keyword := os.Getenv("APP_KEYWORD")

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
	}

	res, err := PrettyStruct(appointments)
	if err != nil {
		fmt.Println("PrettyStruct didn't work: " + err.Error())
	}

	fmt.Println(res)
}

func PrettyStruct(data any) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
