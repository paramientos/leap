package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"filippo.io/age"
)

func Encrypt(data []byte, passphrase string) ([]byte, error) {
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	w, err := age.Encrypt(buf, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decrypt(data []byte, passphrase string) ([]byte, error) {
	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, err
	}

	r, err := age.Decrypt(bytes.NewReader(data), identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// Obfuscate encrypts data using a key derived from system strings
func Obfuscate(data []byte, salt string) ([]byte, error) {
	key := sha256.Sum256([]byte(salt))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Deobfuscate decrypts data using a key derived from system strings
func Deobfuscate(data []byte, salt string) ([]byte, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("invalid data")
	}

	key := sha256.Sum256([]byte(salt))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
