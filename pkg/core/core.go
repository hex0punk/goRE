package core

import (
	"math/rand"
	"os"
	"time"
)

var Debug = false
var Verbose = false
var CurrentDir, _ = os.Getwd()
var src = rand.NewSource(time.Now().UnixNano())