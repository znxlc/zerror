package zerror

import (
	"github.com/stretchr/testify/assert"
	"github.com/znxlc/zerror/errormessage"
	"testing"
)

func TestZError_New(t *testing.T) {
	// testing new error without message
	ze := New()
	assert.Equal(t, 0, len(ze.GetList()))
}

func TestZError_New_WithParams(t *testing.T) {
	// testing new error without message
	ze1 := New(errormessage.ErrorInternal)
	assert.Equal(t, 1, len(ze1.GetList()))
	assert.Equal(t, errormessage.ErrorInternal, ze1.GetList()[0].GetCode())
}

func TestZError_Add(t *testing.T) {
	ze := New()
	ze.Add(errormessage.ErrorInternal)
	assert.Equal(t, 1, len(ze.GetList()))
	assert.Equal(t, errormessage.ErrorInternal, ze.GetList()[0].GetCode())
}

func TestZError_Add_Multiple(t *testing.T) {
	zeLevel1 := New()
	// add registered error by code
	zeLevel1.Add(errormessage.ErrorInternal)

	// add registered error by code with custom elements
	zeLevel1.Add("ERROR_CUSTOM", "new description", map[string]any{"key": "value"})

	// creating new zeLevel2 zerror entity and adding the errors from zeLevel1
	zeLevel2 := New(errormessage.ErrorGeneric)
	zeLevel2.Add(zeLevel1.GetList())

	assert.Equal(t, 3, len(zeLevel2.GetList()))
	assert.Equal(t, errormessage.ErrorGeneric, zeLevel2.GetList()[0].GetCode())
	assert.Equal(t, errormessage.ErrorInternal, zeLevel2.GetList()[1].GetCode())
	assert.Equal(t, "ERROR_CUSTOM", zeLevel2.GetList()[2].GetCode())
}

func TestZError_Get(t *testing.T) {
	zeTest := New("ERROR_1")
	zeTest.Add("ERROR_2")
	zeTest.Add("ERROR_3")

	assert.Equal(t, "ERROR_1", zeTest.Get().GetCode())
	zeTest.SetDefaultElementIndexReturned(FlagReturnLastErrorElement)
	assert.Equal(t, "ERROR_3", zeTest.Get().GetCode())
	assert.Equal(t, "ERROR_2", zeTest.Get(1).GetCode())
}
