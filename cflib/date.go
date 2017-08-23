package cflib

import (
	"time"
	"strings"
	"strconv"
	"regexp"
	"fmt"
)

// Date simplify the usage of time
type Date struct {
	Location string
	timeInstance time.Time
}

// Format set the time format
func (d *Date) Format(format string) string{
	return d.timeInstance.Format(format)
}

// GetTimeInstance return the time object
func (d *Date) GetTimeInstance() time.Time{
	return d.timeInstance
}

// Today set the time object to current time
func (d *Date) Today() *Date {
	d.timeInstance = time.Now()

	if d.Location == "UTC" {
		d.timeInstance = d.timeInstance.UTC()
	}

	return d
}

// SetDate sets the time object to specific time with help of date string and time string
// In the following format:
// dateString: 2017-12-31
// timeString: 15:59:59
func (d *Date) 	SetDate(dateString string, timeString string) (*Date, error) {
	var loc *time.Location

	if d.Location == "UTC" {
		loc = time.UTC
	} else {
		loc = time.Local
	}

	goodDateFormat, err := regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`, dateString)
	if goodDateFormat == false || err != nil {
		err = fmt.Errorf("Incorrect date format. Use the followning format, YYYY-MM-DD.")
		return d, err
	}

	goodTimeFormat, err := regexp.MatchString(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`, timeString)
	if goodTimeFormat == false || err != nil {
		err = fmt.Errorf("Incorrect time format. Use the followning format, hh:mm:ss,")
		return d, err
	}

	dateSlice := strings.Split(dateString, "-")
	timeSlice := strings.Split(timeString, ":")
	year, _ := strconv.Atoi(dateSlice[0])
	month, _ := strconv.Atoi(dateSlice[1])
	day, _ := strconv.Atoi(dateSlice[2])
	hour, _ := strconv.Atoi(timeSlice[0])
	minute, _ := strconv.Atoi(timeSlice[1])
	second, _ := strconv.Atoi(timeSlice[2])

	d.timeInstance = time.Date(
		year,
		time.Month(month),
		day,
		hour,
		minute,
		second,
		0,
		loc,
	)

	return d, err
}

// IncDays increases the current time object with "incDays" days
// Negative number is allowed for decreasing the date
func (d *Date) IncDays(incDays int) *Date {
	d.timeInstance = d.timeInstance.AddDate(0, 0, incDays)
	return d
}
