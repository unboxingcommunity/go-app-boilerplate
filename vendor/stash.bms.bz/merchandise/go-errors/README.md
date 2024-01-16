//EXAMPLE 

package main

import (
	"stash.bms.bz/merchandise/go-errors"
	"stash.bms.bz/merchandise/go-logger"
)

func init() {
	er := errors.Parse(getErrors() // Function should return map string interface of error objects)
	if er != nil {
		panic(er)
	}
}

func main() {
	log := logger.New("go-errors")
	err := DoSomethingAwesome()

	fmt.Println(errors.GetReference(err))
	//Add some references
	errors.AddReference(err, map[string]interface{}{"id": "toooo awesome"})
	// // Get the references
	ref := errors.GetReference(err)
	fmt.Println(ref)

	// Log it
	errors.Log(err)

	// Print the stack trace
	errors.StackTrace(err)

	// Get the stack trace in the string
	s := errors.SstackTrace(err)
	fmt.Println(s)

	// Get the Details of the error
	// Need a better name for details
	details := errors.Details(err)
	// fmt.Println(details)

	log.Error(details.Code, details.Description, details.Severity, errors.GetErrorObj(err))
}

func DoSomethingAwesome() error {
	err := errors.New(PARAMETER_MISSING)
	return err
}

