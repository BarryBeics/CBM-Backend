package shared

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger is a function that configures logging for the application. It
// sets the output format to include caller information, parses command-line
// arguments, sets the log level based on the LOG_LEVEL environment variable
// or the -debug flag (default: Info), and prints the log level to the console
// for debugging purposes. The function uses the zerolog logging library and
// the standard flag package to handle command-line arguments.
func SetupLogger() {
	// Set output format to include caller information
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	// Parse command-line arguments for debug and trace flags
	debug := flag.Bool("debug", false, "enable debug logging")
	trace := flag.Bool("trace", false, "enable trace logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-debug] [-trace] [Report]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Set log level based on environment variable (default: Info)
	logLevel := zerolog.InfoLevel
	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if lvl, err := zerolog.ParseLevel(envLogLevel); err == nil {
			logLevel = lvl
		}
	} else if *debug {
		logLevel = zerolog.DebugLevel
	} else if *trace {
		logLevel = zerolog.TraceLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(logLevel)

}
