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
//	  args[0] [string | map[string]any | error | ErrorElement | []ErrorElement]
//		    depending on type, this parameter will be interpreted as follows:
//		    string - Error Code
//		    error  - will set the Error Code to generic and will set Msg to error.Error()
//		    ErrorElement - will append the ErrorElement to the list, rest of the params will overwrite the initial element
//		    []ErrorElement - will append the ErrorElement to the list, rest of the params will be ignored
//
//	  args[1-3] [string | map[string]any | error]
//			optional parameter list based on type
//			string - ErrorElement.Msg
//			map[string]any - optional ErrorElement.Args
//			error - will set the ErrorElement.Msg to error.Error()
func (ze *ZError) Add(args ...any) {
	itemLen := len(args)

	if itemLen > 0 { // we have at least a parameter
		errorItem := args[0]
		switch element := errorItem.(type) {
		case []errormessage.ErrorElement:
			ze.Errors = append(ze.Errors, element...)
			return
		default: // generate a new error element
			errElement := ze.ElementGenerator(args...)
			ze.Errors = append(ze.Errors, errElement)
			return
		}
	}
	// pushing rest of the args if they are ErrorElement
	if itemLen > 1 {
		for _, errorItem := range args[1:] {
			if element, ok := errorItem.(errormessage.ErrorElement); ok {
				ze.Errors = append(ze.Errors, element)
			}

		}
	}
}

// Clear will reset the Errors list to an empty list
func (ze *ZError) Clear() {
	ze.Errors = []errormessage.ErrorElement{}
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
	return errElement.GetCode()
}

// Get returns a pointer to the ErrorElement specified
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
//	*errormessage.ErrorElement
//	   errors exist and Errors[index] was found or no index is specified
func (ze *ZError) Get(index ...int) errormessage.ErrorElement {
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
func (ze *ZError) GetList() []errormessage.ErrorElement {
	return ze.Errors
}

// Has will return true if the Errors list contains the code specified
func (ze *ZError) Has(errCode string) bool {
	for _, errElement := range ze.Errors {
		if errElement.GetCode() == errCode {
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
