package roomscenter

import "errors"

// ErrOpenRoom represent open room error
var ErrOpenRoom = errors.New("open room failed")

// ErrCloseRoom represet close room error
var ErrCloseRoom = errors.New("close room failed")
