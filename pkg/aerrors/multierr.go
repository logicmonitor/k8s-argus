package aerrors

import "strings"

type MultiError interface {
	Append(error)
	Unwrap() error
	Error() string
}

func GetMultiError(message string, errors ...error) MultiError {
	me := &multiErr{errors: make([]*uerr, 0), message: message}
	for _, err := range errors {
		me.Append(err)
	}
	return me
}

type multiErr struct {
	message string
	errors  []*uerr
}

type uerr struct {
	err error
	*multiErr
}

func (u *uerr) Error() string {
	return u.err.Error()
}

func (u *uerr) Unwrap() error {
	for idx, e := range u.multiErr.errors {
		if e == u {
			if idx < len(u.multiErr.errors)-1 {
				return u.multiErr.errors[idx+1]
			}
			return nil
		}
	}
	return nil
}

func (me *multiErr) Error() string {
	if len(me.errors) > 0 {
		msgs := []string{me.message + ":"}
		for _, e := range me.errors {
			msgs = append(msgs, e.Error())
		}
		return strings.Join(msgs, ", \n")
	}
	return me.message
}

func (me *multiErr) Unwrap() error {
	if len(me.errors) > 0 {
		return me.errors[0]
	}
	return nil
}

func (me *multiErr) Append(err error) {
	me.errors = append(me.errors, &uerr{err: err})
}
