package helper

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicateEntry  = errors.New("duplicate entry")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMissingId       = errors.New("id is required")
	ErrMissingToken    = errors.New("token is required")
	ErrParsingValue    = errors.New("err parsing param value")
)

func InvalidRequest(err string) error {
	ErrInvalidArgument = errors.New(err)
	return ErrInvalidArgument
}

func str2err(err string) error {
	if err == "" {
		return nil
	}
	return errors.New(err)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func ErrEqual(err1 error, err2 error) bool {
	if err1 == nil {
		err1 = errors.New("")
	}
	if err2 == nil {
		err1 = errors.New("")
	}
	return err1.Error() == err2.Error()
}
