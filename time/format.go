package time

import "time"

// Layout is the reference time layout from [time.Layout].
const Layout = time.Layout

// ANSIC is the layout from [time.ANSIC].
const ANSIC = time.ANSIC

// UnixDate is the layout from [time.UnixDate].
const UnixDate = time.UnixDate

// RubyDate is the layout from [time.RubyDate].
const RubyDate = time.RubyDate

// RFC822 is the layout from [time.RFC822].
const RFC822 = time.RFC822

// RFC822Z is the layout from [time.RFC822Z], with a numeric time zone offset.
const RFC822Z = time.RFC822Z

// RFC850 is the layout from [time.RFC850].
const RFC850 = time.RFC850

// RFC1123 is the layout from [time.RFC1123].
const RFC1123 = time.RFC1123

// RFC1123Z is the layout from [time.RFC1123Z], with a numeric time zone offset.
const RFC1123Z = time.RFC1123Z

// RFC3339 is the layout from [time.RFC3339].
const RFC3339 = time.RFC3339

// RFC3339Nano is the layout from [time.RFC3339Nano], with up to nine fractional digits.
const RFC3339Nano = time.RFC3339Nano

// Kitchen is the layout from [time.Kitchen].
const Kitchen = time.Kitchen

// Stamp is the layout from [time.Stamp].
const Stamp = time.Stamp

// StampMilli is the layout from [time.StampMilli].
const StampMilli = time.StampMilli

// StampMicro is the layout from [time.StampMicro].
const StampMicro = time.StampMicro

// StampNano is the layout from [time.StampNano].
const StampNano = time.StampNano

// DateTime is the layout from [time.DateTime], which includes a date and time.
const DateTime = time.DateTime

// DateOnly is the layout from [time.DateOnly], which includes only the date.
const DateOnly = time.DateOnly

// TimeOnly is the layout from [time.TimeOnly], which includes only the time of day.
const TimeOnly = time.TimeOnly

// ParseError is an alias of [time.ParseError].
type ParseError = time.ParseError

// Parse wraps [time.Parse]. It parses value using layout, interprets values
// without a time zone as UTC, and returns a [*ParseError] for invalid input.
func Parse(layout, value string) (Time, error) {
	return time.Parse(layout, value)
}

// ParseInLocation wraps [time.ParseInLocation]. It uses loc to interpret values
// without a time zone and to resolve zone offsets and abbreviations.
func ParseInLocation(layout, value string, loc *Location) (Time, error) {
	return time.ParseInLocation(layout, value, loc)
}
