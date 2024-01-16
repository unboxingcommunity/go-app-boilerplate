package errors

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/mitchellh/mapstructure"
	pkgErrors "github.com/pkg/errors"
)

var errorDetails map[string]*ErrorDetails

//Error ...
type Error struct {
	errorVal     error
	ErrorDetails ErrorDetails
}

//ErrorDetails ...
type ErrorDetails struct {
	Code           *string                `json:"code,omitempty"`
	UserCode       *string                `json:"userCode,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Severity       *string                `json:"severity,omitempty"`
	Message        *string                `json:"message,omitempty"`
	Source         *ErrorSource           `json:"source,omitempty"`
	Action         *string                `json:"action,omitempty"`
	HTTPStatusCode *int                   `json:"httpStatusCode,omitempty"`
	Reference      map[string]interface{} `json:"reference,omitempty"`
}

//EDetails ... used for logging details without the hassle of pointers
type EDetails struct {
	Code           string
	UserCode       string
	Description    string
	Severity       string
	Message        string
	Action         string
	HTTPStatusCode int
	Source         *ErrorSource
	Reference      map[string]interface{}
}

//ErrorSource ...
type ErrorSource struct {
	Caller      string      `json:"caller,omitempty"`
	File        string      `json:"file,omitempty"`
	Line        int         `json:"line,omitempty"`
	ErrorString string      `json:"errorString,omitempty"`
	ErrorStruct interface{} `json:"errorStruct,omitempty"`
	StackTrace  string      `json:"stackTrace,omitempty"`
}

//New ...
func New(err string) error {
	errorDetail := &ErrorDetails{}
	errorSource := &ErrorSource{}
	if errorDetails[err] != nil {
		errorDetail = errorDetails[err]
	}

	// If message is empty , set default message
	if errorDetail.Message == nil || *errorDetail.Message == "" {
		message := "Oops, something went wrong"
		errorDetail.Message = &message
	}

	//Gets error origin path, line , function
	pc, file, line, _ := runtime.Caller(1)
	fn := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	// Traverses one call frame up the stack
	previous := strings.Split(fn[len(fn)-2], "/")
	if strings.ToLower(fn[len(fn)-1]) == "wrap" && strings.ToLower(previous[len(previous)-1]) == "go-errors" {
		pc, file, line, _ = runtime.Caller(2)
		fn = strings.Split(runtime.FuncForPC(pc).Name(), ".")
	}

	//Forms Pkg Erros object
	pkgErr := pkgErrors.New(err)

	//Creates error source object
	errorSource.StackTrace = fmt.Sprintf("%+v", pkgErr)
	errorSource.File = file
	errorSource.Line = line
	errorSource.Caller = fn[len(fn)-1]

	//Assigns error source to error details
	errorDetail.Source = errorSource

	//Creates Error Object
	my := Error{errorVal: pkgErr, ErrorDetails: *errorDetail}

	return &my
}

//Error ...
func (e *Error) Error() string {
	return e.errorVal.Error()
}

//ErrorAsMap ...
func ErrorAsMap(e error) (errMap map[string]interface{}) {
	er, ok := e.(*Error)
	if ok == false {
		return
	}
	errMap = map[string]interface{}{
		"errorDetails": er.ErrorDetails,
	}
	return errMap
}

//FormResponse ...
func FormResponse(e error) (response *ErrorDetails) {
	_, ok := e.(*Error)
	if ok == false {
		return
	}
	er := *(e.(*Error))
	//Gets rid of unnecessary fields
	er.ErrorDetails.Code = nil
	er.ErrorDetails.Severity = nil
	er.ErrorDetails.Source = nil
	er.ErrorDetails.HTTPStatusCode = nil
	er.ErrorDetails.Reference = nil
	er.ErrorDetails.Description = nil

	//Returns formed response
	return &er.ErrorDetails
}

//SstackTrace ...
func SstackTrace(e error) string {
	er := e.(*Error)
	return fmt.Sprintf("%+v", er.errorVal)
}

//Details ...
func Details(e error) (details EDetails) {
	er, ok := e.(*Error)
	if ok == false {
		details.Description = e.Error()
		return details
	}

	eDetails := er.ErrorDetails

	// Assigns details as non pointers to custom EDetails struct
	if eDetails.Action != nil {
		details.Action = *eDetails.Action
	}

	if eDetails.Code != nil {
		details.Code = *eDetails.Code
	}

	if eDetails.UserCode != nil {
		details.UserCode = *eDetails.UserCode
	}

	if eDetails.Description != nil {
		details.Description = *eDetails.Description
	}

	if eDetails.Message != nil {
		details.Message = *eDetails.Message
	}

	if eDetails.Severity != nil {
		details.Severity = *eDetails.Severity
	}

	if eDetails.HTTPStatusCode != nil {
		details.HTTPStatusCode = *eDetails.HTTPStatusCode
	}

	if eDetails.Reference != nil {
		details.Reference = eDetails.Reference
	}

	details.Source = eDetails.Source

	return details
}

//StackTrace ...
func StackTrace(e error) {
	er := e.(*Error)
	fmt.Printf("%+v", er.errorVal)
}

//AddReference ...
func AddReference(e error, code string, reference map[string]interface{}, callerLevel ...int) {
	er, ok := e.(*Error)
	if ok == false {
		return
	}

	if code != "" {
		var cLevel = 1
		//Assigns caller level
		if len(callerLevel) != 0 {
			cLevel = callerLevel[0]
		}

		//Gets reference origin function
		pc, _, _, _ := runtime.Caller(cLevel)
		fn := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcName := fn[len(fn)-1]

		key := code + "." + funcName
		if er.ErrorDetails.Reference == nil {
			er.ErrorDetails.Reference = map[string]interface{}{
				key: reference,
			}
		} else {
			er.ErrorDetails.Reference[key] = reference
		}
	} else {
		er.ErrorDetails.Reference = reference
	}
}

//Add ...
func Add(customError error, sourceError error) {
	er, ok := customError.(*Error)
	if ok == false {
		return
	}

	er.ErrorDetails.Source.ErrorString = sourceError.Error()
	er.ErrorDetails.Source.ErrorStruct = sourceError
}

//Source ... gets source object as a map[string]interface{}
func Source(e error) (source map[string]interface{}) {
	// Validate if error is of type Error
	er, ok := e.(*Error)
	if ok == false {
		return nil
	}
	err := mapstructure.Decode(er.ErrorDetails.Source, &source)
	if err != nil {
		return nil
	}

	return source
}

// //Reference ...
// func Reference(e error) interface{} {
// 	// Validate if error is of type Error
// 	er, ok := e.(*Error)
// 	if ok == false {
// 		return ""
// 	}
// 	return er.ref
// }

//Log ...
func Log(e error) string {
	er, ok := e.(*Error)
	if ok == false {
		return ""
	}

	fmt.Println(er.ErrorDetails)
	// Return a JSON string as output
	// Which can later be logged
	return ""
}

//Map ... maps errMap values into existing err object
func Map(e error, errMap map[string]interface{}) {
	er, ok := e.(*Error)
	if ok == false {
		return
	}

	if errMap == nil {
		return
	}

	var errDetails *ErrorDetails
	err := mapstructure.Decode(errMap, &errDetails)
	if err != nil {
		return
	}

	if errDetails.Action != nil {
		er.ErrorDetails.Action = errDetails.Action
	}

	if errDetails.Code != nil {
		er.ErrorDetails.Code = errDetails.Code
	}

	if errDetails.Description != nil {
		er.ErrorDetails.Description = errDetails.Description
	}

	if errDetails.Message != nil {
		er.ErrorDetails.Message = errDetails.Message
	}

	if errDetails.Severity != nil {
		er.ErrorDetails.Severity = errDetails.Severity
	}

	if errDetails.HTTPStatusCode != nil {
		er.ErrorDetails.HTTPStatusCode = errDetails.HTTPStatusCode
	}

}

//Wrap ...
func Wrap(err string, e error) error {
	er := New(err)
	Add(er, e)
	return er
}

//Parse ... Creates global error map
func Parse(b map[string]interface{}) error {
	if errorDetails == nil {
		errorDetails = make(map[string]*ErrorDetails)
	}

	var errDetails map[string]*ErrorDetails
	err := mapstructure.Decode(b, &errDetails)
	if err != nil {
		return err
	}

	for key, value := range errDetails {
		if _, ok := errDetails[key]; ok {
			errorDetails[key] = value
		}
	}
	return nil
}
