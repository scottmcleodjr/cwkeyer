package cwkeyer

import (
	"time"
	"unicode"
)

type event int

const (
	dit            event = 1
	dah            event = 2
	space          event = 3 // A dit-length space that plays after every event
	charSpace      event = 4 // A dah-length space that plays between characters
	wordSpace      event = 5 // A longer space that plays between words
	ditUnits       int   = 1
	dahUnits       int   = 3
	spaceUnits     int   = 1
	charSpaceUnits int   = 2 // Down one from 3 because it will always be followed by a space
	wordSpaceUnits int   = 6 // Down one from 7 because it will always be followed by a space
)

var runeToEvents = map[rune][]event{
	'A': {dit, dah, charSpace},
	'B': {dah, dit, dit, dit, charSpace},
	'C': {dah, dit, dah, dit, charSpace},
	'D': {dah, dit, dit, charSpace},
	'E': {dit, charSpace},
	'F': {dit, dit, dah, dit, charSpace},
	'G': {dah, dah, dit, charSpace},
	'H': {dit, dit, dit, dit, charSpace},
	'I': {dit, dit, charSpace},
	'J': {dit, dah, dah, dah, charSpace},
	'K': {dah, dit, dah, charSpace},
	'L': {dit, dah, dit, dit, charSpace},
	'M': {dah, dah, charSpace},
	'N': {dah, dit, charSpace},
	'O': {dah, dah, dah, charSpace},
	'P': {dit, dah, dah, dit, charSpace},
	'Q': {dah, dah, dit, dah, charSpace},
	'R': {dit, dah, dit, charSpace},
	'S': {dit, dit, dit, charSpace},
	'T': {dah, charSpace},
	'U': {dit, dit, dah, charSpace},
	'V': {dit, dit, dit, dah, charSpace},
	'W': {dit, dah, dah, charSpace},
	'X': {dah, dit, dit, dah, charSpace},
	'Y': {dah, dit, dah, dah, charSpace},
	'Z': {dah, dah, dit, dit, charSpace},
	'0': {dah, dah, dah, dah, dah, charSpace},
	'1': {dit, dah, dah, dah, dah, charSpace},
	'2': {dit, dit, dah, dah, dah, charSpace},
	'3': {dit, dit, dit, dah, dah, charSpace},
	'4': {dit, dit, dit, dit, dah, charSpace},
	'5': {dit, dit, dit, dit, dit, charSpace},
	'6': {dah, dit, dit, dit, dit, charSpace},
	'7': {dah, dah, dit, dit, dit, charSpace},
	'8': {dah, dah, dah, dit, dit, charSpace},
	'9': {dah, dah, dah, dah, dit, charSpace},
	'?': {dit, dit, dah, dah, dit, dit, charSpace},
	'/': {dah, dit, dit, dah, dit, charSpace},
	' ': {wordSpace},
}

func events(r rune) (events []event, present bool) {
	if !unicode.IsUpper(r) {
		r = unicode.ToUpper(r)
	}
	events, present = runeToEvents[r]
	if !present {
		// Avoid a blow-up if this value is used
		return []event{}, present
	}
	return events, present
}

// IsKeyable accepts a rune and returns true if the rune is keyable as CW.
func IsKeyable(r rune) bool {
	_, present := events(r)
	return present
}

func eventLength(e event, wpm int) time.Duration {
	millisDuration := func(units int) time.Duration {
		millis := units * (1000 / ((5 * wpm) / 6))
		return time.Duration(time.Millisecond * time.Duration(millis))
	}

	switch e {
	case dit:
		return millisDuration(ditUnits)
	case dah:
		return millisDuration(dahUnits)
	case space:
		return millisDuration(spaceUnits)
	case charSpace:
		return millisDuration(charSpaceUnits)
	case wordSpace:
		return millisDuration(wordSpaceUnits)
	}
	return 0
}
