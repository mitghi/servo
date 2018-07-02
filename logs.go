package servo

import (
	"errors"
	"log"
	"os"
)

const (
	// cLogPATH is the default logger output path
	cLogPATH        string = "/tmp/servo.out"
	cDefaultLPrefix string = "servo"
)

// LogFile is the struct that wraps
// `log.Log` struct.
type LogFile struct {
	file          *os.File
	opts          LogOpts
	isInitialized bool
	isClosed      bool
	Logger        *log.Logger
}

// LogOpts provides neccessary options to setup
// `LogFile`.
type LogOpts struct {
	// TODO:
	// . specify log level
	output string
	level  int
}

// - MARK: Alloc/Init section.

// NewLog allocates and initializes a new
// `LogFile` and return a pointer to it.
func NewLog() *LogFile {
	return &LogFile{
		file:          nil,
		isInitialized: false,
		isClosed:      false,
		opts: LogOpts{
			output: cLogPATH,
		},
	}
}

// NewLogFromOpts allocates and initializes
// a new `LogFile` and configs it accordingly
// using options supplied as `opts` argument.
func NewLogFromOpts(opts LogOpts) (lf *LogFile) {
	lf = NewLog()
	lf.opts = opts
	return lf
}

// create either opens or creates a file with
// the given path and returns a file pointer, it
// sets state flags accordingly.
func (lf *LogFile) create(path string) *os.File {
	var (
		file *os.File
	)
	file = FOpenOrCreate(path)
	if file == nil {
		goto ERROR
	}
	lf.isInitialized = true
	lf.file = file
	return file
ERROR:
	lf.isInitialized = false
	lf.file = nil
	return nil
}

// Setup performs internal configuration such as
// opening the file path and return an error to
// indicate any possible failure.
func (lf *LogFile) Setup() (err error) {
	if lf.file == nil {
		/* d e b u g */
		// log.Println("(Setup) lf.file (before initialization): ->", lf.file, lf.isInitialized, lf)
		/* d e b u g */
		lf.file = lf.create(lf.opts.output)
		lf.Logger = log.New(lf.file, cDefaultLPrefix, log.LstdFlags)
		lf.isInitialized = true
		log.SetOutput(lf.file)
		return nil
	}
	return errors.New("logfile: failed to setup LogFile.")
}

// Close closes the log file and returns
// true in case of success.
func (lf *LogFile) Close() bool {
	var (
		err error
	)
	if lf.file != nil && !lf.isClosed {
		return false
	}
	err = lf.file.Close()
	if err != nil {
		return false
	}
	lf.isClosed = true
	return true
}

// Writes writes `input` into log stream and returns
// a boolean to indicate its success.
func (lf *LogFile) Write(input string) (ok bool) {
	if lf.isInitialized && !lf.isClosed {
		lf.Logger.Println(input)
	}
	return false
}

// GetDefault is a factory receiver method
// that allocates and initializes a new
// `LogFile` struct with default options
// and returns a pointer to it.
func GetDefault() *LogFile {
	var (
		l *LogFile = NewLog()
	)
	// discard error value, because default output
	// file is either creates or opened.
	_ = l.Setup()
	return l
}
