package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tommyknows/gocal"
)

const ics = `
BEGIN:VEVENT
DTSTART;TZID=Europe/Paris:20190202T130000
DTEND;TZID=Europe/Paris:20190202T150000
DTSTAMP:20180816T112126Z
UID:4agntbp30gkhdh3cs5ou1jl34q@google.com
RECURRENCE-ID;TZID=Europe/Paris:20190202T090000
CREATED:20180816T110948Z
DESCRIPTION:
LAST-MODIFIED:20180816T112123Z
LOCATION:
SEQUENCE:1
STATUS:CONFIRMED
SUMMARY:1st of month
TRANSP:OPAQUE
END:VEVENT

BEGIN:VEVENT
DTSTART;TZID=Europe/Paris:20190102T090000
DTEND;TZID=Europe/Paris:20190102T110000
DTSTAMP:20180816T112126Z
UID:4agntbp30gkhdh3cs5ou1jl34q@google.com
RECURRENCE-ID;TZID=Europe/Paris:20190102T090000
CREATED:20180816T110948Z
DESCRIPTION:
LAST-MODIFIED:20180816T111457Z
LOCATION:
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:1st of month (edited)
TRANSP:OPAQUE
END:VEVENT

BEGIN:VEVENT
DTSTART;TZID=Europe/Paris:20180802T090000
DTEND;TZID=Europe/Paris:20180802T110000
RRULE:FREQ=MONTHLY;BYMONTHDAY=2
EXDATE;TZID=Europe/Paris:20181202T090000
DTSTAMP:20180816T112126Z
UID:4agntbp30gkhdh3cs5ou1jl34q@google.com
CREATED:20180816T110948Z
DESCRIPTION:
LAST-MODIFIED:20180816T110948Z
LOCATION:
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:1st of month
TRANSP:OPAQUE
END:VEVENT
`

func main() {
	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	c := gocal.NewParser(strings.NewReader(ics))
	c.Start, c.End = &start, &end
	c.Parse()

	for _, e := range c.Events {
		fmt.Printf("%s on %s - %s\n", e.Summary, e.Start, e.End)
	}
}
