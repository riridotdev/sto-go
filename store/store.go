package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const storeFileName = ".sto"

type Store struct {
	EntryList []internalEntry
	Root      string `json:"-"`
}

func Init(path string) (Store, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return Store{}, ErrDirNotExist(path)
	}

	if !stat.Mode().IsDir() {
		return Store{}, ErrNotDir(path)
	}

	s := Store{Root: path}

	_, err = os.Stat(s.storeFilePath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Store{}, ErrStat{source: err, path: s.storeFilePath()}
	}
	if err == nil {
		return Store{}, ErrStoreAlreadyExist(path)
	}

	file, err := os.Create(s.storeFilePath())
	if err != nil {
		return Store{}, ErrCreatingStoreFile{path: s.storeFilePath(), source: err}
	}
	defer file.Close()

	return s, nil
}

func Restore(path string) (Store, error) {
	s := Store{Root: path}

	stat, err := os.Stat(s.storeFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Store{}, ErrStoreNotFound(path)
		}
	}

	if stat.Size() == 0 {
		return s, nil
	}

	file, err := os.Open(s.storeFilePath())
	if err != nil {
		return Store{}, ErrReadingFile{source: err, path: path}
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err := dec.Decode(&s); err != nil {
		return Store{}, ErrReadingFile{source: err, path: path}
	}

	return s, nil
}

func (s Store) Persist() error {
	file, err := os.OpenFile(s.storeFilePath(), os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return ErrWritingFile{source: err, path: s.storeFilePath()}
	}
	defer file.Close()

	buf := bytes.NewBuffer([]byte{})

	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return ErrWritingFile{source: err, path: s.storeFilePath()}
	}

	if _, err := buf.WriteTo(file); err != nil {
		return ErrWritingFile{source: err, path: s.storeFilePath()}
	}

	return nil
}

func (s Store) Entries() []Entry {
	entries := []Entry{}

	for _, entry := range s.EntryList {
		entries = append(entries, entry.toExternal(s.Root))
	}

	return entries
}

func (s *Store) Add(newEntry Entry) error {
	ie := internalEntry{
		Name:        newEntry.Name,
		Source:      newEntry.Source,
		Destination: newEntry.Destination,
	}

	if strings.HasPrefix(ie.Source, s.Root) {
		ie.Source = strings.TrimPrefix(ie.Source, s.Root)
		if ie.Source[0] == '/' {
			ie.Source = ie.Source[1:]
		}
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		panic(fmt.Sprintf("error reading home dir: %s", err))
	}

	if strings.HasPrefix(ie.Source, homeDir) {
		ie.Source = strings.Replace(ie.Source, homeDir, "~", 1)
	}
	if strings.HasPrefix(ie.Destination, homeDir) {
		ie.Destination = strings.Replace(ie.Destination, homeDir, "~", 1)
	}

	for _, entry := range s.EntryList {
		if entry.Source == newEntry.Source {
			return ErrDuplicateEntry{entry: newEntry, conflictingEntry: newEntry}
		}
	}

	s.EntryList = append(s.EntryList, ie)

	return nil
}

func (s *Store) Entry(name string) (Entry, error) {
	for _, entry := range s.EntryList {
		if entry.Name != name {
			continue
		}
		return entry.toExternal(s.Root), nil
	}

	return Entry{}, ErrEntryNotExist(name)
}

func (s Store) storeFilePath() string {
	return fmt.Sprintf("%s/%s", s.Root, storeFileName)
}

type internalEntry struct {
	Name        string
	Source      string
	Destination string
}

func (ie internalEntry) toExternal(root string) Entry {
	entry := Entry{Name: ie.Name, Source: ie.Source}

	var err error

	entry.Destination, err = homedir.Expand(ie.Destination)
	if err != nil {
		panic(fmt.Sprintf("error expanding homedir for %q: %s", ie.Destination, err))
	}
	entry.Source, err = homedir.Expand(ie.Source)
	if err != nil {
		panic(fmt.Sprintf("error expanding homedir for %q: %s", ie.Source, err))
	}

	if entry.Source[0] != '/' {
		entry.Source = fmt.Sprintf("%s/%s", root, ie.Source)
	}

	return entry
}

type Entry struct {
	Name        string
	Source      string
	Destination string
}

func (e Entry) Link() error {
	if _, err := os.Stat(e.Source); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotExists(e.Source)
		}
		return ErrStat{source: err, path: e.Source}
	}

	stat, err := os.Lstat(e.Destination)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ErrStat{source: err, path: e.Destination}
	}
	if err == nil {
		if stat.Mode()&os.ModeSymlink != os.ModeSymlink {
			return ErrItemExists(e.Destination)
		}
		resolvedPath, err := os.Readlink(e.Destination)
		if err != nil {
			return ErrReadingLink{source: err, path: e.Destination}
		}
		if resolvedPath != e.Source {
			return ErrItemExists(e.Destination)
		}
		return ErrAlreadyLinked(e)
	}

	if err := os.Symlink(e.Source, e.Destination); err != nil {
		return ErrLinking{source: err, sourcePath: e.Source, destinationPath: e.Destination}
	}

	return nil
}

func (e Entry) Unlink() error {
	linkState, err := e.Check()
	if err != nil {
		return err
	}

	if linkState != Linked {
		return nil
	}

	if err := os.Remove(e.Destination); err != nil {
		return ErrUnlinking{source: err, sourcePath: e.Source, destinationPath: e.Destination}
	}

	return nil
}

func (e Entry) Check() (EntryState, error) {
	stat, err := os.Lstat(e.Destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Unlinked, nil
		}
		return Error, ErrStat{source: err, path: e.Destination}
	}

	if stat.Mode()&os.ModeSymlink != os.ModeSymlink {
		return ConflictingFile, nil
	}

	resolvedPath, err := os.Readlink(e.Destination)
	if err != nil {
		return Error, ErrReadingLink{source: err, path: e.Destination}
	}

	if resolvedPath != e.Source {
		return ConflictingLink, nil
	}

	return Linked, nil
}

type EntryState int

const (
	_ EntryState = iota
	Error
	Linked
	Unlinked
	ConflictingFile
	ConflictingLink
)

func (es EntryState) Describe(e Entry) string {
	switch es {
	case ConflictingFile:
		return fmt.Sprintf(`conflicting file for entry %q
    %s -> %s
    existing entry at %q`, e.Name, e.Source, e.Destination, e.Destination)

	case ConflictingLink:
		resolvedPath, err := os.Readlink(e.Destination)
		if err != nil {
			panic(fmt.Sprintf("error resolving link at %q", e.Destination))
		}

		return fmt.Sprintf(`conflicting link for entry %q
    %s -> %s
existing link at
    %s -> %s`, e.Name, e.Source, e.Destination, resolvedPath, e.Destination)

	default:
		panic(fmt.Sprintf("should not be called with EntryState value %v", es))
	}
}
