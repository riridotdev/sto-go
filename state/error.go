package state

import (
	"fmt"
)

type ErrReadingStateFile struct {
	source error
	path   string
}

func (e ErrReadingStateFile) Error() string {
	return fmt.Sprintf("error reading state file at %q: %s", e.path, e.source)
}

type ErrStateFileNotFound string

func (e ErrStateFileNotFound) Error() string {
	return fmt.Sprintf("state file at %q not found", string(e))
}

type ErrProfileWithNameAlreadyExists string

func (e ErrProfileWithNameAlreadyExists) Error() string {
	return fmt.Sprintf("profile with name %q already exists", string(e))
}

type ErrProfileWithRootAlreadyExists string

func (e ErrProfileWithRootAlreadyExists) Error() string {
	return fmt.Sprintf("profile with root %q already exists", string(e))
}

type ErrAccessingStateFile struct {
	source error
	path   string
}

func (e ErrAccessingStateFile) Error() string {
	return fmt.Sprintf("error accessing state file at %q: %s", e.path, e.source)
}

type ErrCreatingStateFile struct {
	source error
	path   string
}

func (e ErrCreatingStateFile) Error() string {
	return fmt.Sprintf("error creating state file at %q: %s", e.path, e.source)
}

type ErrWritingStateFile struct {
	source error
	path   string
}

func (e ErrWritingStateFile) Error() string {
	return fmt.Sprintf("error writing state file at %q: %s", e.path, e.source)
}

type ErrProfileNotFound string

func (e ErrProfileNotFound) Error() string {
	return fmt.Sprintf("profile %q not found", string(e))
}

type ErrDuplicateEntry string

func (e ErrDuplicateEntry) Error() string {
	return fmt.Sprintf("duplicate entry: %s", string(e))
}

type ErrEntryNotFound string

func (e ErrEntryNotFound) Error() string {
	return fmt.Sprintf("entry %q not found", string(e))
}
