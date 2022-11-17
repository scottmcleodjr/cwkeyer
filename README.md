# **cwkeyer**

[![Go Report Card](https://goreportcard.com/badge/github.com/scottmcleodjr/cwkeyer)](https://goreportcard.com/report/github.com/scottmcleodjr/cwkeyer)
[![License](https://img.shields.io/badge/License-BSD_2--Clause-blue.svg)](LICENSE)

CWKeyer is a library for sending morse code (CW) in Go.  The library uses an asynchronous send queue that allows the caller to adjust the speed, stop a message, or send additional messages while a previous message is still being keyed.

### **Interface Compatibility**

At the moment, the library has a Key that beeps and a Key that sets the DTR signal on a serial port.  The latter one does everything I need to interface with my radios.  If you need something else, make an issue or pull request.  I'm happy to make this useful to more people.


### **Usage Example & Documentation**

The [usage example](cmd/usage_example/app.go) demonstrates the main functionalities.  The documentation comments are fairly complete - Look to the [documentation](https://pkg.go.dev/github.com/scottmcleodjr/cwkeyer) for everything else you can do.

This library is also used in the [K3GDS REKL](https://github.com/scottmcleodjr/rekl).