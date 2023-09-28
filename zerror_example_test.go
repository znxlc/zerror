package zerror_test

import (
	"errors"
	"fmt"
	"github.com/znxlc/zerror"
	"github.com/znxlc/zerror/errormessage"
)

func ExampleNew() {
	// creating an empty error, elements can be added later via Add()
	ze1 := zerror.New()
	// registered error by code
	ze2 := zerror.New(errormessage.ErrorInternal)
	// custom unregistered error
	ze3 := zerror.New("ERROR_CUSTOM", "This is a custom error", map[string]any{"key": "some value"})
	// clone errors from a different zerror by passing a []errormessage.ErrorElement returned via ze2.GetList()
	ze4 := zerror.New(ze2.GetList())
	// pass a builtin error to zerror
	ze5 := zerror.New(errors.New("error text"))

	fmt.Printf("ze1.len=%d\n\n", len(ze1.GetList()))

	ze2Element := ze2.Get()
	fmt.Printf("ze2.len=%d\n", len(ze2.GetList()))
	fmt.Printf("ze2.Error=%s\n", ze2.Error())
	fmt.Printf("ze2.Code=%s\n", ze2Element.GetCode())
	fmt.Printf("ze2.Msg=%s\n\n", ze2Element.GetMsg())

	ze3Element := ze3.Get()
	fmt.Printf("ze3.len=%d\n", len(ze3.GetList()))
	fmt.Printf("ze3.Error=%s\n", ze3.Error())
	fmt.Printf("ze3.Code=%s\n", ze3Element.GetCode())
	fmt.Printf("ze3.Msg=%s\n", ze3Element.GetMsg())
	fmt.Printf("ze3.Args.key=%s\n\n", ze3Element.GetArgs()["key"])

	ze4Element := ze4.Get()
	fmt.Printf("ze4.len=%d\n", len(ze4.GetList()))
	fmt.Printf("ze4.Error=%s\n", ze4.Error())
	fmt.Printf("ze4.Code=%s\n", ze4Element.GetCode())
	fmt.Printf("ze4.Msg=%s\n\n", ze4Element.GetMsg())

	ze5Element := ze5.Get()
	fmt.Printf("ze5.len=%d\n", len(ze5.GetList()))
	fmt.Printf("ze5.Error=%s\n", ze5.Error())
	fmt.Printf("ze5.Code=%s\n", ze5Element.GetCode())
	fmt.Printf("ze5.Msg=%s\n", ze5Element.GetMsg())

	// Output:
	// ze1.len=0
	//
	// ze2.len=1
	// ze2.Error=ERROR_INTERNAL
	// ze2.Code=ERROR_INTERNAL
	// ze2.Msg=An internal error has occurred
	//
	// ze3.len=1
	// ze3.Error=ERROR_CUSTOM
	// ze3.Code=ERROR_CUSTOM
	// ze3.Msg=This is a custom error
	// ze3.Args.key=some value
	//
	// ze4.len=1
	// ze4.Error=ERROR_INTERNAL
	// ze4.Code=ERROR_INTERNAL
	// ze4.Msg=An internal error has occurred
	//
	// ze5.len=1
	// ze5.Error=ERROR_GENERIC
	// ze5.Code=ERROR_GENERIC
	// ze5.Msg=error text
}

