package strftime

import (
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	fullWeekDayName             = timefmt("Monday")
	abbrvWeekDayName            = timefmt("Mon")
	fullMonthName               = timefmt("January")
	abbrvMonthName              = timefmt("Jan")
	centuryDecimal              = appenderFn(appendCentury)
	timeAndDate                 = timefmt("Mon Jan _2 15:04:05 2006")
	mdy                         = timefmt("01/02/06")
	dayOfMonthZeroPad           = timefmt("02")
	dayOfMonthSpacePad          = timefmt("_2")
	ymd                         = timefmt("2006-01-02")
	twentyFourHourClockZeroPad  = timefmt("15")
	twelveHourClockZeroPad      = timefmt("3")
	dayOfYear                   = appenderFn(appendDayOfYear)
	twentyFourHourClockSpacePad = hourwblank(false)
	twelveHourClockSpacePad     = hourwblank(true)
	minutesZeroPad              = timefmt("04")
	monthNumberZeroPad          = timefmt("01")
	newline                     = verbatim("\n")
	ampm                        = timefmt("PM")
	hm                          = timefmt("15:04")
	imsp                        = timefmt("3:04:05 PM")
	secondsNumberZeroPad        = timefmt("05")
	hms                         = timefmt("15:04:05")
	tab                         = verbatim("\t")
	weekNumberSundayOrigin      = weeknumberOffset(0) // week number of the year, Sunday first
	weekdayMondayOrigin         = weekday(1)
	// monday as the first day, and 01 as the first value
	weekNumberMondayOriginOneOrigin = appenderFn(appendWeekNumber)
	eby                             = timefmt("_2-Jan-2006")
	// monday as the first day, and 00 as the first value
	weekNumberMondayOrigin = weeknumberOffset(1) // week number of the year, Monday first
	weekdaySundayOrigin    = weekday(0)
	natReprTime            = timefmt("15:04:05") // national representation of the time XXX is this correct?
	natReprDate            = timefmt("01/02/06") // national representation of the date XXX is this correct?
	year                   = timefmt("2006")     // year with century
	yearNoCentury          = timefmt("06")       // year w/o century
	timezone               = timefmt("MST")      // time zone name
	timezoneOffset         = timefmt("-0700")    // time zone ofset from UTC
	percent                = verbatim("%")
)

func lookupDirective(key byte) (appender, bool) {
	switch key {
	case 'A':
		return fullWeekDayName, true
	case 'a':
		return abbrvWeekDayName, true
	case 'B':
		return fullMonthName, true
	case 'b', 'h':
		return abbrvMonthName, true
	case 'C':
		return centuryDecimal, true
	case 'c':
		return timeAndDate, true
	case 'D':
		return mdy, true
	case 'd':
		return dayOfMonthZeroPad, true
	case 'e':
		return dayOfMonthSpacePad, true
	case 'F':
		return ymd, true
	case 'H':
		return twentyFourHourClockZeroPad, true
	case 'I':
		return twelveHourClockZeroPad, true
	case 'j':
		return dayOfYear, true
	case 'k':
		return twentyFourHourClockSpacePad, true
	case 'l':
		return twelveHourClockSpacePad, true
	case 'M':
		return minutesZeroPad, true
	case 'm':
		return monthNumberZeroPad, true
	case 'n':
		return newline, true
	case 'p':
		return ampm, true
	case 'R':
		return hm, true
	case 'r':
		return imsp, true
	case 'S':
		return secondsNumberZeroPad, true
	case 'T':
		return hms, true
	case 't':
		return tab, true
	case 'U':
		return weekNumberSundayOrigin, true
	case 'u':
		return weekdayMondayOrigin, true
	case 'V':
		return weekNumberMondayOriginOneOrigin, true
	case 'v':
		return eby, true
	case 'W':
		return weekNumberMondayOrigin, true
	case 'w':
		return weekdaySundayOrigin, true
	case 'X':
		return natReprTime, true
	case 'x':
		return natReprDate, true
	case 'Y':
		return year, true
	case 'y':
		return yearNoCentury, true
	case 'Z':
		return timezone, true
	case 'z':
		return timezoneOffset, true
	case '%':
		return percent, true
	}
	return nil, false
}

