package zerror

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/znxlc/zerror/errormessage"
)

func PtrInt(i int) *int {
	return &i
}

func TestUnit_New(t *testing.T) {
	tests := []struct {
		name              string
		errorList         []any
		expectedLen       int
		expectedErrorCode string
		expectedMessage   string
		checkArgs         bool
		expectedArgs      map[string]any
	}{
		{
			name:      "empty",
			errorList: []any{},
		},
		{
			name:              "with internal error",
			errorList:         []any{errormessage.ErrorInternal},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorInternal,
			expectedMessage:   "An internal error has occurred",
		},
		{
			name:              "with internal error and errorList",
			errorList:         []any{errormessage.ErrorInternal, map[string]any{"user_id": 123, "action": "delete"}},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorInternal,
			expectedMessage:   "An internal error has occurred",
			checkArgs:         true,
			expectedArgs:      map[string]any{"user_id": 123, "action": "delete"},
		},
		{
			name:              "with internal error, custom msg and errorList",
			errorList:         []any{errormessage.ErrorInternal, "Custom internal error message", map[string]any{"module": "auth", "severity": "high"}},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorInternal,
			expectedMessage:   "Custom internal error message",
			checkArgs:         true,
			expectedArgs:      map[string]any{"module": "auth", "severity": "high"},
		},
		{
			name:              "with generic error",
			errorList:         []any{errormessage.ErrorGeneric},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorGeneric,
			expectedMessage:   "An error has occurred",
		},
		{
			name:              "with generic error and errorList",
			errorList:         []any{errormessage.ErrorGeneric, map[string]any{"endpoint": "/api/users", "method": "POST"}},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorGeneric,
			expectedMessage:   "An error has occurred",
			checkArgs:         true,
			expectedArgs:      map[string]any{"endpoint": "/api/users", "method": "POST"},
		},
		{
			name:              "with generic error, custom msg and errorList",
			errorList:         []any{errormessage.ErrorGeneric, "API request failed", map[string]any{"status_code": 500, "response_time": "2.5s"}},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorGeneric,
			expectedMessage:   "API request failed",
			checkArgs:         true,
			expectedArgs:      map[string]any{"status_code": 500, "response_time": "2.5s"},
		},
		{
			name:              "with div by zero error",
			errorList:         []any{errormessage.ErrorDivByZero},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorDivByZero,
			expectedMessage:   "Division by zero",
		},
		{
			name:              "with div by zero and errorList",
			errorList:         []any{errormessage.ErrorDivByZero, map[string]any{"numerator": 10, "denominator": 0}},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorDivByZero,
			expectedMessage:   "Division by zero",
			checkArgs:         true,
			expectedArgs:      map[string]any{"numerator": 10, "denominator": 0},
		},
		{
			name:              "with auth failed error",
			errorList:         []any{errormessage.ErrorAuthFailed},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorAuthFailed,
			expectedMessage:   "Auth Failed",
		},
		{
			name:              "with auth failed and custom msg",
			errorList:         []any{errormessage.ErrorAuthFailed, "Authentication token expired"},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorAuthFailed,
			expectedMessage:   "Authentication token expired",
		},
		{
			name:              "with custom code and message",
			errorList:         []any{"CUSTOM_ERROR", "custom message"},
			expectedLen:       1,
			expectedErrorCode: "CUSTOM_ERROR",
			expectedMessage:   "custom message",
		},
		{
			name:              "with code, message and errorList",
			errorList:         []any{"ERROR_WITH_ARGS", "message with errorList", map[string]any{"key": "value", "number": 42}},
			expectedLen:       1,
			expectedErrorCode: "ERROR_WITH_ARGS",
			expectedMessage:   "message with errorList",
			checkArgs:         true,
			expectedArgs:      map[string]any{"key": "value", "number": 42},
		},
		{
			name:              "with error type",
			errorList:         []any{assert.AnError},
			expectedLen:       1,
			expectedErrorCode: errormessage.ErrorGeneric,
			expectedMessage:   assert.AnError.Error(),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New(testCase.errorList...)
			assert.Equal(t, testCase.expectedLen, len(zerror.GetList()))
			if testCase.expectedLen > 0 {
				errorElement := zerror.GetList()[0]
				assert.Equal(t, testCase.expectedErrorCode, errorElement.Code())
				assert.Equal(t, testCase.expectedMessage, errorElement.Msg())

				if testCase.checkArgs && testCase.expectedArgs != nil {
					args := errorElement.Args()
					assert.NotNil(t, args)
					for key, expectedVal := range testCase.expectedArgs {
						assert.Equal(t, expectedVal, args[key], "Args key '%s' should match expected value", key)
					}
				}
			}
		})
	}
}

