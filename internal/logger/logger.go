package logger

import (
	"log"
	"os"
)

var Info *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var Err *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
