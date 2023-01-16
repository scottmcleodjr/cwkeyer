// CWKeyer is a library for sending morse code (CW).
// The library uses an asynchronous send queue that allows
// the caller to adjust the speed, stop a message, or send
// additional messages while a previous message is still
// being keyed.
package cwkeyer

import (
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
	"go.bug.st/serial"
)

// 2048 should be more than enough for any reasonable message
const eventChanLength = 2048

type SerialDTRKey struct {
	port serial.Port
}

// NewSerialDTRKey creates a SerialDTRKey on the specified port.
// Suggested values based on testing: baudrate=115200.
func NewSerialDTRKey(portName string, baudrate int) (*SerialDTRKey, error) {
	port, err := serial.Open(portName, &serial.Mode{BaudRate: baudrate})
	if err != nil {
		return &SerialDTRKey{}, err
	}
	key := SerialDTRKey{port: port}
	return &key, key.Up() // Always start in the up position
}

func (s *SerialDTRKey) ClosePort() error {
	return s.port.Close()
}

func (s *SerialDTRKey) Down() error {
	return s.port.SetDTR(true)
}

func (s *SerialDTRKey) Up() error {
	return s.port.SetDTR(false)
}

type BeepKey struct {
	streamer beep.Streamer
}

// NewBeepKey creates a BeepKey.
// Suggested values based on testing: freq=700, sampleRate=48000, bufferSize=1200.
func NewBeepKey(freq, sampleRate, bufferSize int) (*BeepKey, error) {
	speaker.Init(beep.SampleRate(sampleRate), bufferSize)
	s, err := generators.SinTone(beep.SampleRate(sampleRate), freq)
	if err != nil {
		return &BeepKey{}, err
	}
	return &BeepKey{streamer: s}, nil
}

func (b *BeepKey) CloseSpeaker() {
	speaker.Close()
}

func (b *BeepKey) Down() error {
	speaker.Play(b.streamer)
	return nil
}

func (b *BeepKey) Up() error {
	speaker.Clear()
	return nil
}

// SpeedProvider provides the speed each time a CW event is keyed.
type SpeedProvider interface {
	Speed() int
}

// Key abstracts a CW key into a type with methods
// called when the key is circuit opened or closed.
type Key interface {
	Down() error // Called when the keyer keys down
	Up() error   // Called when the keyer keys up
}

type Keyer struct {
	key       Key
	speed     SpeedProvider
	sendQueue chan event
}

// New creates a Keyer.
func New(speed SpeedProvider, key Key) Keyer {
	return Keyer{
		speed:     speed,
		key:       key,
		sendQueue: make(chan event, eventChanLength),
	}
}

// ProcessSendQueue processes the send queue.  A bool is accepted
// to indicate if ProcessSendQueue should return when the queue is
// empty.  A method error from the Key will always cause a return.
func (k Keyer) ProcessSendQueue(returnOnEmptyQueue bool) error {
	for {
		if returnOnEmptyQueue && k.SendQueueIsEmpty() {
			return nil
		}
		e := <-k.sendQueue // Blocks loop when chan is empty
		err := k.keyEvent(e)
		if err != nil {
			return err
		}
	}
}

// SendQueueIsEmpty returns true when there is nothing waiting to be sent.
func (k Keyer) SendQueueIsEmpty() bool {
	return len(k.sendQueue) == 0
}

// QueueMessage adds a string to the send queue.  An error is returned
// if the string contains an unsupported rune.
func (k Keyer) QueueMessage(message string) error {
	// Check IsKeyable before any are put on the queue
	for _, r := range message {
		if !IsKeyable(r) {
			return fmt.Errorf("unsupported character: %q", r)
		}
	}
	for _, r := range message {
		k.QueueRune(r)
	}
	return nil
}

// QueueRune adds a rune to the send queue.  An error is returned
// if the rune is unsupported.
func (k Keyer) QueueRune(r rune) error {
	events, present := events(r)
	if !present {
		return fmt.Errorf("unsupported character: %q", r)
	}
	for _, e := range events {
		k.sendQueue <- e
	}
	return nil
}

// DrainSendQueue interrupts the current message by
// draining the send queue, returning when it is empty.
func (k Keyer) DrainSendQueue() {
	for {
		select {
		case <-k.sendQueue:
		default:
			return
		}
	}
}

func (k Keyer) keyEvent(e event) error {
	if e == dit || e == dah {
		err := k.key.Down()
		if err != nil {
			return err
		}
	}
	time.Sleep(eventLength(e, k.speed.Speed()))
	if e == dit || e == dah {
		err := k.key.Up()
		if err != nil {
			return err
		}
	}
	time.Sleep(eventLength(space, k.speed.Speed()))
	return nil
}
