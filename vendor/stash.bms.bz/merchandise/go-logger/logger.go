package logger

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
)

// Logger ...
type Logger struct {
	Log    *Log
	config *Config
	masker *Masker
}

// ErrorSource ...
type ErrorSource struct {
	Caller      *string     `json:"caller,omitempty" valid:"required"`
	File        *string     `json:"file,omitempty" valid:"required"`
	Line        *int        `json:"line,omitempty" valid:"required"`
	StackTrace  *string     `json:"stackTrace,omitempty" valid:"required"`
	ErrorString *string     `json:"errorString,omitempty"`
	ErrorSource interface{} `json:"errorSource,omitempty"`
}

// Log ...
type Log struct {
	Host        string      `json:"host,omitempty"`
	App         string      `json:"app,omitempty"`
	Description string      `json:"description,omitempty"`
	Code        string      `json:"code,omitempty"`
	Severity    string      `json:"severity,omitempty"`
	Type        string      `json:"type,omitempty"`
	File        string      `json:"file,omitempty"`
	Line        int         `json:"line,omitempty"`
	Caller      string      `json:"caller,omitempty"`
	CallerLevel int         `json:"callerLevel,omitempty"`
	Source      interface{} `json:"source,omitempty"`
	Reference   interface{} `json:"reference,omitempty"`
	Ts          string      `json:"ts,omitempty"`
}

