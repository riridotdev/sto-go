package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/riridotdev/sto-go/store"
)

type State struct {
	DefaultProfile string
	Profiles       []Profile

	StateFilePath string `json:"-"`
}

func Restore(path string) (State, error) {
	stateFile, err := os.Open(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return State{}, ErrReadingStateFile{source: err, path: path}
		}
		return State{}, ErrStateFileNotFound(path)
	}

	state := State{StateFilePath: path}

	dec := json.NewDecoder(stateFile)
	if err := dec.Decode(&state); err != nil {
		return State{}, ErrReadingStateFile{source: err, path: path}
	}

	return state, nil
}

func (s *State) Persist() error {
	stateFile, err := os.OpenFile(s.StateFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return ErrAccessingStateFile{source: err, path: s.StateFilePath}
		}
		dir, _ := filepath.Split(s.StateFilePath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return ErrCreatingStateFile{source: err, path: s.StateFilePath}
		}
		if stateFile, err = os.Create(s.StateFilePath); err != nil {
			return ErrCreatingStateFile{source: err, path: s.StateFilePath}
		}
	}
	defer stateFile.Close()

	buf := bytes.NewBuffer([]byte{})

	enc := json.NewEncoder(buf)
	enc.SetIndent("", " ")
	if err := enc.Encode(s); err != nil {
		return ErrWritingStateFile{source: err, path: s.StateFilePath}
	}

	if _, err := buf.WriteTo(stateFile); err != nil {
		return ErrWritingStateFile{source: err, path: s.StateFilePath}
	}

	return nil
}

func (s *State) ActiveProfiles() []Profile {
	profiles := []Profile{}

	for _, profile := range s.Profiles {
		if !profile.Active {
			continue
		}
		profiles = append(profiles, profile)
	}

	return profiles
}

func (s *State) AddProfile(profile Profile) error {
	for _, p := range s.Profiles {
		if profile.Name == p.Name {
			return ErrProfileWithNameAlreadyExists(profile.Name)
		}
		if profile.Root == p.Root {
			return ErrProfileWithRootAlreadyExists(profile.Root)
		}
	}

	if s.DefaultProfile == "" {
		profile.Active = true
		s.DefaultProfile = profile.Name
	}

	s.Profiles = append(s.Profiles, profile)

	return nil
}

func (s *State) GetProfile(name string) (*Profile, error) {
	for i, profile := range s.Profiles {
		if profile.Name != name {
			continue
		}
		return &s.Profiles[i], nil
	}
	return nil, ErrProfileNotFound(name)
}

type Profile struct {
	Name          string
	Root          string
	Active        bool
	LinkedEntries []string
}

func (p *Profile) AddLinkedEntry(name string) error {
	for _, linkedEntry := range p.LinkedEntries {
		if linkedEntry == name {
			return ErrDuplicateEntry(name)
		}
	}

	p.LinkedEntries = append(p.LinkedEntries, name)

	return nil
}

func (p *Profile) RemoveLinkedEntry(name string) error {
	for i, linkedEntry := range p.LinkedEntries {
		if linkedEntry != name {
			continue
		}

		oldEntries := p.LinkedEntries
		p.LinkedEntries = p.LinkedEntries[:i]
		p.LinkedEntries = append(p.LinkedEntries, oldEntries[i+1:]...)

		return nil
	}

	return nil
}

func (p *Profile) Enable() error {
	s, err := store.Restore(p.Root)
	if err != nil {
		return err
	}

	for _, entryName := range p.LinkedEntries {
		entry, err := s.Entry(entryName)
		if err != nil {
			return err
		}

		entryState, err := entry.Check()
		if entryState != store.Unlinked {
			return fmt.Errorf("unable to enable profile %q:\n%s", p.Name, entryState.Describe(entry))
		}
	}

	for _, entryName := range p.LinkedEntries {
		entry, err := s.Entry(entryName)
		if err != nil {
			return err
		}

		if err := entry.Link(); err != nil {
			return err
		}
	}

	p.Active = true

	return nil
}

func (p *Profile) Disable() error {
	s, err := store.Restore(p.Root)
	if err != nil {
		return err
	}

	for _, entryName := range p.LinkedEntries {
		entry, err := s.Entry(entryName)
		if err != nil {
			return err
		}

		if err := entry.Unlink(); err != nil {
			return err
		}
	}

	p.Active = false

	return nil
}
