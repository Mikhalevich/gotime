package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormatLayout = time.RFC3339
)

func formatTime(t time.Time, utc bool) string {
	if utc {
		t = t.UTC()
	}
	return t.Format(timeFormatLayout)
}

// makeTime converts argument into time struct
// returns time struct, is argument in unix format flag, error
func makeTime(arg string) (time.Time, bool, error) {
	secs, err := strconv.Atoi(arg)
	if err == nil {
		return time.Unix(int64(secs), 0), true, nil
	}

	t, err := time.Parse(timeFormatLayout, arg)
	if err != nil {
		return time.Time{}, false, err
	}

	return t, false, nil
}

func convertTime(arg string, utc bool) (string, error) {
	t, isUnix, err := makeTime(arg)
	if err != nil {
		return "", err
	}

	if isUnix {
		return formatTime(t, utc), nil
	}

	return strconv.Itoa(int(t.Unix())), nil
}

func addTimes(times []string) (time.Time, error) {
	if len(times) != 2 {
		return time.Time{}, fmt.Errorf("invalid arguents: %s", strings.Join(times, " "))
	}

	t, _, err := makeTime(times[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time in first argument: %w", err)
	}

	d, err := strconv.Atoi(times[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time duration in second argument: %w", err)
	}

	return t.Add(time.Duration(d) * time.Second), nil
}

func main() {
	var (
		requestCurrentTime bool
		localTimeZone      bool
		requestAddTime     bool
	)

	flag.BoolVar(&requestCurrentTime, "c", false, "request current time")
	flag.BoolVar(&localTimeZone, "l", false, "use local timezone")
	flag.BoolVar(&requestAddTime, "a", false, "add time in first argument duration in second argument")
	flag.Parse()

	if requestCurrentTime {
		t := time.Now()
		fmt.Printf("%d\n%s\n", t.Unix(), formatTime(t, !localTimeZone))
		return
	}

	if requestAddTime {
		t, err := addTimes(flag.Args())
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Println(formatTime(t, !localTimeZone))
		return
	}

	for _, arg := range flag.Args() {
		t, err := convertTime(arg, !localTimeZone)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Println(t)
	}
}