func ExampleZError_Add() {
	// creating an empty zerror entity
	zeLevel1 := zerror.New()
	// add custom unregistered error
	zeLevel1.Add("ERROR_LEVEL1", "This is a level 1 error", map[string]any{"key": "some value"})
	// add registered error with no custom data
	zeLevel1.Add(errormessage.ErrorInternal)
	// add registered error with custom data, overwriting Msg and adding Args
	zeLevel1.Add(errormessage.ErrorGeneric, "different generic error message", map[string]any{"key": "generic key"})
	// add error, code will be defaulted to ErrorGeneric
	zeLevel1.Add(errors.New("some error"))

	// creating a new zerror for level2
	zeLevel2 := zerror.New("ERROR_LEVEL2_FIRST", "This is the first level 2 error", map[string]any{"key": "level2 value"})
	// importing errors from level1
	zeLevel2.Add(zeLevel1.GetList())
	// add another error at the end of the list
	zeLevel2.Add("ERROR_LEVEL2_LAST", "This is the last level 2 error", map[string]any{"key": "level2 last value"})

	// adding a level 3 zerror
	zeLevel3 := zerror.New()
	// set Level3 error
	zeLevel3.Add("ERROR_LEVEL3", "this is a level 3 error")
	// adding main error from Level 2 (as set via zerror ElementIndexReturned)
	zeLevel2.SetDefaultElementIndexReturned(zerror.FlagReturnLastErrorElement)
	zeLevel3.Add(zeLevel2.Get()) // will add the last error entry from Level2
	// adding first error from Level 1
	zeLevel3.Add(zeLevel1.Get()) // will add the first error entry from Level1

	// printing Level1 errors
	fmt.Printf("zeLevel1.len: %d\n", len(zeLevel1.GetList()))
	for idx, errElement := range zeLevel1.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  Code: %s\n", errElement.GetCode())
		fmt.Printf("  Msg:  %s\n", errElement.GetMsg())
		fmt.Printf("  Args: %v\n", errElement.GetArgs())
	}
	// printing Level2 errors
	fmt.Printf("\nzeLevel2.len: %d\n", len(zeLevel2.GetList()))
	for idx, errElement := range zeLevel2.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  Code: %s\n", errElement.GetCode())
		fmt.Printf("  Msg:  %s\n", errElement.GetMsg())
		fmt.Printf("  Args: %v\n", errElement.GetArgs())
	}
	fmt.Printf("\nzeLevel3.len: %d\n", len(zeLevel3.GetList()))
	for idx, errElement := range zeLevel3.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  Code: %s\n", errElement.GetCode())
		fmt.Printf("  Msg:  %s\n", errElement.GetMsg())
		fmt.Printf("  Args: %v\n", errElement.GetArgs())
	}

	// Output:
	// zeLevel1.len: 4
	// index: 0
	//   Code: ERROR_LEVEL1
	//   Msg:  This is a level 1 error
	//   Args: map[key:some value]
	// index: 1
	//   Code: ERROR_INTERNAL
	//   Msg:  An internal error has occurred
	//   Args: map[]
	// index: 2
	//   Code: ERROR_GENERIC
	//   Msg:  different generic error message
	//   Args: map[key:generic key]
	// index: 3
	//   Code: ERROR_GENERIC
	//   Msg:  some error
	//   Args: map[]
	//
	// zeLevel2.len: 6
	// index: 0
	//   Code: ERROR_LEVEL2_FIRST
	//   Msg:  This is the first level 2 error
	//   Args: map[key:level2 value]
	// index: 1
	//   Code: ERROR_LEVEL1
	//   Msg:  This is a level 1 error
	//   Args: map[key:some value]
	// index: 2
	//   Code: ERROR_INTERNAL
	//   Msg:  An internal error has occurred
	//   Args: map[]
	// index: 3
	//   Code: ERROR_GENERIC
	//   Msg:  different generic error message
	//   Args: map[key:generic key]
	// index: 4
	//   Code: ERROR_GENERIC
	//   Msg:  some error
	//   Args: map[]
	// index: 5
	//   Code: ERROR_LEVEL2_LAST
	//   Msg:  This is the last level 2 error
	//   Args: map[key:level2 last value]
	//
	// zeLevel3.len: 3
	// index: 0
	//   Code: ERROR_LEVEL3
	//   Msg:  this is a level 3 error
	//   Args: map[]
	// index: 1
	//   Code: ERROR_LEVEL2_LAST
	//   Msg:  This is the last level 2 error
	//   Args: map[key:level2 last value]
	// index: 2
	//   Code: ERROR_LEVEL1
	//   Msg:  This is a level 1 error
	//   Args: map[key:some value]

}
