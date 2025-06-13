package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io"

	"github.com/google/uuid"
)

var (
	ErrNotImplemented = errors.New("feature not implemented")
	ErrAlreadyClosed  = errors.New("storage is closed")
	ErrInvalidStorage = errors.New("invalid storage")
	ErrInvalidKey     = errors.New("invalid key provided")
	ErrNoSuchEntry    = errors.New("no such entry")
	ErrNoIndex        = errors.New("no index present")
)

type Storage interface {
	Close() error

	Index() (Index, error)
	Rename(id uuid.UUID, newName string) (Index, error)
	MoveUp(id uuid.UUID) (Index, error)
	MoveDown(id uuid.UUID) (Index, error)

	Create(name, initialText string) (Entry, Index, error)
	Read(id uuid.UUID) (Entry, error)
	Update(entry Entry) error
	Delete(id uuid.UUID) (Index, error)
}

func Encode(data any) ([]byte, error) {
	var inputBuffer bytes.Buffer
	if err := gob.NewEncoder(&inputBuffer).Encode(data); err != nil {
		return []byte{}, err
	}
	return inputBuffer.Bytes(), nil
}

func Decode[T any](data []byte, target *T) error {
	var outputBuffer bytes.Buffer
	outputBuffer.Write(data)
	return gob.NewDecoder(&outputBuffer).Decode(target)
}

func NewGCM(key []byte) (cipher.AEAD, error) {
	keySum := sha256.Sum256(key)
	block, err := aes.NewCipher(keySum[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func Harden(gcm cipher.AEAD, input any) ([]byte, error) {
	data, err := Encode(input)
	if err != nil {
		return data, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	cipher := gcm.Seal(nonce, nonce, data, nil)
	return cipher, nil
}

func Soften[T any](gcm cipher.AEAD, encrypted []byte, target *T) error {
	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]
	data, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return err
	}
	return Decode(data, target)
}
