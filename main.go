package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var TimeFormat = time.UnixDate

func DisplayUsage() {
	fmt.Println("Usage: ./gots [-h] [date_string|unix_timestamp|[+|-123s|m|h|d|M|y]]")
}

func OperatorHelper(op string, t time.Time, seconds int64) time.Time {
	switch op {
	case "+":
		return t.Add(time.Duration(seconds) * time.Second)
	case "-":
		return t.Add(time.Duration(seconds) * time.Second * -1)
	default:
		return time.Time{}
	}
}

// ShiftTime returns a shifted unix timestamp + rfc time string
func ShiftTime(t time.Time, op, measure string, value int64) time.Time {
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

	return OperatorHelper(op, t, seconds)
}

func DisplayFullTime(t time.Time) {
	fmt.Printf("RFC3339: %v\n", t.Format(time.RFC3339))
	fmt.Printf("Unix: %v\n", t.Unix())
	fmt.Printf("UnixNano: %v\n", t.UnixNano())
	fmt.Println("---------------------")
	fmt.Printf("RFC3339 (UTC): %v\n", t.UTC().Format(time.RFC3339))
	fmt.Printf("Unix (UTC): %v\n", t.UTC().Unix())
	fmt.Printf("UnixNano (UTC): %v\n", t.UTC().UnixNano())
}

func ConvertTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse int ts: %v", err)
	}

	return time.Unix(i, 0), nil
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

func ParseTimeshiftArgs(timeshift string) []string {
	shiftRegex, regexErr := regexp.Compile(`(\+|\-)(\d+)(s|m|h|d|M|y)$`)
	if regexErr != nil {
		fmt.Printf("ERROR: Unable to create time shift regex (E: %v)\n", regexErr)
		os.Exit(1)
	}

	return shiftRegex.FindStringSubmatch(timeshift)
}

func HandleTimeshift(t time.Time, operator, value, measure string) (time.Time, error) {
	fmt.Println("handling time shift")
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("ERROR: Unable to parse int from string: %s", err)
	}

	return ShiftTime(t, operator, measure, v), nil
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
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		DisplayUsage()
		os.Exit(1)
	}

	var (
		t   time.Time
		err error
	)

	if len(os.Args) == 1 {
		// If no args, display current TS
		t = time.Now()
	} else {
		if IsTimestamp(os.Args[1]) {
			t, err = ConvertTimestamp(os.Args[1])
		} else if IsDateString(os.Args[1]) {
			t, err = ConvertDate(os.Args[1])
		} else if IsNanoTimestamp(os.Args[1]) {
			t, err = ConvertNanoTimestamp(os.Args[1])
		} else {
			fmt.Printf("ERROR: Unable to parse argument: %v\n", os.Args[1])
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("ERROR: Unable to convert input time: %s\n", err)
			os.Exit(1)
		}

		if len(os.Args) == 3 {
			if tsArgs := ParseTimeshiftArgs(os.Args[2]); len(tsArgs) == 4 {
				operator := tsArgs[1]
				value := tsArgs[2]
				measure := tsArgs[3]

				t, err = HandleTimeshift(t, operator, value, measure)
				if err != nil {
					fmt.Printf("ERROR: Unable to shift time: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	DisplayFullTime(t)
	os.Exit(0)
}

func IsNanoTimestamp(timestamp string) bool {
	if match, _ := regexp.MatchString(`^\d{19,19}$`, timestamp); match {
		return true
	}

	return false
}

func ConvertNanoTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse nano ts: %v", err)
	}

	return time.Unix(0, i), nil
}
