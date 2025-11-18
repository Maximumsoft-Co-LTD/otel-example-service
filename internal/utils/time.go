package utils

import "time"

var Now time.Time

func init() {
	Now = time.Now()
}

func TimeNowUTC() time.Time {
	return Now.UTC()
}
