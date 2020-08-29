package encryptor

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"

var (
	genericHashSize int = int(C.crypto_generichash_bytes())
	// GenericKeyMin min key size
	GenericKeyMin int = int(C.crypto_generichash_bytes_min())
	// GenericKeyMax max key size
	GenericKeyMax int = int(C.crypto_generichash_bytes_max())
)

func validGenericKey(key []byte) bool {
	kL := len(key)
	if kL > GenericKeyMax || kL < GenericKeyMin {
		return false
	}
	return true
}

// GenericHash simple algo to generate fingerprint of message, key can not nil
func GenericHash(message, key []byte) ([]byte, error) {
	out := make([]byte, genericHashSize)

	if len(key) > 0 && !validGenericKey(key) {
		return out, ErrKeySize{inRange: []int{GenericKeyMin, GenericKeyMax}}
	}

	_ = C.crypto_generichash(
		(*C.uchar)(BytePointer(out)),
		(C.ulong)(len(out)),
		(*C.uchar)(BytePointer(message)),
		(C.ulonglong)(len(message)),
		(*C.uchar)(BytePointer(key)),
		(C.ulong)(len(key)),
	)

	return out, nil
}
