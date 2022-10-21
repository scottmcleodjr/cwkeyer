package cwkeyer

import (
	"reflect"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {
	tests := []struct {
		input         rune
		eventsWanted  []event
		presentWanted bool
	}{
		{
			input:         'A',
			eventsWanted:  []event{dit, dah, charSpace},
			presentWanted: true,
		}, {
			input:         'J',
			eventsWanted:  []event{dit, dah, dah, dah, charSpace},
			presentWanted: true,
		}, {
			input:         '3',
			eventsWanted:  []event{dit, dit, dit, dah, dah, charSpace},
			presentWanted: true,
		}, {
			input:         't',
			eventsWanted:  []event{dah, charSpace},
			presentWanted: true,
		}, {
			input:         '!',
			eventsWanted:  []event{},
			presentWanted: false,
		}, {
			input:         '\'',
			eventsWanted:  []event{},
			presentWanted: false,
		}, {
			input:         ' ',
			eventsWanted:  []event{wordSpace},
			presentWanted: true,
		},
	}

	for _, test := range tests {
		events, present := events(test.input)
		if !reflect.DeepEqual(events, test.eventsWanted) {
			t.Errorf("got %+v, want %+v for input %q", events, test.eventsWanted, test.input)
		}
		if present != test.presentWanted {
			t.Errorf("got %t, want %t for input %q", present, test.presentWanted, test.input)
		}
	}
}

func TestIsKeyable(t *testing.T) {
	tests := []struct {
		input rune
		want  bool
	}{
		{input: 'R', want: true},
		{input: 'z', want: true},
		{input: '?', want: true},
		{input: '7', want: true},
		{input: ',', want: false},
		{input: '&', want: false},
	}

	for _, test := range tests {
		got := IsKeyable(test.input)
		if got != test.want {
			t.Errorf("got %t, want %t for input %q", got, test.want, test.input)
		}
	}
}

func TestEventLength(t *testing.T) {
	tests := []struct {
		inputEvent event
		inputSpeed int
		want       int64 // Match output of time.Duration.Milliseconds()
	}{
		{inputEvent: dit, inputSpeed: 26, want: 47},
		{inputEvent: dit, inputSpeed: 15, want: 83},
		{inputEvent: dah, inputSpeed: 12, want: 300},
		{inputEvent: dah, inputSpeed: 41, want: 87},
		{inputEvent: space, inputSpeed: 36, want: 33},
		{inputEvent: charSpace, inputSpeed: 17, want: 142},
		{inputEvent: wordSpace, inputSpeed: 28, want: 258},
	}

	for _, test := range tests {
		got := time.Duration.Milliseconds(eventLength(test.inputEvent, test.inputSpeed))
		if got != test.want {
			t.Errorf("got %d, want %d for input event %d at speed %d",
				got, test.want, test.inputEvent, test.inputSpeed)
		}
	}
}
