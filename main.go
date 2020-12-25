package main

import (
	"flag"
	"fmt"
	"strconv"
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

func convertTime(arg string, utc bool) (string, error) {
	secs, err := strconv.Atoi(arg)
	if err == nil {
		t := time.Unix(int64(secs), 0)
		return formatTime(t, utc), nil
	}

	t, err := time.Parse(timeFormatLayout, arg)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(t.Unix())), nil
}

func main() {
	var (
		requestCurrentTime bool
		localTimeZone      bool
	)

	flag.BoolVar(&requestCurrentTime, "c", false, "request current time")
	flag.BoolVar(&localTimeZone, "l", false, "use local timezone")
	flag.Parse()

	if requestCurrentTime {
		t := time.Now()
		fmt.Printf("%d\n%s\n", t.Unix(), formatTime(t, !localTimeZone))
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
