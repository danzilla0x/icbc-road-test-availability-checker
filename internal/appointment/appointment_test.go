package appointment

import (
	"encoding/json"
	"testing"
)

func TestFindTargetAppointment(t *testing.T) {
	testResponseBody := `
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
		]`

	lastDate := "2024-08-10"

	t.Run("Taget appointment present", func(t *testing.T) {
		var appointments []Appointment
		err := json.Unmarshal([]byte(testResponseBody), &appointments)
		if err != nil {
			t.Errorf("failed to unmarshel test json into []Appointment, err: " + err.Error())
		}

		availableExamDate, err := FindExamAppointment(appointments, lastDate)
		if err != nil {
			t.Errorf("expected appointment present: %v, got: %v", true, err)
		}

		if availableExamDate != "2024-07-23" {
			t.Errorf("expected date: 2024-07-23, got: %v", err)
		}
	})
}
