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

	return sub, sub.StateFileLoad()
}

// StateFileLoad data from a json file.
func (s *Subscribe) StateFileLoad() error {
	if s.stateFile == "" {
		return nil
	}

	if buf, err := os.ReadFile(s.stateFile); os.IsNotExist(err) {
		return s.StateFileSave()
	} else if err != nil {
		return fmt.Errorf("file problem: %w", err)
	} else if err := json.Unmarshal(buf, s); err != nil {
		return fmt.Errorf("json problem: %w", err)
	}

	return nil
}

// StateGetJSON returns the state data in json format.
func (s *Subscribe) StateGetJSON() (string, error) {
	s.Events.RLock()
	defer s.Events.RUnlock()

	b, err := json.Marshal(s)

	return string(b), err
}

// StateFileSave writes out the state file.
func (s *Subscribe) StateFileSave() error {
	if s.stateFile == "" {
		return nil
	}

	s.Events.RLock()
	defer s.Events.RUnlock()

	if buf, err := json.Marshal(s); err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	} else if err = os.WriteFile(s.stateFile, buf, 0o600); err != nil { //nolint:gomnd
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// StateFileRelocate writes the state file to a new location.
func (s *Subscribe) StateFileRelocate(newPath string) error {
	s.stateFile, newPath = newPath, s.stateFile // swap places

	if err := s.StateFileLoad(); err != nil {
		s.stateFile = newPath // got an error, put it back.
		return err
	}

	return nil
}
