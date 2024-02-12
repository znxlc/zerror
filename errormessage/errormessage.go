package errormessage

import (
	"encoding/json"
)

// Message is the bare error message struct
type Message struct {
	Code string `json:"code"` // error code
	Msg  string `json:"msg"`  // error message
}

// tElement represents a single error element
type tElement struct {
	Args map[string]any `json:"args"` // error optional args
	Code string         `json:"code"` // error code
	Msg  string         `json:"msg"`  // error message
	//Trace []TraceElement `json:"-"`    // TODO trace list
}

// IElement represents the interface for the tElement
type IElement interface {
	Error() string
	Get() IElement
	GetCode() string
	GetMsg() string
	GetArgs() map[string]any
	Load(string) bool
	Set(args ...any) bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// ErrorElementGenerator is an alias for the constructor function New
type ErrorElementGenerator = func(args ...any) IElement

// New will generate a new tElement starting from ErrorGeneric
func New(args ...any) IElement {
	errElement := new(tElement)
	// setting default value
	errElement.Load(ErrorGeneric)
	errElement.Set(args...)

	return errElement
}

// Error returns an error from current element
func (ee *tElement) Error() string {
	return ee.Msg
}

// Get returns the tElement packed in the interface
func (ee *tElement) Get() IElement {
	return ee
}

// Set will populate the tElement Fields based on a dynamic combination of parameters
//
// @Params
//
//	 the following combinations are supported:
//	   IElement, Msg string(optional), Args map[string]any(optional), error(optional)
//	     provide a prefilled IElement with the ability to overwrite params
//	   ErrorCode string, Msg string(optional), Args map[string]any(optional), error(optional)
//	     define a new IElement from scratch
//
//		args[0] [ string | error | errormessage.IElement ]
//		  string
//		    the error code we wish to use
//		    if found in the registered error list, the entire element will be loaded from there
//		  errormessage.IElement
//		    a prefilled IElement we wish to edit
//		  error
//			the errElement.Msg will be set to errorItem.Error()
//		args
//		  will represent the rest of the params needed to create a new IElement (based on type)
//		  string
//		     will set the IElement.Msg field to the specified value
//		  map[string]any
//		     will add the keys to IElement.Args
//		  TraceElement, []TraceElement
//		     will append the TraceElement to IElement.Trace
//		  other
//		     will be ignored
//
// @Returns
//
//	errElement [ errormessage.IElement ]
//	   the error element
//	result [ bool ]
//	   true
//	      the element was generated successfully
//	   false
//	      the element could not be generated from the provided parameters
//	      errElement will contain a more detailed error
func (ee *tElement) Set(args ...any) bool {
	itemLen := len(args)

	if itemLen > 0 { // we have at least a parameter
		for idx, arg := range args {
			if idx == 0 {
				// processing the args[0] as errorItem
				errorItem := arg
				switch eItem := errorItem.(type) {
				case string:
					ee.Code = eItem
					ee.Load(eItem) // load entire IElement if found in registered list, element will remain unchanged if not found
				case Message:
					ee.Code = eItem.Code
					ee.Msg = eItem.Msg
				case IElement:
					ee.Code = eItem.GetCode()
					ee.Msg = eItem.GetMsg()
					ee.Args = eItem.GetArgs()
				case error:
					ee.Msg = eItem.Error()
				default: // parameter not supported, the error message will contain the actual error
					ee.Load(ErrorGenerateParameterInvalid)
					ee.Args = map[string]any{
						"errorItem":     errorItem,
						"args":          args[1:],
						"expected_type": "string | error | errormessage.IElement",
					}
					return false
				}
				continue
			}

			switch element := arg.(type) {
			case string: // overwriting the Msg
				ee.Msg = element
			case error:
				ee.Msg = element.Error()
			case map[string]any: // add the arguments
				ee.Args = element
			}
		}
	}

	return true
}

// GetCode returns the errorMessage.Code
func (ee *tElement) GetCode() string {
	return ee.Code
}

// GetMsg returns the errorMessage.Msg
func (ee *tElement) GetMsg() string {
	return ee.Msg
}

// GetArgs returns the errorMessage.Args
func (ee *tElement) GetArgs() map[string]any {
	return ee.Args
}

// Load will attempt to create a copy of a registered error and populate the object with its fields
func (ee *tElement) Load(code string) bool {
	errElement, found := registeredErrorsMap[code]
	if found {
		ee.Code = errElement.Code
		ee.Msg = errElement.Msg
	}
	return found
}

// UnmarshalJSON is a function to make IElement compatible with json.Marshal.
func (ee *tElement) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, ee)
}

// MarshalJSON is a function to make IElement compatible with json.Marshal.
//
// Inputs:
//
//	(none)
//
// Outputs:
//
//	[]byte
//	  The JSON representation of the IElement struct
//	error
//	  Marshal error, if any occurred
func (ee *tElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(ee)
}
