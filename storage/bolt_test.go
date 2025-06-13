package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/DemmyDemon/hardnote/storage"
	"github.com/DemmyDemon/hardnote/test"
)

func TestBolt(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "hardnote.test")
	store, err := storage.NewBoltStorage(filename, []byte("Please don't tell anyone my secret key!"))
	test.Result(t, err, "open file", filename)

	defer func() {
		err = store.Close()
		test.Result(t, err, "close file", filename)
	}()

	idx, err := store.Index()
	test.Result(t, err, "read empty index", idx)

	entry, idx, err := store.Create("Test entry name", "Test entry body")
	test.Result(t, err, "create initial entry", entry, idx)

	entry.Text = "Updated text body of the entry"
	err = store.Update(entry)
	test.Result(t, err, "update initial entry", entry)

	compareEntry, err := store.Read(entry.Id)
	test.Result(t, err, "read initial entry back", compareEntry)
	test.Compare(t, "compare read entry to original", entry, compareEntry)

	idx, err = store.Rename(compareEntry.Id, "Renamed")
	test.Result(t, err, "rename initial entry", idx)

	additionalEntry, idx, err := store.Create("Additional test entry name", "Additional test entry body")
	test.Result(t, err, "create additional entry", additionalEntry, idx)

	idx, err = store.MoveUp(additionalEntry.Id)
	test.Result(t, err, "move additional entry up", idx)

	idx, err = store.MoveUp(additionalEntry.Id)
	test.Result(t, err, "move top entry up", idx)

	idx, err = store.MoveDown(additionalEntry.Id)
	test.Result(t, err, "move additional entry down", idx)

	idx, err = store.MoveDown(additionalEntry.Id)
	test.Result(t, err, "move bottom entry down", idx)

	idx, err = store.Delete(compareEntry.Id)
	test.Result(t, err, "delete initial entry", idx)

	idx, err = store.Delete(additionalEntry.Id)
	test.Result(t, err, "delete additional entry", idx)
}