func TestUnit_Add(t *testing.T) {
	tests := []struct {
		name               string
		errorList          []any
		expectedLen        int
		expectedErrorCodes []string
		expectedMsg        map[int]string         // index -> expected message
		expectedArgs       map[int]map[string]any // index -> expected args
	}{
		// Single error additions
		{
			name:               "single error",
			errorList:          []any{[]any{errormessage.ErrorInternal, "Internal error occurred"}},
			expectedLen:        1,
			expectedErrorCodes: []string{errormessage.ErrorInternal},
			expectedMsg:        map[int]string{0: "Internal error occurred"},
			expectedArgs:       nil,
		},
		{
			name:               "map error",
			errorList:          []any{map[string]any{"code": "ERROR_1", "msg": "error 1", "args": map[string]any{"k": "v"}}},
			expectedLen:        1,
			expectedErrorCodes: []string{"ERROR_1"},
			expectedMsg:        map[int]string{0: "error 1"},
			expectedArgs:       map[int]map[string]any{0: map[string]any{"k": "v"}},
		},
		{
			name:      "nil parameter",
			errorList: []any{nil},
		},
		{
			name:               "error type",
			errorList:          []any{assert.AnError},
			expectedLen:        1,
			expectedErrorCodes: []string{errormessage.ErrorGeneric},
			expectedMsg:        map[int]string{0: assert.AnError.Error()},
		},

		// Array additions
		{
			name: "string map array",
			errorList: []any{
				map[string]string{"code": "ERROR_1", "msg": "error 1"},
				map[string]string{"code": "ERROR_2", "msg": "error 2"},
			},
			expectedLen:        2,
			expectedErrorCodes: []string{"ERROR_1", "ERROR_2"},
			expectedMsg:        map[int]string{0: "error 1", 1: "error 2"},
		},
		{
			name: "element array",
			errorList: []any{
				errormessage.New(errormessage.ErrorInternal),
				errormessage.New(errormessage.ErrorGeneric),
			},
			expectedLen:        2,
			expectedErrorCodes: []string{errormessage.ErrorInternal, errormessage.ErrorGeneric},
			expectedMsg:        map[int]string{0: "An internal error has occurred", 1: "An error has occurred"},
		},

		// Mixed error types
		{
			name: "mixed types",
			errorList: []any{
				errormessage.ErrorInternal,
				map[string]any{"code": "CUSTOM_1", "msg": "custom message 1", "args": map[string]any{"type": "custom", "id": 123}},
				"ERROR_STRING",
				assert.AnError,
				errormessage.New(errormessage.ErrorDivByZero),
				[]any{"SLICE_ERROR", "slice error message", map[string]any{"slice": true, "nested": map[string]any{"inner": "value"}}},
			},
			expectedLen:        6,
			expectedErrorCodes: []string{errormessage.ErrorInternal, "CUSTOM_1", "ERROR_STRING", errormessage.ErrorGeneric, errormessage.ErrorDivByZero, "SLICE_ERROR"},
			expectedMsg:        map[int]string{0: "An internal error has occurred", 1: "custom message 1", 3: assert.AnError.Error(), 5: "slice error message"},
			expectedArgs:       map[int]map[string]any{1: map[string]any{"type": "custom", "id": 123}, 5: map[string]any{"slice": true, "nested": map[string]any{"inner": "value"}}},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			assert.Equal(t, testCase.expectedLen, len(zerror.GetList()))

			// Validate expected error codes
			if len(testCase.expectedErrorCodes) > 0 {
				for _, expectedErrorCode := range testCase.expectedErrorCodes {
					assert.True(t, zerror.Has(expectedErrorCode))
				}
			}

			// Validate expected messages
			if testCase.expectedMsg != nil {
				for index, expectedMsg := range testCase.expectedMsg {
					actualElement := zerror.Get(index)
					assert.NotNil(t, actualElement)
					assert.Equal(t, expectedMsg, actualElement.Msg())
				}
			}

			// Validate expected args
			if testCase.expectedArgs != nil {
				for index, expectedArgs := range testCase.expectedArgs {
					actualElement := zerror.Get(index)
					assert.NotNil(t, actualElement)
					actualArgs := actualElement.Args()
					assert.NotNil(t, actualArgs)
					for key, expectedValue := range expectedArgs {
						assert.Equal(t, expectedValue, actualArgs[key], "Args key '%s' at index %d should match expected value", key, index)
					}
				}
			}
		})
	}
}

