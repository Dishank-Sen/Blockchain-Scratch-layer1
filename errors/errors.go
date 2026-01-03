package errors

import "errors"

var ErrDaemonTimeout = errors.New("daemon did not become ready in time")
