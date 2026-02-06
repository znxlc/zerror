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
	// clone errors from a different zerror by passing a []errormessage.IElement returned via ze2.GetList()
	ze4 := zerror.New(ze2.GetList())
	// pass a builtin error to zerror
	ze5 := zerror.New(errors.New("error text"))

	fmt.Printf("ze1.len=%d\n\n", len(ze1.GetList()))

	ze2Element := ze2.Get()
	fmt.Printf("ze2.len=%d\n", len(ze2.GetList()))
	fmt.Printf("ze2.Error=%s\n", ze2.Error())
	fmt.Printf("ze2.code=%s\n", ze2Element.Code())
	fmt.Printf("ze2.msg=%s\n\n", ze2Element.Msg())

	ze3Element := ze3.Get()
	fmt.Printf("ze3.len=%d\n", len(ze3.GetList()))
	fmt.Printf("ze3.Error=%s\n", ze3.Error())
	fmt.Printf("ze3.code=%s\n", ze3Element.Code())
	fmt.Printf("ze3.msg=%s\n", ze3Element.Msg())
	fmt.Printf("ze3.errorList.key=%s\n\n", ze3Element.Args()["key"])

	ze4Element := ze4.Get()
	fmt.Printf("ze4.len=%d\n", len(ze4.GetList()))
	fmt.Printf("ze4.Error=%s\n", ze4.Error())
	fmt.Printf("ze4.code=%s\n", ze4Element.Code())
	fmt.Printf("ze4.msg=%s\n\n", ze4Element.Msg())

	ze5Element := ze5.Get()
	fmt.Printf("ze5.len=%d\n", len(ze5.GetList()))
	fmt.Printf("ze5.Error=%s\n", ze5.Error())
	fmt.Printf("ze5.code=%s\n", ze5Element.Code())
	fmt.Printf("ze5.msg=%s\n", ze5Element.Msg())

	// Output:
	// ze1.len=0
	//
	// ze2.len=1
	// ze2.Error=ERROR_INTERNAL
	// ze2.code=ERROR_INTERNAL
	// ze2.msg=An internal error has occurred
	//
	// ze3.len=1
	// ze3.Error=ERROR_CUSTOM
	// ze3.code=ERROR_CUSTOM
	// ze3.msg=This is a custom error
	// ze3.errorList.key=some value
	//
	// ze4.len=1
	// ze4.Error=ERROR_INTERNAL
	// ze4.code=ERROR_INTERNAL
	// ze4.msg=An internal error has occurred
	//
	// ze5.len=1
	// ze5.Error=ERROR_GENERIC
	// ze5.code=ERROR_GENERIC
	// ze5.msg=error text
}

func ExampleZError_Add() {
	// creating an empty zerror entity
	zeLevel1 := zerror.New()
	// add custom unregistered error
	zeLevel1.Add("ERROR_LEVEL1_FIRST", "This is a level 1 error", map[string]any{"key": "some value"})
	// add registered error with no custom data
	zeLevel1.Add(errormessage.ErrorInternal)
	// add registered error with custom data, overwriting msg and adding errorList
	zeLevel1.Add(errormessage.ErrorGeneric, "err 1.3", map[string]any{"key": "generic key"})
	// add error, code will be defaulted to ErrorGeneric
	zeLevel1.Add(errors.New("err 1.4"))

	// creating a new zerror for level2
	zeLevel2 := zerror.New("ERROR_LEVEL2_FIRST", "This is the first level 2 error", map[string]any{"key": "level2 value"})
	zeLevel2.SetDefaultElementIndexReturned(zerror.FlagReturnLastErrorElement)
	// importing errors from level1
	zeLevel2.Add(zeLevel1.GetList())
	// add another error at the end of the list
	zeLevel2.Add("ERROR_LEVEL2_LAST", "This is the last level 2 error", map[string]any{"key": "level2 last value"})

	// adding a level 3 zerror
	zeLevel3 := zerror.New()
	// set Level3 error
	zeLevel3.Add("ERROR_LEVEL3", "this is a level 3 error")
	// adding main error from Level 2 (as set via zerror ElementIndexReturned)
	zeLevel3.Add(zeLevel2.Get()) // will add the last error entry from Level2
	// adding first error from Level 1
	zeLevel3.Add(zeLevel1.Get()) // will add the first error entry from Level1

	// printing Level1 errors
	fmt.Printf("zeLevel1.len: %d\n", len(zeLevel1.GetList()))
	for idx, errElement := range zeLevel1.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  code: %s\n", errElement.Code())
		fmt.Printf("  msg:  %s\n", errElement.Msg())
		fmt.Printf("  errorList: %v\n", errElement.Args())
	}
	// printing Level2 errors
	fmt.Printf("\nzeLevel2.len: %d\n", len(zeLevel2.GetList()))
	for idx, errElement := range zeLevel2.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  code: %s\n", errElement.Code())
		fmt.Printf("  msg:  %s\n", errElement.Msg())
		fmt.Printf("  errorList: %v\n", errElement.Args())
	}
	fmt.Printf("\nzeLevel3.len: %d\n", len(zeLevel3.GetList()))
	for idx, errElement := range zeLevel3.GetList() {
		fmt.Printf("index: %d\n", idx)
		fmt.Printf("  code: %s\n", errElement.Code())
		fmt.Printf("  msg:  %s\n", errElement.Msg())
		fmt.Printf("  errorList: %v\n", errElement.Args())
	}

	// Output:
	// zeLevel1.len: 4
	// index: 0
	//   code: ERROR_LEVEL1_FIRST
	//   msg:  This is a level 1 error
	//   errorList: map[key:some value]
	// index: 1
	//   code: ERROR_INTERNAL
	//   msg:  An internal error has occurred
	//   errorList: map[]
	// index: 2
	//   code: ERROR_GENERIC
	//   msg:  err 1.3
	//   errorList: map[key:generic key]
	// index: 3
	//   code: ERROR_GENERIC
	//   msg:  err 1.4
	//   errorList: map[]
	//
	// zeLevel2.len: 6
	// index: 0
	//   code: ERROR_LEVEL2_FIRST
	//   msg:  This is the first level 2 error
	//   errorList: map[key:level2 value]
	// index: 1
	//   code: ERROR_LEVEL1_FIRST
	//   msg:  This is a level 1 error
	//   errorList: map[key:some value]
	// index: 2
	//   code: ERROR_INTERNAL
	//   msg:  An internal error has occurred
	//   errorList: map[]
	// index: 3
	//   code: ERROR_GENERIC
	//   msg:  err 1.3
	//   errorList: map[key:generic key]
	// index: 4
	//   code: ERROR_GENERIC
	//   msg:  err 1.4
	//   errorList: map[]
	// index: 5
	//   code: ERROR_LEVEL2_LAST
	//   msg:  This is the last level 2 error
	//   errorList: map[key:level2 last value]
	//
	// zeLevel3.len: 3
	// index: 0
	//   code: ERROR_LEVEL3
	//   msg:  this is a level 3 error
	//   errorList: map[]
	// index: 1
	//   code: ERROR_LEVEL2_LAST
	//   msg:  This is the last level 2 error
	//   errorList: map[key:level2 last value]
	// index: 2
	//   code: ERROR_LEVEL1_FIRST
	//   msg:  This is a level 1 error
	//   errorList: map[key:some value]

}
