package utils

import (
	"errors"
	"strconv"
	"time"
)

func ValidateSecurityCode(code string) error {
	if len(code) != 3 {
		return errors.New("the security code must be exactly 3 digits long")
	}
	return nil
}
func ValidateExpiryMonth(month string) error {
	if len(month) != 2 {
		return errors.New("the due month must have exactly 2 digits")
	}
	return nil
}

func ValidateExpiryYear(year string) error {
	if len(year) != 4 {
		return errors.New("the due year must have exactly 4 digits")
	}
	return nil
}

func ValidateExpiryDate(year string, month string) error {
	yearInt, err := strconv.ParseInt(year, 10, 64)
	if err != nil {
		return errors.New("error converting year to integer")
	}
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		return errors.New("error converting month to integer")
	}
	firstDayOfNextMonth := time.Date(int(yearInt), time.Month(int(monthInt))+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfNextMonth.Add(-24 * time.Hour)
	if !lastDayOfMonth.After(time.Now()) {
		return errors.New("the expiration date has expired")
	}
	return nil
}
