package util

import (
	"git.multiverse.io/eventkit/kit/constant"
	"time"
)

// ValidDate is used to check whether the format of the input time is `YYYY-MM-DD`
func ValidDate(timeStr string) error {
	_, err := time.Parse("2006-01-02", timeStr)
	return err
}

// GetFirstDateOfMonth is used to get first day of month
func GetFirstDateOfMonth() string {
	d := time.Now()
	timestamp := time.Now().AddDate(0, 0, -d.Day()+1).Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}

// GetCurrentTimeFmt is used to get current time ,fmt: hh:mm:ss
func GetCurrentTimeFmt() string {
	return time.Now().Format("15:04:05")
}

// GetTimeFormat is used to format timestamp to string(format:"yyyy-mm-dd hh:mm:ss")
func GetTimeFormat(timestmp int64) string {
	t := time.Unix(timestmp, 0).Format(constant.TimeStamp)
	return t
}

// GetTimeParse is used to parse the input time to time.Time, error return if the input parameter is invalid.
func GetTimeParse(timestr string) (time.Time, error) {
	t, err := time.ParseInLocation(constant.TimeStamp, timestr, time.Local)
	return t, err
}

// GetCurTime returns the current time(call time.Now to get current time).
func GetCurTime() time.Time {
	t := time.Now()
	return t
}

// GetCurrentTimeStamp returns a pair of current date and current time(call time.Now to get current time), format the pair result as <`YYYY-MM-DD`, `hh:mm:ss`>
func GetCurrentTimeStamp() (string, string) {
	timestamp := time.Now()
	return timestamp.Format("2006-01-02"), timestamp.Format("15:04:05")
}

// GetCurrentDateTime returns formatted current date and time, The format of result is `YYYY-MM-DD hh:mm:ss`
func GetCurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurrentDate returns formatted current date, The format of date is `YYYY-MM-DD`
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentTime returns formatted current time, The format of time is `hh:mm:ss`
func GetCurrentTime() string {
	return time.Now().Format("15:04:05")
}

// GetPreMonthDate is used to get the date of months ago
func GetPreMonthDate(monthNum int) string {
	now := time.Now()
	yesterday := now.AddDate(0, -monthNum, 0)
	return yesterday.Format("2006-01-02")
}

// GetTHCurrentDateTime is used to get current datetime of Thailand, the format of result is `YYYY-MM-DD hh:mm:ss`
func GetTHCurrentDateTime() string {
	var cstZone = time.FixedZone("CST", 7*3600)
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05")
}
