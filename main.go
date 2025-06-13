package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DemmyDemon/hardnote/storage"
	"github.com/DemmyDemon/hardnote/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

func must(code int, what string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", what, err)
	os.Exit(code)
}
func same[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Specify a filename!")
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Println("SECURITY NOTE: KEY AND CURRENT NOTE ARE UNENCRYPTED IN MEMORY!")
	fmt.Println("DO NOT ENTER YOUR PASSPHEASE IN AN UNTRUSTED ENVIRONMENT!")
	fmt.Printf("Enter passprase for %s> ", filepath.Base(filename))

	key, err := term.ReadPassword(os.Stdin.Fd())
	fmt.Println("")
	must(2, "Reading password failed", err)

	_, err = os.Stat(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			must(3, "Could not get information about the specified file", err)
		}
		fmt.Println("This is a new file. Please repeat the passphrase.")
		fmt.Printf("Enter passprase for %s> ", filepath.Base(filename))
		keyAgain, err := term.ReadPassword(os.Stdin.Fd())
		fmt.Println("")
		must(4, "Reading password failed", err)
		if !same(key, keyAgain) {
			must(5, "Could not create file", errors.New("passwords did not match"))
		}
	}

	store, err := storage.NewBoltStorage(filename, key)
	must(6, "Could not open storage", err)
	defer func() {
		err := store.Close()
		must(7, "Error while closing storage", err)
		fmt.Println("\033c")
		fmt.Println("OKAY BYE!")
	}()

	p := tea.NewProgram(ui.New(filepath.Base(filename), store))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "OH NO, I TOTALLY %v\n", err)
		os.Exit(8)
	}
}
