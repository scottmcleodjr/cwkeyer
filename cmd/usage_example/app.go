package main

import (
	"log"
	"time"

	"github.com/scottmcleodjr/cwkeyer"
)

type config struct {
	speed int
}

func (c *config) Speed() int {
	return c.speed
}

func main() {

	// We'll need a cwkeyer.Key.  This type abstracts the open and close
	// functionality of a straight key into two methods.  There are two
	// keys in the library.  This one makes beeping sounds.  The other one
	// sets the DTR bit on a serial port.

	key, err := cwkeyer.NewBeepKey(700, 48000, 1200)
	if err != nil {
		log.Fatal(err)
	}

	// We'll need a type with the method Speed() for the cwkeyer.Keyer.
	// This will be called before every dit or dah so that the speed can
	// be adjusted while sending a message.

	cfg := config{speed: 18}

	// We can create the cwkeyer.Keyer by passing in the SpeedProvder (cfg)
	// and the key.

	keyer := cwkeyer.New(&cfg, key)

	// We will start the keyer in a new goroutine.  The ProcessSendQueue
	// method starts the processing of the send queue and will return if
	// an error is returned by the cwkeyer.Key.  ProcessSendQueue may also
	// return when the send queue is empty, depending on its argument.

	go func() {
		err := keyer.ProcessSendQueue(false)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// We add a message to the send queue.  This call will return before the
	// message completes sending.

	err = keyer.QueueMessage("CQ CQ CQ de K3GDS K3GDS KN")
	if err != nil {
		log.Fatal(err)
	}

	// After 3 seconds, we will increase the speed to 25wpm.  You will hear this
	// change happen in real time.

	time.Sleep(time.Second * 3)
	cfg.speed = 25

	// After another 3 seconds, we will interrupt the current message.  This will
	// stop any output immediately and drain the send queue.

	time.Sleep(time.Second * 3)
	keyer.DrainSendQueue()
}
