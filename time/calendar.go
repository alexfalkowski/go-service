package time

import "time"

// January is an alias of [time.January].
const January = time.January

// February is an alias of [time.February].
const February = time.February

// March is an alias of [time.March].
const March = time.March

// April is an alias of [time.April].
const April = time.April

// May is an alias of [time.May].
const May = time.May

// June is an alias of [time.June].
const June = time.June

// July is an alias of [time.July].
const July = time.July

// August is an alias of [time.August].
const August = time.August

// September is an alias of [time.September].
const September = time.September

// October is an alias of [time.October].
const October = time.October

// November is an alias of [time.November].
const November = time.November

// December is an alias of [time.December].
const December = time.December

// Sunday is an alias of [time.Sunday].
const Sunday = time.Sunday

// Monday is an alias of [time.Monday].
const Monday = time.Monday

// Tuesday is an alias of [time.Tuesday].
const Tuesday = time.Tuesday

// Wednesday is an alias of [time.Wednesday].
const Wednesday = time.Wednesday

// Thursday is an alias of [time.Thursday].
const Thursday = time.Thursday

// Friday is an alias of [time.Friday].
const Friday = time.Friday

// Saturday is an alias of [time.Saturday].
const Saturday = time.Saturday

// Month is an alias of [time.Month].
type Month = time.Month

// Weekday is an alias of [time.Weekday].
type Weekday = time.Weekday

// Location is an alias of [time.Location].
type Location = time.Location

// UTC is initialized from [time.UTC]. Assigning it changes only this package's variable.
var UTC = time.UTC

// Date wraps [time.Date] to construct a calendar time in loc. It normalizes
// out-of-range calendar fields and panics when loc is nil.
func Date(year int, month Month, day, hour, minute, second, nanosecond int, loc *Location) Time {
	return time.Date(year, month, day, hour, minute, second, nanosecond, loc)
}

// FixedZone wraps [time.FixedZone] to construct a location that is always
// offset seconds east of UTC.
func FixedZone(name string, offset int) *Location {
	return time.FixedZone(name, offset)
}

// LoadLocation wraps [time.LoadLocation] to load a named IANA time zone from
// Go's time zone database search path. The name "Local" selects the system time zone.
func LoadLocation(name string) (*Location, error) {
	return time.LoadLocation(name)
}

// LoadLocationFromTZData wraps [time.LoadLocationFromTZData] to load a named
// time zone from IANA TZif data.
func LoadLocationFromTZData(name string, data []byte) (*Location, error) {
	return time.LoadLocationFromTZData(name, data)
}
