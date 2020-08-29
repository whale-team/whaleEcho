package encryptor

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import (
	"errors"
	"unsafe"
)

var (
	// ErrPasswordFailedVerify customize error represent password verification failed
	ErrPasswordFailedVerify = errors.New("password verification failed")
	// ErrEncryptedPasswordInvalid customize error represent invalid encrypted password len
	ErrEncryptedPasswordInvalid = errors.New("encrypted password has incorrect length")
)

var (
	// PWDSaltSize for password hashing
	PWDSaltSize int = int(C.crypto_pwhash_saltbytes())
	pwdHashStr  int = int(C.crypto_pwhash_strbytes())
	pwdhashMin  int = int(C.crypto_pwhash_bytes_min())
	pwdPolMin   int = int(C.crypto_pwhash_opslimit_min())
	pwdmemMin   int = int(C.crypto_pwhash_memlimit_min())
	pwdalt      int = int(C.crypto_pwhash_alg_default())
)

// PasswordHashing generate password hasn
func PasswordHashing(password string, salt []byte) []byte {
	out := make([]byte, pwdhashMin)
	C.crypto_pwhash(
		(*C.uchar)(&out[0]),
		(C.ulonglong)(len(out)),
		C.CString(password),
		(C.ulonglong)(len(password)),
		(*C.uchar)(&salt[0]),
		(C.ulonglong)(pwdPolMin),
		(C.ulong)(pwdmemMin),
		(C.int)(pwdalt),
	)

	return out
}

// EncryptPassword encrypt password that generate salt automatically
func EncryptPassword(password string) []byte {
	out := make([]byte, pwdHashStr)

	C.crypto_pwhash_str(
		(*C.char)(unsafe.Pointer(&out[0])),
		C.CString(password),
		(C.ulonglong)(len(password)),
		(C.ulonglong)(pwdPolMin),
		(C.ulong)(pwdmemMin),
	)
	return out
}

// VerifyPassword verify password
func VerifyPassword(encrypted []byte, password string) error {
	if len(encrypted) != pwdHashStr {
		return ErrEncryptedPasswordInvalid
	}

	res := int(C.crypto_pwhash_str_verify(
		(*C.char)(unsafe.Pointer(&encrypted[0])),
		C.CString(password),
		(C.ulonglong)(len(password)),
	))
	if res < 0 {
		return ErrVerifiactionFailed{alg: "password encryption"}
	}
	return nil
}
