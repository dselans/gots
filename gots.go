package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var TIME_FORMAT string = time.UnixDate

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

// Returns a shifted unix timestamp + rfc time string
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
	return strconv.FormatInt(newTs, 10), time.Unix(newTs, 0).Format(TIME_FORMAT)
}

func DisplayCurrentTimestamp() {
	t := time.Now()
	fmt.Printf("%v <=> %v\n", t.Unix(), t.Format(TIME_FORMAT))
}

func ConvertTimestamp(ts string) (string, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return "", err
	}

	return time.Unix(i, 0).Format(TIME_FORMAT), nil
}

func ConvertDate(dateString string) (time.Time, error) {
	return time.Parse("01/02/2006 15:04:05", dateString)
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" {
			DisplayUsage()
			os.Exit(0)
		}

		// Is this a timestamp?
		if match, _ := regexp.MatchString("^\\d{10,10}$", os.Args[1]); match {
			ts, err := ConvertTimestamp(os.Args[1])
			if err != nil {
				fmt.Printf("ERROR: Unable to convert timestamp. (%v)\n", err)
				os.Exit(1)
			}

			fmt.Printf("%v <=> %v\n", os.Args[1], ts)
			os.Exit(1)
		}

		// Is this a timeshift request?
		shiftRegex, regexErr := regexp.Compile(`(\+|\-)(\d+)(s|m|h|d|M|y)$`)
		if regexErr != nil {
			fmt.Printf("ERROR: Unable to create shift regex. (%v)\n", regexErr)
			os.Exit(1)
		}

		match := shiftRegex.FindStringSubmatch(os.Args[1])
		if len(match) == 4 {
			value, parseErr := strconv.ParseInt(match[2], 10, 64)
			if parseErr != nil {
				fmt.Printf("ERROR: Unable to parse int from string. (%v)\n", parseErr)
				os.Exit(1)
			}

			tsTime, rfcTime := ShiftTime(match[1], match[3], value)
			fmt.Printf("%v <=> %v\n", tsTime, rfcTime)
			os.Exit(1)
		}

		// Last chance, maybe a date string?
		if match, _ := regexp.MatchString(`^\d+/\d+/\d+\s+\d+:\d+:\d+$`, os.Args[1]); match {
			date, dateErr := ConvertDate(os.Args[1])
			if dateErr != nil {
				fmt.Printf("ERROR: Unable to convert date argument. (%v)\n", dateErr)
				os.Exit(1)
			}

			fmt.Printf("%v <=> %v\n", date.Unix(), date.Format(TIME_FORMAT))
			os.Exit(1)
		}

		// Sorry, nothing else we can do
		fmt.Println("ERROR: Unable to determine argument; see usage (-h)")
		os.Exit(1)
	}

	// No args, display current timestamp
	DisplayCurrentTimestamp()
	os.Exit(0)
}
