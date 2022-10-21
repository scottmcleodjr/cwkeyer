# cwkeyer

[![Go Report Card](https://goreportcard.com/badge/github.com/scottmcleodjr/cwkeyer)](https://goreportcard.com/report/github.com/scottmcleodjr/cwkeyer)
[![License](https://img.shields.io/badge/License-BSD_2--Clause-blue.svg)](LICENSE)

CWKeyer is a library for sending morse code (CW) in Go.  The library uses an asynchronous send queue that allows the caller to adjust the speed, stop a message, or send additional messages while a previous message is still being keyed.

This was a part of a larger, unfinished project.  I wasn't getting around to finishing the whole thing, but this part of the application was mostly complete and could be a standalone library for use elsewhere.

## Usage Example & Documentation

The [usage example](cmd/usage_example/app.go) demonstrates the main functionalities.  The documentation comments are fairly complete - Look to the godoc for everything else you can do.