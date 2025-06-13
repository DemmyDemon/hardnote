package storage_test

import (
	"testing"

	"github.com/DemmyDemon/hardnote/storage"
	"github.com/DemmyDemon/hardnote/test"
)

type SillyTestType struct {
	String  string
	Integer int
	Slice   []string
}

var testObject = SillyTestType{
	String:  "Hopefully the integer is 42",
	Integer: 42,
	Slice:   []string{"one", "two", "three"},
}

func TestEncodeDecode(t *testing.T) {
	data, err := storage.Encode(testObject)
	test.Result(t, err, "encoding test object", len(data))

	target := SillyTestType{}
	err = storage.Decode(data, &target)
	test.Result(t, err, "decoding test object", target)

	test.Compare(t, "compare initial object to decoded data", testObject, target)
}

func TestHardenSoften(t *testing.T) {
	gcm, err := storage.NewGCM([]byte("obvious-key"))
	test.Result(t, err, "instantiate GCM")

	babble, err := storage.Harden(gcm, testObject)
	test.Result(t, err, "hardening", len(babble))

	target := SillyTestType{}
	err = storage.Soften(gcm, babble, &target)
	test.Result(t, err, "softening", target)

	test.Compare(t, "compare initial object to softened data", testObject, target)
}
