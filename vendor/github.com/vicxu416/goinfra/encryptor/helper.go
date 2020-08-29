package encryptor

// BytePointer returns a pointer to the start of a byte slice, or nil when the slice is empty.
func BytePointer(b []byte) *uint8 {
	if len(b) > 0 {
		return &b[0]
	}
	return nil
}


