package storage

import (
	"crypto/cipher"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type BoltStorage struct {
	bolt   *bolt.DB
	cipher cipher.AEAD
}

var (
	indexKey  = []byte("index")
	bucketKey = []byte("hardnote")
)

func NewBoltStorage(filename string, keyText []byte) (Storage, error) {

	gcm, err := NewGCM(keyText)

	if err != nil {
		return BoltStorage{}, err
	}

	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return BoltStorage{}, err
	}

	store := BoltStorage{
		bolt:   db,
		cipher: gcm,
	}

	_, err = store.Index()
	if err != nil {
		return BoltStorage{}, err
	}

	return store, nil
}

func (b BoltStorage) Close() error {
	return b.bolt.Close()
}

func (b BoltStorage) bucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketKey)
	if bucket == nil {
		return tx.CreateBucket(bucketKey)
	}
	return bucket, nil
}

func (b BoltStorage) Index() (Index, error) {
	idx := Index{}
	err := b.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := b.bucket(tx)
		if err != nil {
			return err
		}
		value := bucket.Get(indexKey)
		if value == nil {
			data, err := Harden(b.cipher, idx)
			if err != nil {
				return err
			}
			return bucket.Put(indexKey, data)
		}
		return Soften(b.cipher, value, &idx)
	})
	return idx, err
}
func (b BoltStorage) updateIndex(idx Index) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := b.bucket(tx)
		if err != nil {
			return err
		}
		data, err := Harden(b.cipher, idx)
		if err != nil {
			return err
		}
		return bucket.Put(indexKey, data)
	})
}
func (b BoltStorage) Rename(id uuid.UUID, newName string) (Index, error) {
	idx, err := b.Index()
	if err != nil {
		return idx, err
	}
	for i, meta := range idx {
		if meta.Id == id {
			idx[i].Name = newName
			return idx, b.updateIndex(idx)
		}
	}
	return idx, ErrNoSuchEntry
}

func (b BoltStorage) MoveUp(id uuid.UUID) (Index, error) {
	idx, err := b.Index()
	if err != nil {
		return idx, err
	}
	for i, entry := range idx {
		if entry.Id == id {
			if i == 0 {
				return idx, nil // Already being at the top is not an error
			}
			idx[i-1], idx[i] = idx[i], idx[i-1]
			err = b.updateIndex(idx)
			return idx, err
		}
	}
	return idx, ErrNoSuchEntry
}
func (b BoltStorage) MoveDown(id uuid.UUID) (Index, error) {
	idx, err := b.Index()
	if err != nil {
		return idx, err
	}
	for i, entry := range idx {
		if entry.Id == id {
			if i == len(idx)-1 {
				return idx, nil // Already being at the bottom is not an error
			}
			idx[i+1], idx[i] = idx[i], idx[i+1]
			err = b.updateIndex(idx)
			return idx, err
		}
	}
	return idx, ErrNoSuchEntry
}

func (b BoltStorage) Create(name, initialText string) (Entry, Index, error) {
	entry := Entry{
		Text: initialText,
	}

	idx, err := b.Index()
	if err != nil {
		return entry, idx, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return entry, idx, err
	}
	entry.Id = id

	idx = append(idx, EntryMeta{
		Name: name,
		Id:   entry.Id,
	})

	err = b.updateIndex(idx)
	if err != nil {
		return entry, idx, err
	}

	err = b.storeEntry(entry)

	return entry, idx, err
}

func (b BoltStorage) storeEntry(entry Entry) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := b.bucket(tx)
		if err != nil {
			return err
		}
		data, err := Harden(b.cipher, entry)
		if err != nil {
			return err
		}
		return bucket.Put(entry.Id[:], data)
	})
}

func (b BoltStorage) Read(id uuid.UUID) (Entry, error) {
	entry := Entry{}
	err := b.bolt.View(func(tx *bolt.Tx) error {
		bucket, err := b.bucket(tx)
		if err != nil {
			return err
		}
		raw := bucket.Get(id[:])
		if raw == nil {
			return ErrNoSuchEntry
		}
		return Soften(b.cipher, raw, &entry)

	})
	return entry, err
}

func (b BoltStorage) Update(entry Entry) error {
	idx, err := b.Index()
	if err != nil {
		return err
	}
	for _, candidate := range idx {
		if candidate.Id == entry.Id {
			return b.storeEntry(entry)
		}
	}
	return ErrNoSuchEntry
}

func (b BoltStorage) Delete(id uuid.UUID) (Index, error) {
	idx, err := b.Index()
	if err != nil {
		return idx, err
	}
	remove := -1
	for i, candidate := range idx {
		if candidate.Id == id {
			remove = i
			break
		}
	}
	if remove < 0 {
		return idx, ErrNoSuchEntry
	}
	idx = append(idx[:remove], idx[remove+1:]...)
	err = b.updateIndex(idx)
	if err != nil {
		return idx, err
	}
	return idx, b.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := b.bucket(tx)
		if err != nil {
			return err
		}
		return bucket.Delete(id[:])
	})
}
