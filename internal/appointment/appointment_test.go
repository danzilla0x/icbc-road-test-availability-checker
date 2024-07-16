package appointment

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFindTargetAppointment(t *testing.T) {
	testCases := []struct {
		name              string
		examLastDate      string
		availableExamDate string
		responseBody      string
		expectedFound     bool
	}{
		{
			name:              "Appointment found",
			examLastDate:      "2024-08-10",
			availableExamDate: "2024-07-23",
			responseBody: `
				[
					{
						"appointmentDt": {
							"date": "2024-07-23",
							"dayOfWeek": "Wednesday"
						},
						"endTm": "09:55",
						"lemgMsgId": 1,
						"posId": 0,
						"resourceId": 0,
						"signature": "aaaa",
						"startTm": "09:20"
					}
				]`,
			expectedFound: true,
		},
		{
			name:              "Appointment NOT found",
			examLastDate:      "2024-08-10",
			availableExamDate: "2024-08-23", // FIX duplication
			responseBody: `
				[
					{
						"appointmentDt": {
							"date": "2024-08-23",
							"dayOfWeek": "Wednesday"
						},
						"endTm": "09:55",
						"lemgMsgId": 1,
						"posId": 0,
						"resourceId": 0,
						"signature": "aaaa",
						"startTm": "09:20"
					}
				]`,
			expectedFound: false,
		},
	}

	today := time.Now().Format("2006-01-02")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var appointments []Appointment
			err := json.Unmarshal([]byte(tc.responseBody), &appointments)
			if err != nil {
				t.Errorf("failed to unmarshel test json into []Appointment, err: " + err.Error())
			}

			availableExamDate, err := FindExamAppointment(appointments, today, tc.examLastDate)
			if (err != nil) == tc.expectedFound {
				t.Errorf("expected appointment present: %v, got: %v", tc.expectedFound, err)
			}

			if tc.expectedFound && availableExamDate != tc.availableExamDate {
				t.Errorf("expected date %v, got: %v", tc.availableExamDate, err)
			}
		})
	}
}
