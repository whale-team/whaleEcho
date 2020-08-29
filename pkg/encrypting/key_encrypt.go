package encrypting

import (
	"encoding/base64"

	"github.com/vicxu416/goinfra/encryptor"
)

func Encrypt(plain, additional []byte, nonce, key string) (string, error) {
	nonceDecoded, keyDecoded, err := decodeNonceAndKey(nonce, key)
	if err != nil {
		return "", err
	}
	aesEcnryptor, err := encryptor.NewAES256GCM(keyDecoded)
	if err != nil {
		return "", err
	}

	cipher := aesEcnryptor.Seal(plain[:0], nonceDecoded, plain, additional)
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func Decrypt(cipher string, additional []byte, nonce, key string) ([]byte, error) {
	nonceDecoded, keyDecoded, err := decodeNonceAndKey(nonce, key)
	if err != nil {
		return []byte{}, err
	}
	aesEcnryptor, err := encryptor.NewAES256GCM(keyDecoded)
	if err != nil {
		return []byte{}, err
	}

	cipherDecoded, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return []byte{}, err
	}

	return aesEcnryptor.Open(cipherDecoded[:0], nonceDecoded, cipherDecoded, additional)
}

func decodeNonceAndKey(nonce, key string) ([]byte, []byte, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	nonceDecoded, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return nonceDecoded, keyDecoded, nil
}
