package cwkeyer_test

import (
	"errors"
	"testing"
	"time"

	. "github.com/scottmcleodjr/cwkeyer"
)

type testSender struct {
	downMethod func() error
	upMethod   func() error
}

func (ts testSender) Down() error {
	return ts.downMethod()
}

func (ts testSender) Up() error {
	return ts.upMethod()
}

type testSpeedProvider struct{}

func (ts testSpeedProvider) Speed() int {
	return 40 // Make the tests quick
}

// pollSendQueueUntilEmpty blocks until k.SendQueueIsEmpty() returns true
// or the timeout is reached.  Returns true if the timeout was reached.
func pollSendQueueUntilEmpty(k *Keyer, timeout time.Duration) bool {
	done := make(chan interface{})
	go func() {
		// Poll over a short interval.  3ms should be okay.
		for range time.Tick(time.Millisecond * 3) {
			if k.SendQueueIsEmpty() {
				close(done)
				return
			}
		}
	}()

	select {
	case <-done:
		return false
	case <-time.After(timeout):
		return true
	}
}

func TestProcessSendQueueErrorBehavior(t *testing.T) {
	downErrorKeyer := New(testSpeedProvider{}, testSender{
		downMethod: func() error { return errors.New("test error") },
		upMethod:   func() error { return nil },
	})
	upErrorKeyer := New(testSpeedProvider{}, testSender{
		downMethod: func() error { return nil },
		upMethod:   func() error { return errors.New("test error") },
	})

	tests := []struct {
		id        int // To tell which test failed
		keyer     *Keyer
		sendInput func(*Keyer) // Wrapping for type reasons
	}{
		{id: 1, keyer: downErrorKeyer, sendInput: func(k *Keyer) { k.QueueMessage("A") }},
		{id: 2, keyer: downErrorKeyer, sendInput: func(k *Keyer) { k.QueueRune('A') }},
		{id: 3, keyer: upErrorKeyer, sendInput: func(k *Keyer) { k.QueueMessage("A") }},
		{id: 4, keyer: upErrorKeyer, sendInput: func(k *Keyer) { k.QueueRune('A') }},
	}

	for _, test := range tests {
		test.sendInput(test.keyer)
		err := test.keyer.ProcessSendQueue(true)
		if err == nil {
			t.Errorf("got nil, want error for test case %d", test.id)
		}
	}
}

func TestProcessSendQueueEmptyQueueBehavior(t *testing.T) {
	keyer := New(testSpeedProvider{}, testSender{
		downMethod: func() error { return nil },
		upMethod:   func() error { return nil },
	})

	done := make(chan interface{})
	go func() {
		keyer.ProcessSendQueue(true)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Millisecond * 100): // 100ms timeout, should be immediate.
		t.Error("got timeout, expected immediate return with true argument")
	}

	done = make(chan interface{})
	go func() {
		keyer.ProcessSendQueue(false)
		close(done)
	}()
	select {
	case <-done:
		t.Error("got return before timeout, expected timeout with false argument")
	case <-time.After(time.Millisecond * 100): // Same timeout as before
	}
}

func TestSendQueueIsEmpty(t *testing.T) {
	keyer := New(testSpeedProvider{}, testSender{
		downMethod: func() error { return nil },
		upMethod:   func() error { return nil },
	})

	if !keyer.SendQueueIsEmpty() {
		t.Error("got false, want true before queuing message")
	}

	keyer.QueueMessage("Test Message")
	if keyer.SendQueueIsEmpty() {
		t.Error("got true, want false after queuing message")
	}
}

func TestQueueMessage(t *testing.T) {
	tests := []struct {
		input           string
		keyEventsWanted int  // Total number of key up and down events wanted
		errorWanted     bool // If an error should be returned
	}{
		{input: "PARIS", keyEventsWanted: 28, errorWanted: false},
		{input: "PAR!S", keyEventsWanted: 0, errorWanted: true},
	}

	for _, test := range tests {
		var eventCount int
		keyFunc := func() error {
			eventCount++
			return nil
		}

		keyer := New(testSpeedProvider{}, testSender{
			downMethod: keyFunc,
			upMethod:   keyFunc,
		})
		err := keyer.QueueMessage(test.input)
		go keyer.ProcessSendQueue(true)

		if (err != nil) && !test.errorWanted {
			t.Errorf("got error, want nil queuing message %q", test.input)
		}
		if (err == nil) && test.errorWanted {
			t.Errorf("got nil, want error queuing message %q", test.input)
		}

		timedOut := pollSendQueueUntilEmpty(keyer, time.Second*3) // "PARIS" should take about half the timeout
		if timedOut {
			t.Errorf("timeout reached waiting on message %q to send", test.input)
		}
		if eventCount != test.keyEventsWanted {
			t.Errorf("got %d, want %d CW events for message %q",
				eventCount, test.keyEventsWanted, test.input)
		}
	}
}

func TestQueueRune(t *testing.T) {
	tests := []struct {
		input           rune
		keyEventsWanted int  // Total number of key up and down events wanted
		errorWanted     bool // If an error should be returned
	}{
		{input: 'E', keyEventsWanted: 2, errorWanted: false},
		{input: 'N', keyEventsWanted: 4, errorWanted: false},
		{input: 'R', keyEventsWanted: 6, errorWanted: false},
		{input: '$', keyEventsWanted: 0, errorWanted: true},
	}

	for _, test := range tests {
		var eventCount int
		keyFunc := func() error {
			eventCount++
			return nil
		}

		keyer := New(testSpeedProvider{}, testSender{
			downMethod: keyFunc,
			upMethod:   keyFunc,
		})
		err := keyer.QueueRune(test.input)
		go keyer.ProcessSendQueue(true)

		if (err != nil) && !test.errorWanted {
			t.Errorf("got error, want nil queuing rune %q", test.input)
		}
		if (err == nil) && test.errorWanted {
			t.Errorf("got nil, want error queuing rune %q", test.input)
		}

		timedOut := pollSendQueueUntilEmpty(keyer, time.Second*2)
		if timedOut {
			t.Errorf("timeout reached waiting on rune %q to send", test.input)
		}
		if eventCount != test.keyEventsWanted {
			t.Errorf("got %d, want %d CW events for rune %q",
				eventCount, test.keyEventsWanted, test.input)
		}
	}
}

func TestDrainSendQueue(t *testing.T) {
	keyer := New(testSpeedProvider{}, testSender{
		downMethod: func() error { return nil },
		upMethod:   func() error { return nil },
	})
	keyer.QueueMessage("This message would take a while to send")
	go keyer.ProcessSendQueue(true)

	keyer.DrainSendQueue()
	if !keyer.SendQueueIsEmpty() {
		t.Error("send queue not empty after DrainSendQueue call")
	}
}
