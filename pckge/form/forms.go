package form

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) HasRequired(tagIDS ...string) {
	for _, tagID := range tagIDS {
		value := f.Get(tagID)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(tagID, "This field can't be blank")
		}
	}
}

func (f *Form) HasValue(tagID string, r *http.Request) bool {
	x := r.Form.Get(tagID)
	// if x == "" {
	// 	f.Errors.AddError(tagID, "Field Empty")
	// 	return false
	// }
	return x != ""
}

func (f *Form) MinLength(tagID string, length int, r *http.Request) bool {
	x := r.Form.Get(tagID)
	if len(x) < length {
		f.Errors.AddError(tagID, fmt.Sprintf("This field must be %d charcaters long or more", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(tagID string) {
	if !govalidator.IsEmail(f.Get(tagID)) {
		f.Errors.AddError(tagID, "Invalid Email")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
