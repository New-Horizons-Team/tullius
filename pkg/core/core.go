package core

import (
	"math/rand"
	"os"
	"time"
)

// Debug puts Tullius into debug mode and displays debug messages
var Debug = false

// Verbose puts Tullius into verbose mode and displays verbose messages
var Verbose = false

// CurrentDir is the current directory where Merlin was executed from
var CurrentDir, _ = os.Getwd()
var src = rand.NewSource(time.Now().UnixNano())

// Constants
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

