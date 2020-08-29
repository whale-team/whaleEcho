package encryptor

import "strconv"

// ErrKeySize customeized error represent key size error
type ErrKeySize struct {
	correctSize int
	inRange     []int
}

func (err ErrKeySize) rangeStr() string {
	if len(err.inRange) == 0 {
		return ""
	}

	return strconv.Itoa(err.inRange[0]) + " to " + strconv.Itoa(err.inRange[1])
}

func (err ErrKeySize) Error() string {
	if len(err.inRange) > 0 {
		return "given key size is invalid, should be in range of " + err.rangeStr()
	}

	return "given key size is invalid, should be " + strconv.Itoa(err.correctSize)
}

// ErrVerifiactionFailed customeized error represent verification error
type ErrVerifiactionFailed struct {
	alg string
}

func (err ErrVerifiactionFailed) Error() string {
	return err.alg + " verification failed"
}
