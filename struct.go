package zerror

import "github.com/znxlc/zerror/errormessage"

// this file contains types and variables used by the package

// Error flags and predefined codes
const (
	FlagReturnLastErrorElement  = "LAST"  // returns the last error in the list by default (Get() and Error() functions)
	FlagReturnFirstErrorElement = "FIRST" // returns the first element in the list when calling Get() and Error()
	FlagReturnErrorCode         = "CODE"  // default text returned when calling Error()
	FlagReturnErrorMsg          = "MSG"   // return Msg field when calling Error()
)

var (
	// ElementIndexReturned represents the default element to be returned when calling single ErrorElement functions like Get() and Error()
	ElementIndexReturned = FlagReturnFirstErrorElement
	// ElementTextReturned is used by Error() to select which text to return (default is Error Code)
	ElementTextReturned = FlagReturnErrorCode
	// ElementGenerator will be used to create new error elements and should be a pointer to the constructor of the errorElement used
	DefaultElementGenerator = errormessage.New
)

// ZError is the main error structure of the package
type ZError struct {
	ElementIndexReturned string                             `json:"-"` // set the default element to be returned when calling Get() or Error()
	ElementGenerator     errormessage.ErrorElementGenerator // the generator for the error elements (pointer to the New() constructor)
	Errors               []errormessage.ErrorElement        `json:"errors"` // the error list
}

type Error interface {
	Add(...any)
	Clear()
	Error() string
	GetList() []errormessage.ErrorElement
	Get(...int) errormessage.ErrorElement
	HasErrors() bool
	SetDefaultElementIndexReturned(string)
}
