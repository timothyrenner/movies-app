package cmd

import (
	"database/sql"
	"regexp"
	"strconv"
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
