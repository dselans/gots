package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var TimeFormat string = time.UnixDate

func DisplayUsage() {
	fmt.Println("Usage: ./gots [-h] [date_string|unix_timestamp|[+|-123s|m|h|d|M|y]]")
}

func OperatorHelper(op string, val1, val2 int64) int64 {
	switch op {
	case "+":
		return val1 + val2
	case "-":
		return val1 - val2
	default:
		return 0
	}
}

// ShiftTime returns a shifted unix timestamp + rfc time string
func ShiftTime(op, measure string, value int64) (string, string) {
	t := time.Now()
	ts := t.Unix()

	var multiplier int64

	switch measure {
	case "s":
		multiplier = 1 // 1 sec
	case "m":
		multiplier = 60 // 60 sec
	case "h":
		multiplier = 60 * 60 // seconds in an hour
	case "d":
		multiplier = (60 * 60) * 24
	case "M":
		multiplier = ((60 * 60) * 24) * 30 // crude, but oh well
	case "y":
		multiplier = ((60 * 60) * 24) * 365
	}

	seconds := multiplier * value
	newTs := OperatorHelper(op, ts, seconds)
	return strconv.FormatInt(newTs, 10), time.Unix(newTs, 0).Format(TimeFormat)
}

func DisplayCurrentTimestamp() {
	t := time.Now()

	fmt.Printf("RFC3339: %v\n", t.Format(time.RFC3339))
	fmt.Printf("Unix: %v\n", t.Unix())
	fmt.Printf("UnixNano: %v\n", t.UnixNano())
	fmt.Println("---------------------")
	fmt.Printf("RFC3339 (UTC): %v\n", t.UTC().Format(time.RFC3339))
	fmt.Printf("Unix (UTC): %v\n", t.UTC().Unix())
	fmt.Printf("UnixNano (UTC): %v\n", t.UTC().UnixNano())
}

func ConvertTimestamp(ts string) (string, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return "", err
	}

	return time.Unix(i, 0).Format(TimeFormat), nil
}

func ConvertDate(dateString string) (time.Time, error) {
	return time.Parse("01/02/2006 15:04:05", dateString)
}

func IsTimestamp(timestamp string) bool {
	if match, _ := regexp.MatchString(`^\d{10,10}$`, timestamp); match {
		return true
	}

	return false
}

func HandleTimestamp(timestamp string) {
	ts, err := ConvertTimestamp(timestamp)
	if err != nil {
		fmt.Printf("ERROR: Unable to convert timestamp (E: %v)\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v <=> %v\n", timestamp, ts)
	os.Exit(0)
}

func ParseTimeshiftArgs(timeshift string) []string {
	shiftRegex, regexErr := regexp.Compile(`(\+|\-)(\d+)(s|m|h|d|M|y)$`)
	if regexErr != nil {
		fmt.Printf("ERROR: Unable to create time shift regex (E: %v)\n", regexErr)
		os.Exit(1)
	}

	return shiftRegex.FindStringSubmatch(timeshift)
}

func HandleTimeshift(match []string) {
	value, parseErr := strconv.ParseInt(match[2], 10, 64)
	if parseErr != nil {
		fmt.Printf("ERROR: Unable to parse int from string (E: %v)\n", parseErr)
		os.Exit(1)
	}

	tsTime, rfcTime := ShiftTime(match[1], match[3], value)
	fmt.Printf("%v <=> %v\n", tsTime, rfcTime)
	os.Exit(0)
}

func IsDateString(date string) bool {
	if match, _ := regexp.MatchString(`^\d+/\d+/\d+\s+\d+:\d+:\d+$`, date); match {
		return true
	}

	return false
}

func HandleDate(dateArg string) {
	date, dateErr := ConvertDate(dateArg)
	if dateErr != nil {
		fmt.Printf("ERROR: Unable to convert date argument. (%v)\n", dateErr)
		os.Exit(1)
	}

	fmt.Printf("%v <=> %v\n", date.Unix(), date.Format(TimeFormat))
	os.Exit(0)
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" {
			DisplayUsage()
			os.Exit(0)
		}

		if IsTimestamp(os.Args[1]) {
			HandleTimestamp(os.Args[1])
		}

		if tsArgs := ParseTimeshiftArgs(os.Args[1]); len(tsArgs) == 4 {
			HandleTimeshift(tsArgs)
		}

		if IsDateString(os.Args[1]) {
			HandleDate(os.Args[1])
		}

		fmt.Println("ERROR: Unable to determine argument; see usage (-h)")
		os.Exit(1)
	}

	// No args, display current timestamp
	DisplayCurrentTimestamp()
	os.Exit(0)
}
