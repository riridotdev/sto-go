package store

import (
	"fmt"
)

type ErrDirNotExist string

func (e ErrDirNotExist) Error() string {
	return fmt.Sprintf("directory %q does not exist", string(e))
}

type ErrNotDir string

func (e ErrNotDir) Error() string {
	return fmt.Sprintf("target at path %q es not a directory", string(e))
}

type ErrStat struct {
	source error
	path   string
}

func (e ErrStat) Error() string {
	return fmt.Sprintf("error reading file stat for %q: %s", e.path, e.source)
}

type ErrStoreAlreadyExist string

func (e ErrStoreAlreadyExist) Error() string {
	return fmt.Sprintf("store at %q already exists", string(e))
}

type ErrStoreNotExist string

func (e ErrStoreNotExist) Error() string {
	return fmt.Sprintf("no store found at %q", string(e))
}

type ErrCreatingStoreFile struct {
	source error
	path   string
}

func (e ErrCreatingStoreFile) Error() string {
	return fmt.Sprintf("error creating store file at %q: %s", e.path, e.source)
}

type ErrDuplicateEntry struct {
	entry            Entry
	conflictingEntry Entry
}

func (e ErrDuplicateEntry) Error() string {
	return fmt.Sprintf(
		"duplicate entry for %q, %q -> %q",
		e.entry.Source, e.conflictingEntry.Source, e.conflictingEntry.Destination,
	)
}

type ErrItemExists string

func (e ErrItemExists) Error() string {
	return fmt.Sprintf("item already exists at %q", string(e))
}

type ErrNotExists string

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("item at %q does not exist", string(e))
}

type ErrLinking struct {
	source          error
	sourcePath      string
	destinationPath string
}

func (e ErrLinking) Error() string {
	return fmt.Sprintf(
		"error creating symlink %q -> %q: %s",
		e.sourcePath, e.destinationPath, e.source,
	)
}

type ErrUnlinking struct {
	source          error
	sourcePath      string
	destinationPath string
}

func (e ErrUnlinking) Error() string {
	return fmt.Sprintf(
		"error removing symlink %q -> %q: %s",
		e.sourcePath, e.destinationPath, e.source,
	)
}

type ErrReadingFile struct {
	source error
	path   string
}

func (e ErrReadingFile) Error() string {
	return fmt.Sprintf("error reading file at %q: %s", e.path, e.source)
}

type ErrWritingFile struct {
	source error
	path   string
}

func (e ErrWritingFile) Error() string {
	return fmt.Sprintf("error writing file at %q: %s", e.path, e.source)
}

type ErrEntryNotExist string

func (e ErrEntryNotExist) Error() string {
	return fmt.Sprintf("entry %q not found", string(e))
}

type ErrReadingLink struct {
	source error
	path   string
}

func (e ErrReadingLink) Error() string {
	return fmt.Sprintf("error reading link at %q: %s", e.path, e.source)
}

type ErrAlreadyLinked Entry

func (e ErrAlreadyLinked) Error() string {
	return fmt.Sprintf("entry %q -> %q is already linked", e.Source, e.Destination)
}

type ErrStoreNotFound string

func (e ErrStoreNotFound) Error() string {
	return fmt.Sprintf("store at %q not found", string(e))
}
