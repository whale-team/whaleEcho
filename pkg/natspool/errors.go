package natspool

import "errors"

// ErrPoolClose ...
var ErrPoolClose = errors.New("pool has been already closed")

// ErrConnOverLimit ...
var ErrConnOverLimit = errors.New("number of connection is over the limit")

// ErrGetConnTimeout ...
var ErrGetConnTimeout = errors.New("get connection timeout")
