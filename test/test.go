package test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Smush(parts ...any) string {
	if parts == nil || len(parts) == 0 {
		return "{nothing}"
	}
	var sb strings.Builder
	for i, part := range parts {
		str := fmt.Sprintf("%v", part)
		if strings.Contains(str, "\n") {
			sb.WriteString("\n\t")
		}
		str = strings.ReplaceAll(str, "\n", "\n\t")
		if part == nil {
			str = "nil"
		}
		if str == "" {
			str = "{empty}"
		}
		sb.WriteString(str)
		if i < len(parts)-1 {
			sb.WriteString(" | ")
		}
	}
	return sb.String()
}

func Result(t *testing.T, err error, operation string, more ...any) {
	_, file, no, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("TESTING IS RUINED BECAUSE THE RESULT CHECKER COULDN'T GET THE CALL STACK")
	}
	file = filepath.Base(file)
	if err == nil {
		fmt.Printf("✓ %s:%d %s → %s\n", file, no, operation, Smush(more...))
		return
	}
	fmt.Printf("✗ %s:%d %s: %s → %s\n", file, no, operation, err, Smush(more...))
	t.FailNow()
}

func Compare[T any](t *testing.T, context string, one T, other T) {
	_, file, no, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("TESTING IS RUINED BECAUSE THE VALUE COMPARATOR COULDN'T GET THE CALL STACK")
	}
	file = filepath.Base(file)
	if reflect.DeepEqual(one, other) {
		fmt.Printf("✓ %s:%d %s → compared values are DeepEqual\n", file, no, context)
		return
	}
	fmt.Printf("✗ %s:%d %s → compared values DIFFER\n\t%#v\n\t%#v\n", file, no, context, one, other)
	t.Fail()
}
