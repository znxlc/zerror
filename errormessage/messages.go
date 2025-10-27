// The package contains base structs and predefined error messages to be used with zerror package
// additional errormessages can be registered in the main registeredErrorsMap using zerror.RegisterErrors
package errormessage

// Error messages have ErrCode in the format ENTITY_<ATTRIBUTE/VERB>_LIST
const (
  // Main errors
  ErrorGeneric                  = "ERROR_GENERIC"
  ErrorGenerateParameterInvalid = "ERROR_GENERATE_PARAMETER_INVALID"
  ErrorInternal                 = "ERROR_INTERNAL"
  ErrorPanic                    = "ERROR_PANIC"
  ErrorDivByZero                = "ERROR_DIV_BY_ZERO"
  ErrorAuthFailed               = "ERROR_AUTH_FAILED"

  // generic errors
  ErrorJSONParse = "ERROR_JSON_PARSE"
  ErrorRegEx     = "ERROR_REGEX"
)

// RegisteredErrorMap is the main map
var (
  registeredErrorsMap = map[string]Message{
    ErrorGeneric: {
      Code: ErrorGeneric,
      Msg:  "An error has occurred",
    },
    ErrorGenerateParameterInvalid: {
      Code: ErrorGenerateParameterInvalid,
      Msg:  "Unable to generate error element, parameter invalid",
    },
    ErrorInternal: {
      Code: ErrorInternal,
      Msg:  "An internal error has occurred",
    },
    ErrorPanic: {
      Code: ErrorPanic,
      Msg:  "A fatal error has occurred",
    },
    ErrorDivByZero: {
      Code: ErrorDivByZero,
      Msg:  "Division by zero",
    },

    // generic errors
    ErrorJSONParse: {
      Code: ErrorJSONParse,
      Msg:  "Error parsing JSON",
    },
    ErrorRegEx: {
      Code: ErrorRegEx,
      Msg:  "Error in regex",
    },
    ErrorAuthFailed: {
      Code: ErrorAuthFailed,
      Msg:  "Auth Failed",
    },
  }
)