type combiningAppend struct {
	list           appenderList
	prev           appender
	prevCanCombine bool
}

func (ca *combiningAppend) Append(w appender) {
	if ca.prevCanCombine {
		if wc, ok := w.(combiner); ok && wc.canCombine() {
			ca.prev = ca.prev.(combiner).combine(wc)
			ca.list[len(ca.list)-1] = ca.prev
			return
		}
	}

	ca.list = append(ca.list, w)
	ca.prev = w
	ca.prevCanCombine = false
	if comb, ok := w.(combiner); ok {
		if comb.canCombine() {
			ca.prevCanCombine = true
		}
	}
}

func compile(wl *appenderList, p string) error {
	var ca combiningAppend
	for l := len(p); l > 0; l = len(p) {
		i := strings.IndexByte(p, '%')
		if i < 0 {
			ca.Append(verbatim(p))
			// this is silly, but I don't trust break keywords when there's a
			// possibility of this piece of code being rearranged
			p = p[l:]
			continue
		}
		if i == l-1 {
			return errors.New(`stray % at the end of pattern`)
		}

		// we found a '%'. we need the next byte to decide what to do next
		// we already know that i < l - 1
		// everything up to the i is verbatim
		if i > 0 {
			ca.Append(verbatim(p[:i]))
			p = p[i:]
		}

		directive, ok := lookupDirective(p[1])
		if !ok {
			return errors.Errorf(`unknown time format specification '%c'`, p[1])
		}
		ca.Append(directive)
		p = p[2:]
	}

	*wl = ca.list

	return nil
}

// Format takes the format `s` and the time `t` to produce the
// format date/time. Note that this function re-compiles the
// pattern every time it is called.
//
// If you know beforehand that you will be reusing the pattern
// within your application, consider creating a `Strftime` object
// and reusing it.
func Format(p string, t time.Time) (string, error) {
	var dst []byte
	// TODO: optimize for 64 byte strings
	dst = make([]byte, 0, len(p)+10)
	// Compile, but execute as we go
	for l := len(p); l > 0; l = len(p) {
		i := strings.IndexByte(p, '%')
		if i < 0 {
			dst = append(dst, p...)
			// this is silly, but I don't trust break keywords when there's a
			// possibility of this piece of code being rearranged
			p = p[l:]
			continue
		}
		if i == l-1 {
			return "", errors.New(`stray % at the end of pattern`)
		}

		// we found a '%'. we need the next byte to decide what to do next
		// we already know that i < l - 1
		// everything up to the i is verbatim
		if i > 0 {
			dst = append(dst, p[:i]...)
			p = p[i:]
		}

		directive, ok := lookupDirective(p[1])
		if !ok {
			return "", errors.Errorf(`unknown time format specification '%c'`, p[1])
		}
		dst = directive.Append(dst, t)
		p = p[2:]
	}

	return string(dst), nil
}

// Strftime is the object that represents a compiled strftime pattern
type Strftime struct {
	pattern  string
	compiled appenderList
}

// New creates a new Strftime object. If the compilation fails, then
// an error is returned in the second argument.
func New(f string) (*Strftime, error) {
	var wl appenderList
	if err := compile(&wl, f); err != nil {
		return nil, errors.Wrap(err, `failed to compile format`)
	}
	return &Strftime{
		pattern:  f,
		compiled: wl,
	}, nil
}

// Pattern returns the original pattern string
func (f *Strftime) Pattern() string {
	return f.pattern
}

// Format takes the destination `dst` and time `t`. It formats the date/time
// using the pre-compiled pattern, and outputs the results to `dst`
func (f *Strftime) Format(dst io.Writer, t time.Time) error {
	const bufSize = 64
	var b []byte
	max := len(f.pattern) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	if _, err := dst.Write(f.format(b, t)); err != nil {
		return err
	}
	return nil
}

func (f *Strftime) format(b []byte, t time.Time) []byte {
	for _, w := range f.compiled {
		b = w.Append(b, t)
	}
	return b
}

// FormatString takes the time `t` and formats it, returning the
// string containing the formated data.
func (f *Strftime) FormatString(t time.Time) string {
	const bufSize = 64
	var b []byte
	max := len(f.pattern) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	return string(f.format(b, t))
}
