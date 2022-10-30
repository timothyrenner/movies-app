package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var runtimeRegex = regexp.MustCompile("([0-9]+)( min)?")

func ParseRuntime(runtimeString string) (int, error) {
	runtimeMatch := runtimeRegex.FindStringSubmatch(runtimeString)
	if runtimeMatch == nil {
		return 0, fmt.Errorf("unable to parse runtime %v", runtimeString)
	} else if len(runtimeMatch) <= 1 {
		return 0, fmt.Errorf(
			"got multiple matches for runtime %v: %v",
			runtimeString, len(runtimeMatch),
		)
	} else {
		runtimeStr := runtimeMatch[1]
		runtimeInt, err := strconv.Atoi(runtimeStr)
		if err != nil {
			return 0, fmt.Errorf(
				"error converting match %v to int: %v", runtimeStr, err,
			)
		}
		return runtimeInt, nil
	}
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

func SplitOnCommaAndTrim(toSplit string) []string {
	splitStrings := strings.Split(toSplit, ",")
	stringSlice := make([]string, len(splitStrings))
	for ii := range stringSlice {
		spacesTrimmed := strings.TrimSpace(splitStrings[ii])
		leftBracketsTrimmed := strings.Trim(spacesTrimmed, "[")
		rightBracketsTrimmed := strings.Trim(leftBracketsTrimmed, "]")
		stringSlice[ii] = rightBracketsTrimmed
	}
	return stringSlice
}
