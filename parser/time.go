package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TimeStart = iota
	TimeEnd
)

var (
	TZMapper func(s string) (*time.Location, error)
)

func ParseTime(s string, params map[string]string, ty int, allday bool) (*time.Time, error) {
	var err error
	var tz *time.Location

	format := ""

	if params["VALUE"] == "DATE" || len(s) == 8 {
		t, err := time.Parse("20060102", s)
		if ty == TimeStart {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		} else if ty == TimeEnd {
			if allday {
				t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, time.UTC)
			} else {
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(-1 * time.Millisecond)
			}
		}

		return &t, err
	}

	if strings.HasSuffix(s, "Z") {
		// If string end in 'Z', timezone is UTC
		format = "20060102T150405Z"
		tz, _ = time.LoadLocation("UTC")
	} else if params["TZID"] != "" {
		var err error

		// If TZID param is given, parse in the timezone unless it is not valid
		format = "20060102T150405"
		if TZMapper != nil {
			tz, err = TZMapper(params["TZID"])
		}
		if TZMapper == nil || err != nil {
			tz, err = LoadTimezone(params["TZID"])
		}

		if err != nil {
			tz, _ = time.LoadLocation("UTC")
		}
	} else {
		// Else, consider the timezone is local the parser
		format = "20060102T150405"
		tz = time.Local
	}

	t, err := time.ParseInLocation(format, s, tz)

	return &t, err
}

var (
	full = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)
	week = regexp.MustCompile(`P((?P<week>\d+)W)`)
)

func ParseDuration(dur string) (*time.Duration, error) {
	var (
		match []string
		re    *regexp.Regexp
	)

	if week.MatchString(dur) {
		match = week.FindStringSubmatch(dur)
		re = week
	} else if full.MatchString(dur) {
		match = full.FindStringSubmatch(dur)
		re = full
	} else {
		return nil, fmt.Errorf("bad format string")
	}

	var years, weeks, days, hours, minutes, seconds int

	for i, name := range re.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		switch name {
		case "year":
			years = val
		case "month":
			return nil, fmt.Errorf("no months allowed")
		case "week":
			weeks = val
		case "day":
			days = val
		case "hour":
			hours = val
		case "minute":
			minutes = val
		case "second":
			seconds = val
		default:
			return nil, fmt.Errorf("unknown field %s", name)
		}
	}

	day := time.Hour * 24
	year := day * 365

	tot := time.Duration(0)
	tot += year * time.Duration(years)
	tot += day * 7 * time.Duration(weeks)
	tot += day * time.Duration(days)
	tot += time.Hour * time.Duration(hours)
	tot += time.Minute * time.Duration(minutes)
	tot += time.Second * time.Duration(seconds)

	return &tot, nil
}

func LoadTimezone(tzid string) (*time.Location, error) {
	tz, err := time.LoadLocation(tzid)
	if err == nil {
		return tz, err
	}

	tokens := strings.Split(tzid, "_")
	for idx, t := range tokens {
		t = strings.ToLower(t)

		if t != "of" && t != "es" {
			tokens[idx] = strings.Title(t)
		} else {
			tokens[idx] = t
		}
	}

	return time.LoadLocation(strings.Join(tokens, "_"))
}
