package encryptor

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import (
	"encoding/base64"
	"unsafe"
)

var (
	// KeySize for sodium secretbox key
	KeySize int = int(C.crypto_secretbox_keybytes())
	// NonceSize for sodium secretbox nonce key
	NonceSize int = int(C.crypto_secretbox_noncebytes())
)

// GenRandomBtyes generate random butes
func GenRandomBtyes(size int) []byte {
	buf := make([]byte, size)
	C.randombytes_buf(unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return buf
}

// GenRandomKey generate secretbox key
func GenRandomKey() []byte {
	buf := make([]byte, KeySize)
	C.randombytes_buf(unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return buf
}

// GenNonce generate secretbox nonce
func GenNonce(size int) []byte {
	buf := make([]byte, NonceSize)
	C.randombytes_buf(unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return buf
}

// GenRandomString generate random string
func GenRandomString(size int) string {
	buf := make([]byte, size)
	C.randombytes_buf(unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return base64.StdEncoding.EncodeToString(buf)
}
