package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrRecordExists = errors.New("record exists")
var ErrWrongContentType = errors.New("wrong content type in request")
