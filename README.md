# zerror

This library explores an alternate approach to error handling when the projects require it.

### Major concepts:

- Provide more detail to errors by using Codes and dynamic Args alongside the Msg
- Error lists for a more complex context
- use interfaces to enable easy extension creation
- 

## install
```
  go  get github.com/znxlc/zerror
```

## usage

Below is an example which shows some common use cases

```go
package main

import (
  "fmt"
  "github.com/znxlc/zerror"
)

type TestRequest struct {
    User string `json:"user"`
    Email string `json:"email"`
}
  
  
// simple example that validates and logs a request  
func main(){
  req := TestRequest{
    User: "test",
    Email: "example@test.com",
  }
  
  err := Process(req)
  
  if err != nil {
    fmt.Printf("Error processing request: %s", err.Error())
  }
}

// function that returns a standard error from a zerror
func Process(req TestRequest) error {
  err := LogRequest(req)
  
  return err
}

// simple function that will return a zerror in case validation fails
func LogRequest(req TestRequest) (ze zerror.Error){

  err := validateRequest(req)
  // validate will always return a non nil zerror so we check if it has error messages
  if err.HasErrors() { 
    ze.Add("ERROR_REQUEST_VALIDATION", "Request validation error")
    ze.Add(err.GetList()) // we add the errors from validateRequest to our return
    fmt.Printf("Request: user=%s, email=%s, error: %#v\n", req.User, req.Email)
    
    return
  }
  
  fmt.Printf("Request: user=%s, email=%s\n", req.User, req.Email)
  
  return nil // we can return a nil error
}

// Validate request will return a zerror element with errors if the request contains invalid fields 
func validateRequest (req TestRequest) (ze zerror.Error){
    ze = zerror.New() // returning a non nil zerror even if validation succeeds
    
    if req.User == ""{
      ze.Add(
        "ERROR_USER_INVALID", // the code 
        "User name is invalid", // optional message
        map[string]any{ // optional arguments to make error more descriptive
          "user": req.User,
          "expected": "non empty user",          
        })
    } else if len(req.User) < 8 { // expecting user lentgth to be at least 8 
      ze.Add(
        "ERROR_USER_LENGTH", // the code 
        "User length is less than minimum required", // optional message
        map[string]any{ // optional arguments to make error more descriptive
          "user": req.User,
          "user_length": len(req.User),
          "expected_length": 8,
        })
    }
    
    return
}

```