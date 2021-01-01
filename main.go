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

func formatTime(t time.Time, unixFormat bool, utc bool) string {
	if unixFormat {
		return strconv.Itoa(int(t.Unix()))
	}

	if utc {
		t = t.UTC()
	}
	return t.Format(timeFormatLayout)
}

func formatDuration(d time.Duration, unixFormat bool) string {
	if unixFormat {
		return strconv.Itoa(int(d.Seconds()))
	}
	return d.String()
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
	return formatTime(t, !isUnix, utc), nil
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

func diffTimes(times []string) (time.Duration, error) {
	if len(times) != 2 {
		return 0, fmt.Errorf("invalid arguents: %s", strings.Join(times, " "))
	}

	startTime, _, err := makeTime(times[0])
	if err != nil {
		return 0, fmt.Errorf("invalid time in first argument: %w", err)
	}

	finishTime, _, err := makeTime(times[1])
	if err != nil {
		return 0, fmt.Errorf("invalid time in second argument: %w", err)
	}

	return startTime.Sub(finishTime), nil
}

func main() {
	var (
		requestCurrentTime bool
		localTimeZone      bool
		unixTimeFormat     bool
		requestAddTime     bool
		requestDiffTime    bool
	)

	flag.BoolVar(&requestCurrentTime, "c", false, "request current time")
	flag.BoolVar(&localTimeZone, "l", false, "use local timezone")
	flag.BoolVar(&unixTimeFormat, "u", false, "print time in unix format")
	flag.BoolVar(&requestAddTime, "a", false, "add time in first argument with duration in second argument")
	flag.BoolVar(&requestDiffTime, "d", false, "calculate time difference")
	flag.Parse()

	if requestCurrentTime {
		fmt.Println(formatTime(time.Now(), unixTimeFormat, !localTimeZone))
		return
	}

	if requestAddTime {
		t, err := addTimes(flag.Args())
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Println(formatTime(t, unixTimeFormat, !localTimeZone))
		return
	}

	if requestDiffTime {
		diff, err := diffTimes(flag.Args())
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Println(formatDuration(diff, unixTimeFormat))
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
