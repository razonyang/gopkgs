package helper

import "time"

func CurrentUTC() time.Time {
	return time.Now().UTC()
}
