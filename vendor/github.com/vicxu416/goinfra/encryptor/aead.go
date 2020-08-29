package encryptor

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import (
	"crypto/cipher"
	"errors"
)

var (
	// AEADKeySize aead key size
	AEADKeySize int = int(C.crypto_aead_aes256gcm_keybytes())
	// PubNonceSize public nonce size
	PubNonceSize int = int(C.crypto_aead_aes256gcm_npubbytes())
	secretSize   int = int(C.crypto_aead_aes256gcm_abytes())
)

var (
	// ErrAEADEncryptionError error represent aead encryption failed
	ErrAEADEncryptionError = errors.New("aead encryption failed")
)

// NewAES256GCM construct AES256GCM that implement AEAD interface
func NewAES256GCM(key []byte) (cipher.AEAD, error) {
	if len(key) != AEADKeySize {
		return nil, ErrKeySize{correctSize: AEADKeySize}
	}

	return aeadCryptor{
		key:       key,
		nonceSize: PubNonceSize,
		overhead:  secretSize,
	}, nil
}

type aeadCryptor struct {
	key       []byte
	nonceSize int
	overhead  int
}

func (cryptor aeadCryptor) NonceSize() int {
	return cryptor.nonceSize
}

func (cryptor aeadCryptor) Overhead() int {
	return cryptor.overhead
}

func (cryptor aeadCryptor) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return aeadAES256Encrypt(plaintext, additionalData, nonce, cryptor.key)
}

func (cryptor aeadCryptor) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return aeadAES256Decrypt(ciphertext, additionalData, nonce, cryptor.key)
}

func aeadAES256Encrypt(plain, additional, publicNonce, key []byte) []byte {
	out := make([]byte, len(plain)+secretSize)

	l := C.ulonglong(len(out))

	_ = C.crypto_aead_aes256gcm_encrypt(
		(*C.uchar)(BytePointer(out)),
		(*C.ulonglong)(&l),
		(*C.uchar)(BytePointer(plain)),
		(C.ulonglong)(len(plain)),
		(*C.uchar)(BytePointer(additional)),
		(C.ulonglong)(len(additional)),
		(*C.uchar)(nil),
		(*C.uchar)(&publicNonce[0]),
		(*C.uchar)(&key[0]))

	return out
}

func aeadAES256Decrypt(cipher, additional, publicNonce, key []byte) ([]byte, error) {
	out := make([]byte, len(cipher)-secretSize)
	l := (C.ulonglong)(len(out))

	exit := int(C.crypto_aead_aes256gcm_decrypt(
		(*C.uchar)(BytePointer(out)),
		(*C.ulonglong)(&l),
		(*C.uchar)(nil),
		(*C.uchar)(&cipher[0]),
		(C.ulonglong)(len(cipher)),
		(*C.uchar)(BytePointer(additional)),
		(C.ulonglong)(len(additional)),
		(*C.uchar)(&publicNonce[0]),
		(*C.uchar)(&key[0])))

	if exit < 0 {
		return out, ErrVerifiactionFailed{alg: "AEAD"}
	}
	return out, nil
}
