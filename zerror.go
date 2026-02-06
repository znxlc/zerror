// Package zerror explores an alternate error structure for whoever needs to use multiple error responses
// without having to create complex structures and wrappers.
//
// Functions use variable number of parameters so anyone can reuse the interface with different structures.
//
// The New function will create a new zerror object and can accept a variable number of parameters (up to 4).
//
// Multiple errors can be added to the list using the Add function in the desired order, creating something similar to a trace list.
//
// Errors can be retrieved via Error(), Get() and GetList().
package zerror

import (
	"encoding/json"

	errormessage "github.com/znxlc/zerror/errormessage"
)

// New creates a new zerror instance, see Add() for the parameter format.
func New(args ...any) Error {
	ze := &ZError{}
	ze.Clear() // generate a clear error list
	ze.ElementIndexReturned = ElementIndexReturned
	ze.ElementGenerator = DefaultElementGenerator
	if len(args) > 0 {
		ze.Add(args...)
	}

	return ze
}

// Add will append an error element to the Errors list based on the type of the parameters.
//
// Using variable number of parameters for more versatility.
//
// @Params
//
//	  errorList[0] [string | map[string]any | []any | error | IElement | []IElement]
//	    depending on type, this parameter will be interpreted as follows:
//	  string - Error code
//	  Error - will add the errors in the object
//	  error  - will set the Error code to generic and will set msg to error.Error()
//	  IElement - will append the IElement to the list, rest of the params will overwrite the initial element
//	  []IElement - will append the IElement to the list, rest of the params will be ignored
//	  map[string]any - will convert it to IElement and add it to the list
//	  []any          - will trigger Add(element...)
//	  []map[string]any - will convert it to []IElement and add it to the list
//
//	  errorList[1-3] [string | map[string]any | error]
//		optional parameter list based on type
//	    string - IElement.msg
//		map[string]any - optional IElement.errorList
//		error - will set the IElement.msg to error.Error()
func (ze *ZError) Add(args ...any) {
	itemLen := len(args)

	if itemLen > 0 { // we have at least a parameter
		errorItem := args[0]
		if errorItem == nil { // we skip adding an element if the param is nil
			return
		}

		switch element := errorItem.(type) {
		case []errormessage.IElement:
			ze.Errors = append(ze.Errors, element...)
			return
		case errormessage.IElement:
			ze.Errors = append(ze.Errors, element)
		case Error:
			ze.Add(element.GetList())
		case map[string]any, map[string]string:
			em := errormessage.New(element)
			ze.Add(em)
		case []any:
			ze.Add(element...)
		case []map[string]any:
			for _, errMap := range element {
				em := errormessage.New(errMap)
				ze.Add(em)
			}
		case []map[string]string:
			for _, errMap := range element {
				em := errormessage.New(errMap)
				ze.Add(em)
			}
		case error: // needs to be the last before default since error interface will match others ( like Error)
			ze.Add(errormessage.ErrorGeneric, element.Error())
		default: // generate a new error element
			errElement := ze.ElementGenerator(args...)
			ze.Errors = append(ze.Errors, errElement)
			return
		}
	}
	// pushing rest of the errorList if they are IElement
	if itemLen > 1 {
		for _, errorItem := range args[1:] {
			if element, ok := errorItem.(errormessage.IElement); ok {
				ze.Errors = append(ze.Errors, element)
			}

		}
	}
}

// Clear will reset the Errors list to an empty list
func (ze *ZError) Clear() {
	ze.Errors = []errormessage.IElement{}
}

// Error will return a specific element (based on ElementIndexReturned and ElementTextReturned) wrapped as an error string
func (ze *ZError) Error() string {
	errElement := ze.Get()
	if errElement == nil {
		return ""
	}
	if ElementTextReturned == FlagReturnErrorMsg {
		return errElement.Error()
	}
	return errElement.Code()
}

// Get returns a pointer to the IElement specified
//
// @Params
//
//	no param
//	   gets first or last element as specified in zerror.ElementIndexReturned
//	index [ int ]
//	   gets the element specified by index from the Errors list
//
// @Returns
//
//	nil
//	   no errors exist or index out of bounds
//	*errormessage.IElement
//	   errors exist and Errors[index] was found or no index is specified
func (ze *ZError) Get(index ...int) errormessage.IElement {
	errLen := len(ze.Errors)
	if errLen > 0 {
		if len(index) == 0 {
			if ze.ElementIndexReturned == FlagReturnFirstErrorElement {
				return ze.Errors[0]
			}
			return ze.Errors[errLen-1]
		}
		idx := index[0]
		if idx >= len(ze.Errors) {
			return nil
		}
		return ze.Errors[idx]
	}
	return nil
}

// GetList returns the list of errors
func (ze *ZError) GetList() []errormessage.IElement {
	return ze.Errors
}

// Has will return true if the Errors list contains the code specified
func (ze *ZError) Has(errCode string) bool {
	for _, errElement := range ze.Errors {
		if errElement.Code() == errCode {
			return true
		}
	}
	return false
}

// HasErrors will return true if the Errors list contains elements
func (ze *ZError) HasErrors() bool {
	return len(ze.Errors) > 0
}

// SetDefaultElementIndexReturned will set the default element returned when using Get() or Error()
func (ze *ZError) SetDefaultElementIndexReturned(flag string) {
	switch flag {
	case FlagReturnFirstErrorElement, FlagReturnLastErrorElement:
		ze.ElementIndexReturned = flag
	}
}

// UnmarshalJSON is a function to make IElement compatible with json.Marshal.
func (ze *ZError) UnmarshalJSON(data []byte) (err error) {
	errList := make([]any, 0) // generic potential list
	err = json.Unmarshal(data, &errList)
	if err != nil {
		return err
	}

	for _, errorElement := range errList {
		ze.Add(errorElement)
	}

	return nil
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
func (ze *ZError) MarshalJSON() (result []byte, err error) {
	if ze.HasErrors() {
		return json.Marshal(ze.Errors)
	}
	return json.Marshal(nil)
}
