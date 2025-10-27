package errormessage

import (
  "encoding/json"
  "gopkg.in/yaml.v3"
)

// RegisterErrors will add predefined error codes to the main errormessage.RegisteredErrorMap so they can be accessed by error ErrCode.
// This is a wrapper function that can be used in a generic app that loads predefined error messages from various places(packages, config files, etc...).
//
// @Params
//
//	Arguments ...any
//	  single argument - assumes the parameter contains a fully defined IElement or a list of elements
//	    combination of parameters that allows for registering error messages
//	    []IElement   - adds the list directly
//	    IElement     - adds the element to the list
//	    Message     - adds the element to the list
//	    []map[string]any - imports the list from a generic map if the structure is compatible
//	    map[string]any   - imports the element from map if compatible
//	    string | []byte  - assumes this is a json or yaml map or slice
//	  multiple arguments - assumes you are registering elements that must be compatible with Message structure
//      ErrCode, Message string - register a single message with the properties specfied via these Arguments
//      Message... - array of messages
func RegisterErrors(args ...any) {
  itemLen := len(args)
  if itemLen == 1 { // fully defined message or a list of elements
    switch element := args[0].(type) {
    case []IElement:
      registerErrorElementList(element...)
    case IElement:
      registerErrorElementList(element)
    case Message:
      registeredErrorsMap[element.Code] = element
    case []Message:
      for _, msg := range element {
        registeredErrorsMap[msg.Code] = msg
      }
    case map[string]Message: // most common case
      for key, value := range element {
        registeredErrorsMap[key] = value
      }
    case string: // we have an error ErrCode or a json/yaml
      if !registerProcessStringList(element) { // element was not json/yaml, we assume it is a ErrCode
        newMessage := Message{
          Code: element,
        }
        if itemLen > 1 {
          if msg, msgOk := args[1].(string); msgOk {
            newMessage.Msg = msg
          }
        }
        registeredErrorsMap[element] = newMessage
      }
    }
  }
  if itemLen > 1 { // we have a potential list of Message, ignoring other types
    for _, element := range args {
      switch message := element.(type) {
      case Message:
        registeredErrorsMap[message.Code] = message
      }
    }
  }
}

// registerErrorElementList adds IElement items to the registeredErrorsMap
func registerErrorElementList(args ...IElement) {
  if len(args) > 0 {
    for _, element := range args {
      registeredErrorsMap[element.Code()] = Message{element.Code(), element.Msg()}
    }
  }
}

func registerProcessStringList[T string | []byte](config T) bool {
  listData := ([]byte)(config)
  resultMap := map[string]Message{}
  resultSlice := make([]Message, 0)

  if err := json.Unmarshal(listData, &resultMap); err == nil {
    for key, value := range resultMap {
      registeredErrorsMap[key] = value
    }
    return true
  }
  // try to marshal to a slice
  if err := json.Unmarshal(listData, &resultSlice); err != nil {
    for _, elem := range resultSlice {
      registeredErrorsMap[elem.Code] = elem
    }
    return true
  }
  // try to see if it is yaml
  if err := yaml.Unmarshal(listData, &resultMap); err == nil {
    for key, value := range resultMap {
      registeredErrorsMap[key] = value
    }
    return true
  }
  // try to marshal to a slice
  if err := yaml.Unmarshal(listData, &resultSlice); err != nil {
    for _, elem := range resultSlice {
      registeredErrorsMap[elem.Code] = elem
    }
    return true
  }

  return false
}

// GetRegisteredElement returns the element defined by key or will return the ErrorGeneric otherwise
func GetRegisteredElement(key string) IElement {
  elem := New(ErrorGeneric)
  if !elem.Load(key) {
    elem.Load(ErrorGeneric)
  }

  return elem
}

// GenerateErrorElement will create a new IElement from various parameter combinations
//
// @Params
//
//	 the following combinations are supported:
//	   IElement, Message string(optional), Arguments map[string]any(optional), error(optional)
//	     provide a prefilled IElement with the ability to overwrite params
//	   ErrorCode string, Message string(optional), Arguments map[string]any(optional), error(optional)
//	     define a new IElement from scratch
//
//		errorCode [ string | error | errormessage.IElement ]
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
func GenerateErrorElement(args ...any) (errElement IElement, result bool) {
  errElement = New(args)
  return errElement, true
}