// Config ...
type Config struct {
	CallerLevel *int   `json:"callerLevel,omitempty"`
	Debug       *bool  `json:"debug,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

// Masker ...
type Masker struct {
	Enabled    bool
	PatternMap map[string]string
}

// GLog ... global logger object
var GLog *Logger

////////////////////////////////////////////////////////////

//----Exposed Functions-----------//

// New ... returns new instance of Logger
func New(app string, callerLevel ...int) (log *Logger) {
	var logger Logger
	logger.Log = new(Log)
	logger.config = new(Config)
	logger.config.Reference = "object"

	name := os.Getenv("LOG_APP_NAME")
	if name == "" {
		name = app
		err := os.Setenv("LOG_APP_NAME", name)
		if err != nil {
			defaultLog("Error setting LOG_APP_NAME environment variable ", "", err)
		}
	}
	logger.Log.App = name

	host, err := os.Hostname()
	if err != nil {
		logger.Log.Host = "Could not register host. " + "LOGGER ERR @Func 'New' : " + err.Error()
	} else {
		logger.Log.Host = host
	}

	if len(callerLevel) != 0 {
		logger.Log.CallerLevel = callerLevel[0]
	} else {
		logger.Log.CallerLevel = 1
	}

	//Returns
	return &logger
}

//Global ... assignes supplied logger instance to global logger object (GLog)
func Global(logger *Logger) {
	GLog = logger
}

// Info ...
func (logger *Logger) Info(description string, reference ...interface{}) {
	if reference != nil {
		err := validateReference(reference)
		if err != nil {
			reference[0] = map[string]interface{}{"error": "Reference not displayed. " + err.Error()}
		}
	}

	var tmpLog Logger
	tmpLog.Log = new(Log)
	tmpLog.Log.Description = description
	tmpLog.Log.CallerLevel = logger.Log.CallerLevel
	tmpLog.Log.Type = "info"
	tmpLog.Log.App = logger.Log.App
	tmpLog.Log.Host = logger.Log.Host
	tmpLog.config = logger.config
	tmpLog.masker = logger.masker

	//Forms new Log objects
	clog := newLog(&tmpLog, reference)

	//Sets log output path and flags
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// Creates Json object
	body, err := json.Marshal(clog)
	if err != nil {
		defaultLog(description, "LOGGER ERR @Func 'Info - 2' : ", err)
		return
	}

	// Logs json object
	log.Println(string(body))

}

// Debug ...
func (logger *Logger) Debug(description string, reference ...interface{}) {
	if logger.config == nil || logger.config.Debug == nil || *logger.config.Debug == false {
		return
	}

	if reference != nil {
		err := validateReference(reference)
		if err != nil {
			reference[0] = map[string]interface{}{"error": "Reference not displayed. " + err.Error()}
		}
	}

	var tmpLog Logger
	tmpLog.Log = new(Log)
	tmpLog.Log.Description = description
	tmpLog.Log.CallerLevel = logger.Log.CallerLevel
	tmpLog.Log.Type = "debug"
	tmpLog.Log.App = logger.Log.App
	tmpLog.Log.Host = logger.Log.Host
	tmpLog.config = logger.config
	tmpLog.masker = logger.masker

	// Forms new Log objects
	clog := newLog(&tmpLog, reference)

	// Sets log output path
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// Creates Json object
	body, err := json.Marshal(clog)
	if err != nil {
		defaultLog(description, "LOGGER ERR @Func 'Debug - 2' : ", err)
		return
	}

	// Logs json object
	log.Println(string(body))
}

// Warn ...
func (logger *Logger) Warn(description string, reference ...interface{}) {
	if reference != nil {
		err := validateReference(reference)
		if err != nil {
			reference[0] = map[string]interface{}{"error": "Reference not displayed. " + err.Error()}
		}
	}

	var tmpLog Logger
	tmpLog.Log = new(Log)
	tmpLog.Log.Description = description
	tmpLog.Log.CallerLevel = logger.Log.CallerLevel
	tmpLog.Log.Type = "warning"
	tmpLog.Log.App = logger.Log.App
	tmpLog.Log.Host = logger.Log.Host
	tmpLog.config = logger.config
	tmpLog.masker = logger.masker

	clog := newLog(&tmpLog, reference)

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	body, err := json.Marshal(clog)
	if err != nil {
		defaultLog(description, "LOGGER ERR @Func 'Warn - 2' : ", err)
		return
	}

	log.Println(string(body))
}

// Error ...
func (logger *Logger) Error(code string, description string, severity string, source interface{}, reference ...interface{}) {
	if reference != nil {
		err := validateReference(reference)
		if err != nil {
			reference[0] = map[string]interface{}{"error": "Reference not displayed. " + err.Error()}
		}
	}

	var errSrc *ErrorSource
	var errorSource interface{}
	var sourcePresent = true
	err := mapstructure.Decode(source, &errSrc)
	if err != nil {
		sourcePresent = false
		errorSource = "Source not displayed. " + err.Error()
	}

	if errSrc != nil && sourcePresent == true {
		_, err = validator.ValidateStruct(*errSrc)
		if err != nil {
			errorSource = "Source not displayed. " + err.Error()
		}
		errorSource = errSrc
	}

	var tmpLog Logger
	tmpLog.Log = new(Log)
	tmpLog.Log.Description = description
	tmpLog.Log.CallerLevel = logger.Log.CallerLevel
	tmpLog.Log.Type = "error"
	tmpLog.Log.Code = code
	tmpLog.Log.Severity = severity
	tmpLog.Log.App = logger.Log.App
	tmpLog.Log.Host = logger.Log.Host
	tmpLog.config = logger.config
	tmpLog.masker = logger.masker

	//Assign source to tmpLog
	tmpLog.Log.Source = errorSource

	clog := newLog(&tmpLog, reference)

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	body, err := json.Marshal(clog)
	if err != nil {
		defaultLog(description, "LOGGER ERR @Func 'Error - 4' : ", err)
		return
	}

	log.Print(string(body))
}

// Fatal ... app exits on logging as fatal
func (logger *Logger) Fatal(code string, description string, source interface{}, reference ...interface{}) {
	if reference != nil {
		err := validateReference(reference)
		if err != nil {
			reference[0] = map[string]interface{}{"error": "Reference not displayed. " + err.Error()}
		}
	}

	var errSrc *ErrorSource
	var errorSource interface{}
	var sourcePresent = true
	err := mapstructure.Decode(source, &errSrc)
	if err != nil {
		sourcePresent = false
		errorSource = "Source not displayed. " + err.Error()
	}

	if errSrc != nil && sourcePresent == true {
		_, err = validator.ValidateStruct(*errSrc)
		if err != nil {
			errorSource = "Source not displayed. " + err.Error()
		}
		errorSource = errSrc
	}

	var tmpLog Logger
	tmpLog.Log = new(Log)
	tmpLog.Log.Description = description
	tmpLog.Log.CallerLevel = logger.Log.CallerLevel
	tmpLog.Log.Type = "fatal"
	tmpLog.Log.Code = code
	tmpLog.Log.App = logger.Log.App
	tmpLog.Log.Host = logger.Log.Host
	tmpLog.config = logger.config
	tmpLog.masker = logger.masker

	// Assign source to tmpLog
	tmpLog.Log.Source = errorSource

	clog := newLog(&tmpLog, reference)

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	body, err := json.Marshal(clog)
	if err != nil {
		defaultLog(description, "LOGGER ERR @Func 'Fatal - 4' : ", err)
		return
	}

	log.Println(string(body))
	os.Exit(500)
}

// EnableMask ... enables masking of patterns present inside pattern map
func (logger *Logger) EnableMask(patternMap map[string]string) {
	if len(patternMap) > 0 {
		logger.masker = &Masker{Enabled: true, PatternMap: patternMap}
	}
}

// Set ... sets reference and severity manually
func (logger *Logger) Set(config string) {
	var cnf Config

	err := json.Unmarshal([]byte(config), &cnf)
	if err != nil {
		defaultLog("Set func error", "LOGGER ERR  @Func Set : ", err)
		return
	}

	// Assigns config to master Logger instance
	logger.config = &Config{CallerLevel: getCallerVar(cnf, logger), Debug: getDebugVar(cnf, logger), Reference: getReferenceVar(cnf, logger)}
}

func getCallerVar(source1 Config, source2 *Logger) *int {
	var callerLevel int
	if source1.CallerLevel != nil {
		callerLevel = *source1.CallerLevel
		return &callerLevel
	} else if source2.config.CallerLevel != nil {
		callerLevel = *source2.config.CallerLevel
		return &callerLevel
	}
	return nil
}

func getDebugVar(source1 Config, source2 *Logger) *bool {
	var debug bool
	if source1.Debug != nil {
		debug = *source1.Debug
		return &debug
	} else if source2.config.Debug != nil {
		debug = *source2.config.Debug
		return &debug
	}
	return nil
}

func getReferenceVar(source1 Config, source2 *Logger) string {
	var reference string
	if source1.Reference != "" {
		reference = source1.Reference
		return reference
	} else if source2.config.Reference != "" {
		reference = source2.config.Reference
		return reference
	}
	return ""
}

//---Internal Functions---//

// Assigns data and creates new log object
func newLog(logger *Logger, reference []interface{}) (clog Log) {
	clog.Description = logger.Log.Description
	clog.Code = logger.Log.Code
	clog.App = logger.Log.App
	clog.Host = logger.Log.Host

	// Default Caller level
	clog.CallerLevel = logger.Log.CallerLevel
	clog.Severity = logger.Log.Severity
	clog.Type = logger.Log.Type

	if logger.config != nil && logger.config.CallerLevel != nil {
		clog.CallerLevel = *logger.config.CallerLevel
	}

	if logger.Log.Source != nil {
		clog.Source = logger.Log.Source
	}

	if reference != nil {
		if logger.config.Reference == "object" {
			clog.Reference = reference[0]
		} else if logger.config.Reference == "string" {
			body, err := json.Marshal(reference[0])
			if err != nil {
				clog.Reference = map[string]interface{}{"error": "Error marshaling reference object"}
			} else {
				clog.Reference = map[string]interface{}{"payload": maskData(logger, string(body))}
			}
		} else {
			clog.Reference = map[string]interface{}{"error": "Set invalid config value"}
		}
	}

	// Gets func name , file name and line no.
	caller, file, line := getFrame(clog.CallerLevel + 2)

	clog.Caller = caller
	clog.File = file
	clog.Line = line
	clog.Ts = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	return clog
}

// maskData ...
func maskData(logger *Logger, input string) (output string) {
	// Returns output as is if masking is not enabled
	if logger.masker == nil {
		return input
	}
	if !logger.masker.Enabled {
		return input
	}

	var re *regexp.Regexp
	output = input

	// Masks patterns
	for key, value := range logger.masker.PatternMap {
		re = regexp.MustCompile(key)
		output = re.ReplaceAllString(output, value)
	}

	return output
}

// defaultLog if everything fails
func defaultLog(description string, errString string, err error) {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	jsonStr := `{"description":"` + description + `","error":"` + errString + err.Error() + `","ts":"` + ts + `"}`

	log.Println(jsonStr)
}

// getFrame gets Frame Info from stack
func getFrame(skip int) (caller string, file string, line int) {
	pc, file, line, _ := runtime.Caller(skip)
	fn := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	return fn[len(fn)-1], file, line
}

// validateReference validates the type of reference to a struct or map
func validateReference(reference []interface{}) (err error) {
	if reference[0] != nil {
		if reflect.ValueOf(reference[0]).Kind() == reflect.Ptr {
			if reflect.ValueOf(reference[0]).Elem().Kind() == reflect.Struct || reflect.ValueOf(reference[0]).Elem().Kind() == reflect.Map {
				return nil
			}
		}
		if reflect.TypeOf(reference[0]).Kind() != reflect.Struct && reflect.TypeOf(reference[0]).Kind() != reflect.Map {
			return errors.New("Reference should either be a (map or a struct) / (pointer to a map or a pointer to a struct)")
		}
	}
	return nil
}
