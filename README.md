# zerror

A powerful and flexible error handling library for Go that provides structured error handling with error codes, contextual data, and error chaining capabilities.

## Features

- **Structured Errors**: Define errors with codes, messages, and custom arguments
- **Error Chaining**: Build complex error hierarchies with multiple error messages
- **Type Safety**: Strongly typed error codes and arguments
- **JSON Support**: Built-in JSON marshaling/unmarshaling
- **Flexible API**: Multiple ways to create and handle errors
- **Standard Library Integration**: Implements the standard `error` interface

## Installation

```bash
go get github.com/znxlc/zerror
```

## Basic Usage

### Creating a simple error

```go
package main

import (
    "fmt"
    "github.com/znxlc/zerror"
)

func main() {
    // Create a new error with just a code, if error is registered, the existing description will be used else just the code will exist
    err := zerror.New("INVALID_INPUT")
    
    // Add a second error with more details
    err.Add("AUTH_ERROR", "Invalid credentials", map[string]any{
        "username": "testuser",
        "attempts": 3,
    })
    
    // Print the last error code added
    fmt.Println(err.Error()) // Outputs: AUTH_ERROR
}
```

### Error Validation Example

```go
package main

import (
    "fmt"
    "github.com/znxlc/zerror"
)

type SignupRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func validateSignup(req SignupRequest) (ze zerror.Error) {
    ze = zerror.New()
    
    // Validate username
    if req.Username == "" {
        ze.Add("USERNAME_REQUIRED", "Username is required")
    } else if len(req.Username) < 3 {
        ze.Add("USERNAME_TOO_SHORT", "Username too short", 
            map[string]any{"min_length": 3, "actual": len(req.Username)})
    }
    
    // Validate email
    if req.Email == "" {
        ze.Add("EMAIL_REQUIRED", "Email is required")
    } else if !isValidEmail(req.Email) { // use a custom validation function
        ze.Add("INVALID_EMAIL", "Invalid email format", 
            map[string]any{"email": req.Email})
    }
    
    // Validate password
    if len(req.Password) < 8 {
        ze.Add("PASSWORD_TOO_SHORT", "Password must be at least 8 characters",
            map[string]any{"min_length": 8, "actual": len(req.Password)})
    }
    
    return ze
}

func main() {
    req := SignupRequest{
        Username: "ab",
        Email:    "invalid-email",
        Password: "123",
    }
    
    if err := validateSignup(req); err.HasErrors() {
        // Print all validation errors
        for _, e := range err.GetList() {
            fmt.Printf("Error: %s - %s : %v\n", e.Code(), e.Msg(), e.Args())
        }
      // test a certain error exists
      if err.Has("USERNAME_TOO_SHORT") {
        fmt.Println("Username is too short")
      }
    }
}
```

## Advanced Features

### Error Chaining

```go
func processRequest() (ze zerror.Error) {
    if err := validateRequest(); err.HasErrors() {
        ze.Add("PROCESS_ERROR", "Failed to process request")
        ze.Add(err) // Chain the validation errors
        return
    }
    // ... processing logic ...
    return nil
}
```

### JSON Serialization

```go
func handleAPIError(w http.ResponseWriter, err zerror.Error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": false,
        "errors":  err.GetList(),
    })
}
```

### Custom Error Elements

```go
type CustomError struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Details map[string]any `json:"details,omitempty"`
}

// Implement errormessage.IElement interface
func (e *CustomError) Code() string            { return e.Code }
func (e *CustomError) Error() string           { return e.Message }
func (e *CustomError) Args() map[string]any    { return e.Details }
func (e *CustomError) SetArgs(args map[string]any) { e.Details = args }

// Usage
customErr := &CustomError{
    Code:    "CUSTOM_ERROR",
    Message: "Something went wrong",
    Details: map[string]any{"retry_after": 30},
}

err := zerror.New()
err.Add(customErr)
```

## API Reference

### Creating Errors

- `zerror.New() Error` - Create a new empty error
- `zerror.New(code string) Error` - Create a new error with a code
- `zerror.New(code, message string) Error` - Create a new error with code and message
- `zerror.New(code, message string, args map[string]any) Error` - Create a new error with code, message, and arguments

### Error Methods

- `Add(...interface{})` - Add error(s) to the error list
- `Clear()` - Remove all errors
- `Error() string` - Returns the error message
- `Get() IElement` - Get the primary error (either first or last as defined by err.ElementIndexReturned)
- `Get(index int) IElement` - Get the error at the specified index
- `GetList() []IElement` - Get all errors
- `Has(code string) bool` - Check if error with code exists
- `HasErrors() bool` - Check if there are any errors
- `SetDefaultElementIndexReturned(flag string)` - Set which error to return (FIRST or LAST)

## Best Practices

1. **Use Constants for Error Codes**: Define error codes as constants for better maintainability
2. **Include Context**: Always include relevant context in error arguments
3. **Handle Errors Gracefully**: Check for errors using `HasErrors()` before proceeding
4. **Use Meaningful Error Codes**: Make error codes descriptive and consistent
5. **Document Error Types**: Document what each error code means and when it occurs

## License

MIT