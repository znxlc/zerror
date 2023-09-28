// The package contains base structs and predefined error messages to be used with zerror package
// additional errormessages can be registered in the main RegisteredErrorsMap using zerror.RegisterErrors
package errormessage

// Error messages have code in the format ENTITY_<ATTRIBUTE/VERB>_LIST
const (
  ErrorGeneric                  = "ERROR_GENERIC"
  ErrorGenerateParameterInvalid = "ERROR_GENERATE_PARAMETER_INVALID"
  ErrorInternal                 = "ERROR_INTERNAL"
  ErrorPanic                    = "ERROR_PANIC"
)

// RegisteredErrorMap is the main map
var (
  RegisteredErrorsMap = map[string]TElement{
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
  }
)
