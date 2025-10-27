package errormessage

import (
  "encoding/json"
)

// Message is the bare error message struct
type Message struct {
  Code string `json:"code"` // error Code
  Msg  string `json:"msg"`  // error message
}

// ErrElement represents a single error element
type ErrElement struct {
  Arguments map[string]any `json:"args,omitempty"` // error optional args
  ErrCode   string         `json:"code"`           // error code
  Message   string         `json:"msg"`            // error message
  // Trace []TraceElement `json:"-"`    // TODO trace list
}

// IElement represents the interface for the ErrElement
type IElement interface {
  Error() string                   // returns element.Msg (helper for error interface)
  Get() IElement                   // returns error element
  Code(code ...string) string      // set/get the error ErrCode
  Msg(msg ...string) string        // set/get the error message
  Args(args ...any) map[string]any // set/get the Arguments
  ClearArgs(keys ...string)        // clears the Arguments based on the specified keys
  Load(string) bool                // adds a clone of the registered error
  Set(args ...any) bool            // sets entire error element
  MarshalJSON() ([]byte, error)    // helper for json.Marshall
  UnmarshalJSON([]byte) error      // helper for json.Unmarshall
}

// ErrorElementGenerator is an alias for the constructor function New
type ErrorElementGenerator = func(args ...any) IElement

// New will generate a new ErrElement starting from ErrorGeneric
func New(args ...any) IElement {
  errElement := new(ErrElement)
  // setting default value
  errElement.Load(ErrorGeneric)
  errElement.Set(args...)

  return errElement
}

// Error returns an error from current element
func (ee *ErrElement) Error() string {
  return ee.Message
}

// Get returns the ErrElement packed in the interface
func (ee *ErrElement) Get() IElement {
  return ee
}

// Set will populate the ErrElement Fields based on a dynamic combination of parameters
//
// @Params
//
//	 the following combinations are supported:
//	   IElement, Message string(optional), Arguments map[string]any(optional), error(optional)
//	     provide a prefilled IElement with the ability to overwrite params
//	   ErrorCode string, Message string(optional), Arguments map[string]any(optional), error(optional)
//	     define a new IElement from scratch
//
//		Arguments[0] [ string | error | errormessage.IElement ]
//		  string
//		    the error ErrCode we wish to use
//		    if found in the registered error list, the entire element will be loaded from there
//		  errormessage.IElement
//		    a prefilled IElement we wish to edit
//		  error
//			the errElement.Message will be set to errorItem.Error()
//		Arguments
//		  will represent the rest of the params needed to create a new IElement (based on type)
//		  string
//		     will set the IElement.Message field to the specified value
//		  map[string]any
//		     will add the keys to IElement.Arguments
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
func (ee *ErrElement) Set(args ...any) bool {
  itemLen := len(args)

  if itemLen > 0 { // we have at least a parameter
    for idx, arg := range args {
      if idx == 0 {
        if arg == nil { // nil element, we ignore it
          return true
        }
        // processing the Arguments[0] as errorItem
        errorItem := arg
        switch eItem := errorItem.(type) {
        case string:
          ee.ErrCode = eItem
          ee.Load(eItem) // load entire IElement if found in registered list, element will remain unchanged if not found
        case Message:
          ee.ErrCode = eItem.Code
          ee.Message = eItem.Msg
        case IElement:
          ee.ErrCode = eItem.Code()
          ee.Message = eItem.Msg()
          ee.Arguments = eItem.Args()
        case map[string]any:
          if code, ok := eItem["code"].(string); ok {
            ee.ErrCode = code
          }
          if msg, ok := eItem["msg"].(string); ok {
            ee.Message = msg
          }
          if ar, ok := eItem["args"].(map[string]any); ok {
            ee.Arguments = ar
          }

        case map[string]string:
          if code, ok := eItem["code"]; ok {
            ee.ErrCode = code
          }
          if msg, ok := eItem["msg"]; ok {
            ee.Message = msg
          }
        case error:
          ee.Message = eItem.Error()
        default: // parameter not supported, the error message will contain the actual error
          ee.Load(ErrorGenerateParameterInvalid)
          ee.Arguments = map[string]any{
            "errorItem":     errorItem,
            "Arguments":     args[1:],
            "expected_type": "string | error | errormessage.IElement",
          }
          return false
        }
        continue
      }

      switch element := arg.(type) {
      case string: // overwriting the Message
        ee.Message = element
      case error:
        ee.Message = element.Error()
      case map[string]any: // add the arguments
        ee.Arguments = element
      }
    }
  }

  return true
}

// Code sets/gets the errorMessage.ErrCode
func (ee *ErrElement) Code(code ...string) string {
  if len(code) > 0 {
    ee.ErrCode = code[0]
  }
  return ee.ErrCode
}

// Msg sets/gets the errorMessage.Message
func (ee *ErrElement) Msg(msg ...string) string {
  if len(msg) > 0 {
    ee.Message = msg[0]
  }
  return ee.Message
}

// Args will set/get the errorMessage.Arguments
// params
//  1 param map[string]any - will add/replace the specified fields in the Arguments
//  2 params key string, value any - adds the key,value pair to the Arguments
func (ee *ErrElement) Args(args ...any) map[string]any {
  key := ""
  var value any
  if len(args) > 0 && ee.Arguments == nil {
    ee.Arguments = make(map[string]any)
  }
  if len(args) == 1 {
    if args[0] == nil {
      return ee.Arguments
    } else if mapArgs, ok := args[0].(map[string]any); ok {
      for key, value = range mapArgs {
        ee.Arguments[key] = value
      }
    }
  } else if len(args) > 1 {
    ok := false
    if args[0] == nil {
      return ee.Arguments
    } else if key, ok = args[0].(string); ok {
      if key != "" {
        value = args[1]
        ee.Arguments[key] = value
      }
    }
  }
  return ee.Arguments
}

// ClearArgs will clear the Arguments in the error message
//  if no keys specified will clear all the Arguments
//  otherwise will only clear the specified keys
func (ee *ErrElement) ClearArgs(keys ...string) {
  if len(keys) > 0 {
    for _, key := range keys {
      if key != "" {
        delete(ee.Arguments, key)
      }
    }
  } else {
    ee.Arguments = make(map[string]any)
  }

}

// Load will attempt to create a copy of a registered error and populate the object with its fields
func (ee *ErrElement) Load(code string) bool {
  errElement, found := registeredErrorsMap[code]
  if found {
    ee.ErrCode = errElement.Code
    ee.Message = errElement.Msg
  }
  return found
}

// UnmarshalJSON is a function to make IElement compatible with json.Marshal.
func (ee *ErrElement) UnmarshalJSON(data []byte) error {
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
func (ee *ErrElement) MarshalJSON() ([]byte, error) {
  return json.Marshal(*ee)
}
