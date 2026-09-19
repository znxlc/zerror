package zerror

import "github.com/znxlc/zerror/errormessage"

// this file contains types and variables used by the package

// Error flags and predefined codes
const (
	FlagReturnLastErrorElement  = "LAST"  // returns the last error in the list by default (Get() and Error() functions)
	FlagReturnFirstErrorElement = "FIRST" // returns the first element in the list when calling Get() and Error()
	FlagReturnErrorCode         = "CODE"  // default text returned when calling Error()
	FlagReturnErrorMsg          = "MSG"   // return msg field when calling Error()
)

var (
	// ElementIndexReturned represents the default element to be returned when calling single IElement functions like Get() and Error()
	ElementIndexReturned = FlagReturnFirstErrorElement
	// ElementTextReturned is used by Error() to select which text to return (default is Error code)
	ElementTextReturned = FlagReturnErrorCode
	// DefaultElementGenerator will be used to create new error elements and should be a pointer to the constructor of the errorElement used
	DefaultElementGenerator = errormessage.New
)

// ZError is the main error structure of the package
type ZError struct {
	ElementIndexReturned string                             `json:"-"`                // set the default element to be returned when calling Get() or Error()
	ElementGenerator     errormessage.ErrorElementGenerator `json:"-"`                // the generator for the error elements (pointer to the New() constructor)
	Errors               []errormessage.IElement            `json:"errors,omitempty"` // the error list
}

type Error interface {
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
	Add(...any)
	
	// Clear will reset the Errors list to an empty list
	Clear()
	
	// Error will return a specific element (based on ElementIndexReturned and ElementTextReturned) wrapped as an error string
	Error() string
	
	// GetList returns the list of errors
	GetList() []errormessage.IElement
	
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
	Get(...int) errormessage.IElement
	
	// Has will return true if the Errors list contains the code specified
	Has(errCode string) bool
	
	// HasErrors will return true if the Errors list contains elements
	HasErrors() bool
	
	// SetDefaultElementIndexReturned will set the default element returned when using Get() or Error()
	SetDefaultElementIndexReturned(string)
}