func TestUnit_Basic_Functions(t *testing.T) {
	zerror := New(errormessage.ErrorInternal, "test1")
	zerror.Add(errormessage.ErrorGeneric, "test2")
	zerror.Add("CUSTOM_ERROR", "test3")

	// Test len (GetList length)
	actualLength := len(zerror.GetList())
	assert.Equal(t, 3, actualLength)

	// Test Get(FIRST)
	zerror.SetDefaultElementIndexReturned(FlagReturnFirstErrorElement)
	firstErrorElement := zerror.Get()
	assert.Equal(t, errormessage.ErrorInternal, firstErrorElement.Code())
	assert.Equal(t, "test1", firstErrorElement.Msg())

	// Test Get(LAST)
	zerror.SetDefaultElementIndexReturned(FlagReturnLastErrorElement)
	lastErrorElement := zerror.Get()
	assert.Equal(t, "CUSTOM_ERROR", lastErrorElement.Code())
	assert.Equal(t, "test3", lastErrorElement.Msg())

	// Test GetList()
	errorList := zerror.GetList()
	assert.Equal(t, 3, len(errorList))
	assert.Equal(t, errormessage.ErrorInternal, errorList[0].Code())
	assert.Equal(t, errormessage.ErrorGeneric, errorList[1].Code())
	assert.Equal(t, "CUSTOM_ERROR", errorList[2].Code())
}

func TestUnit_Get(t *testing.T) {
	tests := []struct {
		name              string
		errorList         []any
		getIndex          *int
		defaultFlag       string
		expectedErrorCode string
		shouldNil         bool
	}{
		{
			name:              "default first",
			errorList:         []any{"ERROR_1"},
			defaultFlag:       FlagReturnFirstErrorElement,
			expectedErrorCode: "ERROR_1",
			shouldNil:         false,
		},
		{
			name:              "default last",
			errorList:         []any{"ERROR_1", "ERROR_2"},
			defaultFlag:       FlagReturnLastErrorElement,
			expectedErrorCode: "ERROR_2",
			shouldNil:         false,
		},
		{
			name:              "specific index",
			errorList:         []any{"ERROR_1", "ERROR_2", "ERROR_3"},
			getIndex:          PtrInt(1),
			expectedErrorCode: "ERROR_2",
		},
		{
			name:      "out of bounds",
			errorList: []any{"ERROR_1"},
			getIndex:  PtrInt(999),
			shouldNil: true,
		},
		{
			name:      "empty list",
			errorList: []any{},
			shouldNil: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			// Set default flag if specified
			if testCase.defaultFlag != "" {
				zerror.SetDefaultElementIndexReturned(testCase.defaultFlag)
			}

			result := zerror.Get() // fetch default error
			// Get the error element
			if testCase.getIndex != nil { // fetch error from specified index
				result = zerror.Get(*testCase.getIndex)
			}

			if testCase.shouldNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, testCase.expectedErrorCode, result.Code())
			}
		})
	}
}

func TestUnit_JSON(t *testing.T) {
	tests := []struct {
		name              string
		errorList         []any
		expectedJSON      string
		shouldMarshalNull bool
	}{
		{
			name: "marshal with errors",
			errorList: []any{
				[]any{errormessage.ErrorInternal, "e1"},
				errormessage.ErrorGeneric,
			},
			expectedJSON:      `[{"code":"ERROR_INTERNAL","msg":"e1"},{"code":"ERROR_GENERIC","msg":"An error has occurred"}]`,
			shouldMarshalNull: false,
		},
		{
			name:         "marshal empty",
			errorList:    []any{},
			expectedJSON: "null",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			js, err := json.Marshal(zerror)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, string(js))
		})
	}
}

func TestUnit_Clear(t *testing.T) {
	tests := []struct {
		name           string
		errorList      []any
		expectedBefore int
		expectedAfter  int
	}{
		{
			name:           "clear single error",
			errorList:      []any{errormessage.ErrorInternal},
			expectedBefore: 1,
			expectedAfter:  0,
		},
		{
			name:           "clear multiple errors",
			errorList:      []any{"ERROR_1", "ERROR_2", "ERROR_3"},
			expectedBefore: 3,
			expectedAfter:  0,
		},
		{
			name:      "clear empty list",
			errorList: []any{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			// Verify initial state
			assert.Equal(t, testCase.expectedBefore, len(zerror.GetList()))

			// Clear and verify
			zerror.Clear()
			assert.Equal(t, testCase.expectedAfter, len(zerror.GetList()))
		})
	}
}

