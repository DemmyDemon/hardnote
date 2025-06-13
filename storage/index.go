package storage

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Index []EntryMeta
type EntryMeta struct {
	Name string
	Id   uuid.UUID
}

func (idx Index) String() string {
	var sb strings.Builder
	for i, element := range idx {
		sb.WriteString(fmt.Sprintf("%03d [%s] %s", i, element.Id, element.Name))
		if i < len(idx)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
