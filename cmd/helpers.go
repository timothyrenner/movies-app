package cmd

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var runtimeRegex = regexp.MustCompile("([0-9]+) min")

func ParseRuntime(runtimeString string) *sql.NullInt32 {
	var runtime sql.NullInt32
	runtimeMatch := runtimeRegex.FindStringSubmatch(runtimeString)
	if runtimeMatch == nil {
		runtime.Int32 = 0
		runtime.Valid = false
	} else if len(runtimeMatch) <= 1 {
		runtime.Int32 = 0
		runtime.Valid = false
	} else {
		runtimeStr := runtimeMatch[1]
		runtimeInt, err := strconv.Atoi(runtimeStr)
		if err != nil {
			runtime.Int32 = 0
			runtime.Valid = false
		}
		runtime.Int32 = int32(runtimeInt)
		runtime.Valid = true
	}
	return &runtime
}

func ParseReleased(releasedString string) (string, error) {
	var releasedDate string
	if releasedString == "N/A" {
		releasedDate = releasedString
	} else {
		released, err := time.Parse("2 Jan 2006", releasedString)
		if err != nil {
			return "", fmt.Errorf(
				"error parsing date %v: %v", releasedString, err,
			)
		}
		releasedDate = released.Format("2006-01-02")
	}
	return releasedDate, nil
}
