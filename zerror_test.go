package zerror

import (
  "encoding/json"
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
  assert.Equal(t, errormessage.ErrorInternal, ze1.GetList()[0].Code())
}

func TestZError_Add(t *testing.T) {
  ze := New()
  ze.Add(errormessage.ErrorInternal)
  assert.Equal(t, 1, len(ze.GetList()))
  assert.Equal(t, errormessage.ErrorInternal, ze.GetList()[0].Code())
}

func TestZError_Add_Map(t *testing.T) {
  ze := New()
  // add registered error by code
  ze.Add(map[string]any{
    "code": "ERROR_1",
    "msg":  "error 1",
    "args": map[string]any{"k": "v"},
  })
  ze.Add([]map[string]any{
    {
      "code": "ERROR_2",
      "msg":  "error 2",
    },
    {
      "code": "ERROR_3",
      "msg":  "error 3",
    },
  })
  assert.Equal(t, 3, len(ze.GetList()))
  assert.Equal(t, "ERROR_1", ze.Get().Code()) // default to first error which is index=0
  assert.Equal(t, "ERROR_2", ze.Get(1).Code())
  assert.Equal(t, "ERROR_3", ze.Get(2).Code())
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
  assert.Equal(t, errormessage.ErrorGeneric, zeLevel2.Get().Code()) // default to first error which is index=0
  assert.Equal(t, errormessage.ErrorInternal, zeLevel2.Get(1).Code())
  assert.Equal(t, "ERROR_CUSTOM", zeLevel2.Get(2).Code())
}

func TestZError_Get(t *testing.T) {
  zeTest := New("ERROR_1")
  zeTest.Add("ERROR_2")
  zeTest.Add("ERROR_3")

  assert.Equal(t, "ERROR_1", zeTest.Get().Code())
  zeTest.SetDefaultElementIndexReturned(FlagReturnLastErrorElement)
  assert.Equal(t, "ERROR_3", zeTest.Get().Code())
  assert.Equal(t, "ERROR_2", zeTest.Get(1).Code())
}

func TestUnit_JSONMarshall(t *testing.T) {
  zeTest := New(errormessage.ErrorInternal, "e1")
  zeTest.Add(errormessage.ErrorGeneric, "e2")
  zeTest.Add(errormessage.ErrorDivByZero, "e3")

  js, er := json.Marshal(zeTest)
  if er != nil {
    t.Error(er)
    return
  }

  expected := `[{"code":"ERROR_INTERNAL","msg":"e1"},{"code":"ERROR_GENERIC","msg":"e2"},{"code":"ERROR_DIV_BY_ZERO","msg":"e3"}]`
  assert.Equal(t, expected, string(js))
}
