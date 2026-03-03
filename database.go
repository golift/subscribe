// Package subscribe provides a subscription management system.
package subscribe

import (
	"encoding/json"
	"fmt"
	"os"
)

/************************
 *   Database Methods   *
 ************************/

// GetDB returns an interface to manage events.
func GetDB(stateFile string) (*Subscribe, error) {
	sub := &Subscribe{
		stateFile:   stateFile,
		EnableAPIs:  make([]string, 0),
		Events:      &Events{Map: make(map[string]*Rules)},
		Subscribers: make([]*Subscriber, 0),
	}

	err := sub.StateFileLoad()
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// StateFileLoad data from a json file.
func (s *Subscribe) StateFileLoad() error {
	if s.stateFile == "" {
		return nil
	}

	buf, err := os.ReadFile(s.stateFile)
	if os.IsNotExist(err) {
		return s.StateFileSave()
	}

	if err != nil {
		return fmt.Errorf("failed reading state file: %w", err)
	}

	err = json.Unmarshal(buf, s)
	if err != nil {
		return fmt.Errorf("failed decoding state file: %w", err)
	}

	return nil
}

// StateGetJSON returns the state data in json format.
func (s *Subscribe) StateGetJSON() (string, error) {
	s.Events.mu.RLock()
	defer s.Events.mu.RUnlock()

	b, err := json.Marshal(s)

	return string(b), err
}

// StateFileSave writes out the state file.
func (s *Subscribe) StateFileSave() error {
	const stateFileMode = 0o600

	if s.stateFile == "" {
		return nil
	}

	s.Events.mu.RLock()
	defer s.Events.mu.RUnlock()

	buf, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	err = os.WriteFile(s.stateFile, buf, stateFileMode)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// StateFileRelocate writes the state file to a new location.
func (s *Subscribe) StateFileRelocate(newPath string) error {
	s.stateFile, newPath = newPath, s.stateFile // swap places

	err := s.StateFileLoad()
	if err != nil {
		s.stateFile = newPath // got an error, put it back.
	}

	return err
}