func TestUnit_Error(t *testing.T) {
	tests := []struct {
		name           string
		errorList      []any
		textFlag       string
		expectedResult string
	}{
		{
			name:           "return code",
			errorList:      []any{errormessage.ErrorInternal},
			textFlag:       FlagReturnErrorCode,
			expectedResult: errormessage.ErrorInternal,
		},
		{
			name:           "return message",
			errorList:      []any{errormessage.ErrorInternal},
			textFlag:       FlagReturnErrorMsg,
			expectedResult: "An internal error has occurred",
		},
		{
			name:      "empty list",
			errorList: []any{},
			textFlag:  FlagReturnErrorCode,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			originalFlag := ElementTextReturned
			ElementTextReturned = testCase.textFlag
			defer func() { ElementTextReturned = originalFlag }()

			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			result := zerror.Error()
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestUnit_Has(t *testing.T) {
	tests := []struct {
		name            string
		errorList       []any
		searchErrorCode string
		shouldExist     bool
	}{
		{
			name:            "existing error",
			errorList:       []any{errormessage.ErrorInternal, "test message", "CUSTOM_ERROR", "custom message"},
			searchErrorCode: errormessage.ErrorInternal,
			shouldExist:     true,
		},
		{
			name:            "custom error",
			errorList:       []any{errormessage.ErrorInternal, "test message", "CUSTOM_ERROR", "custom message"},
			searchErrorCode: "CUSTOM_ERROR",
			shouldExist:     true,
		},
		{
			name:            "nonexistent error",
			errorList:       []any{errormessage.ErrorInternal, "test message", "CUSTOM_ERROR", "custom message"},
			searchErrorCode: "NONEXISTENT_ERROR",
			shouldExist:     false,
		},
		{
			name:            "empty list",
			errorList:       []any{},
			searchErrorCode: errormessage.ErrorInternal,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			result := zerror.Has(testCase.searchErrorCode)
			assert.Equal(t, testCase.shouldExist, result)
		})
	}
}

func TestUnit_HasErrors(t *testing.T) {
	tests := []struct {
		name           string
		errorList      []any
		shouldClear    bool
		expectedResult bool
	}{
		{
			name:      "empty list",
			errorList: []any{},
		},
		{
			name:           "with errors",
			errorList:      []any{errormessage.ErrorInternal},
			shouldClear:    false,
			expectedResult: true,
		},
		{
			name:           "after clear",
			errorList:      []any{errormessage.ErrorInternal},
			shouldClear:    true,
			expectedResult: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			// Clear if specified
			if testCase.shouldClear {
				zerror.Clear()
			}

			result := zerror.HasErrors()
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestUnit_SetDefaultElementIndexReturned(t *testing.T) {
	tests := []struct {
		name              string
		errorList         []any
		flag              string
		expectedErrorCode string
	}{
		{
			name:              "invalid flag",
			errorList:         []any{errormessage.ErrorInternal, errormessage.ErrorGeneric},
			flag:              "INVALID_FLAG",
			expectedErrorCode: errormessage.ErrorInternal,
		},
		{
			name:              "first error",
			errorList:         []any{errormessage.ErrorInternal, errormessage.ErrorGeneric},
			flag:              FlagReturnFirstErrorElement,
			expectedErrorCode: errormessage.ErrorInternal,
		},
		{
			name:              "last error",
			errorList:         []any{errormessage.ErrorInternal, errormessage.ErrorGeneric},
			flag:              FlagReturnLastErrorElement,
			expectedErrorCode: errormessage.ErrorGeneric,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()

			// Add each error from the errorList
			for _, errorItem := range testCase.errorList {
				zerror.Add(errorItem)
			}

			zerror.SetDefaultElementIndexReturned(testCase.flag)
			result := zerror.Get().Code()
			assert.Equal(t, testCase.expectedErrorCode, result)
		})
	}
}

func TestUnit_JSON_Unmarshal(t *testing.T) {
	tests := []struct {
		name              string
		jsonData          string
		shouldError       bool
		expectedLen       int
		expectedErrorCode string
		expectedMessage   string
	}{
		{
			name:              "valid JSON",
			jsonData:          `[{"code":"ERROR_INTERNAL","msg":"test message"},{"code":"ERROR_GENERIC","msg":"generic error"}]`,
			shouldError:       false,
			expectedLen:       2,
			expectedErrorCode: "ERROR_INTERNAL",
			expectedMessage:   "test message",
		},
		{
			name:              "valid JSON array",
			jsonData:          `[["ERROR_INTERNAL","test message"],["ERROR_GENERIC","generic error"]]`,
			shouldError:       false,
			expectedLen:       2,
			expectedErrorCode: "ERROR_INTERNAL",
			expectedMessage:   "test message",
		},
		{
			name:        "invalid JSON",
			jsonData:    `{"invalid": json}`,
			shouldError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			zerror := New()
			err := json.Unmarshal([]byte(testCase.jsonData), zerror)

			if testCase.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedLen, len(zerror.GetList()))
				if testCase.expectedLen > 0 {
					assert.Equal(t, testCase.expectedErrorCode, zerror.Get().Code())
					assert.Equal(t, testCase.expectedMessage, zerror.Get().Msg())
				}
			}
		})
	}
}
