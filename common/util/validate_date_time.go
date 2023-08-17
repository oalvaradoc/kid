package util

import (
	"fmt"
	"time"
)

// ValidateDate is used to check whether the format of input date is `YYYY-MM-DD`
func ValidateDate(date string) error {

	_, err := time.Parse("2006-01-02", date)
	//Date checked failed, the correct format is: YYYY-MM-DD
	if err != nil {
		return fmt.Errorf("date format checked failed, the correct format is: YYYY-MM-DD")
	}
	return nil

}

// ValidateTime is used to check whether the format of the input time is `HH:MM:SS`
func ValidateTime(t string) error {

	_, err := time.Parse("15:04:05", t)

	if err != nil {
		return fmt.Errorf("time format checked failed, the correct format is: HH:MM:SS")
	}

	return nil
}

// ValidateDateTime is used to check whether the format of the input data time is `YYYY-MM-DD HH:MM:SS`.
func ValidateDateTime(dateTime string) error {

	_, err := time.Parse("2006-01-02 15:04:05", dateTime)

	if err != nil {
		return fmt.Errorf("datetime format checked failed, the correct format is: YYYY-MM-DD HH:MM:SS")
	}

	return nil
}
