package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r = httptest.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestError_Get(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	form.Required("a", "b", "c")
	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	postedData := url.Values{}
	postedData.Add("some_field", "some value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "some value")
	form = New(postedData)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows minlength of 1 is not met when it is")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData := url.Values{}
	postedData.Add("email", "aaaaa@test.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("email", "a")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got a valid email")
	}
}
