package errormessage

// RegisterErrors will add predefined error codes to the main errormessage.RegisteredErrorMap so they can be accessed by error code.
// This is a wrapper function that can be used in a generic app that loads predefined error messages from various places(packages, config files, etc...).
//
// @Params
//
//	args ...any
//	  single argument - assumes the parameter contains a fully defined ErrorElement or a list of elements
//	    combination of parameters that allows for registering error messages
//	    []ErrorElement   - adds the list directly
//	    ErrorElement     - adds the element to the list
//	    []map[string]any - imports the list from a generic map if the structure is compatible
//	    map[string]any   - imports the element from map if compatible
//	    string | []byte  - assumes this is a json or yaml list
//	  multiple arguments - assumes you are registering a single ErrorElement and it must be compatible with Add() parameters
func RegisterErrors(args ...any) {
	itemLen := len(args)
	if itemLen == 1 { // fully defined errormessage.ErrorElement or a list of elements
		switch element := args[0].(type) {
		case []errorElement:
			registerErrorElementList(element...)
		case errorElement:
			registerErrorElementList(element)
		}
	} else if itemLen > 1 {

	}
}

// registerErrorElementList adds ErrorElement items to the RegisteredErrorsMap
func registerErrorElementList(element ...errorElement) {
	if len(element) > 0 {
		for _, errElement := range element {
			RegisteredErrorsMap[errElement.Code] = errElement
		}
	}
}

// GetRegisteredElement returns the element defined by key or will return the ErrorGeneric otherwise
func GetRegisteredElement(key string) ErrorElement {
	elem := New(ErrorGeneric)
	if !elem.Load(key) {
		elem.Load(ErrorGeneric)
	}

	return elem
}

// GenerateErrorElement will create a new ErrorElement from various parameter combinations
//
// @Params
//
//	 the following combinations are supported:
//	   ErrorElement, Msg string(optional), Args map[string]any(optional), error(optional)
//	     provide a prefilled ErrorElement with the ability to overwrite params
//	   ErrorCode string, Msg string(optional), Args map[string]any(optional), error(optional)
//	     define a new ErrorElement from scratch
//
//		errorCode [ string | error | errormessage.ErrorElement ]
//		  string
//		    the error code we wish to use
//		    if found in the registered error list, the entire element will be loaded from there
//		  errormessage.ErrorElement
//		    a prefilled ErrorElement we wish to edit
//		  error
//			the errElement.Msg will be set to errorItem.Error()
//		args
//		  will represent the rest of the params needed to create a new ErrorElement (based on type)
//		  string
//		     will set the ErrorElement.Msg field to the specified value
//		  map[string]any
//		     will add the keys to ErrorElement.Args
//		  TraceElement, []TraceElement
//		     will append the TraceElement to ErrorElement.Trace
//		  other
//		     will be ignored
//
// @Returns
//
//	errElement [ errormessage.ErrorElement ]
//	   the error element
//	result [ bool ]
//	   true
//	      the element was generated successfully
//	   false
//	      the element could not be generated from the provided parameters
//	      errElement will contain a more detailed error
func GenerateErrorElement(args ...any) (errElement ErrorElement, result bool) {
	errElement = New(args)
	return errElement, true
}
