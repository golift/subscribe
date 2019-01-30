package subscribe

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// GetDB returns an interface to manage events
func GetDB(StateFile string) (*Subscribe, error) {
	s := &Subscribe{
		stateFile:   StateFile,
		EnableAPIs:  make([]string, 0),
		Events:      make(Events),
		Subscribers: make([]*Subscriber, 0),
	}
	return s, s.LoadStateFile()
}

// LoadStateFile data from a json file.
func (s *Subscribe) LoadStateFile() error {
	if s.stateFile == "" {
		return nil
	}
	if buf, err := ioutil.ReadFile(s.stateFile); os.IsNotExist(err) {
		return s.SaveStateFile()
	} else if err != nil {
		return err
	} else if err := json.Unmarshal(buf, s); err != nil {
		return err
	}
	return nil
}

// GetStateJSON returns the state data in json format.
func (s *Subscribe) GetStateJSON() (string, error) {
	s.RLock()
	defer s.RUnlock()
	b, err := json.Marshal(s)
	return string(b), err
}

// SaveStateFile writes out the state file.
func (s *Subscribe) SaveStateFile() error {
	if s.stateFile == "" {
		return nil
	}
	s.RLock()
	defer s.RUnlock()
	if buf, err := json.Marshal(s); err != nil {
		return err
	} else if err := ioutil.WriteFile(s.stateFile, buf, 0640); err != nil {
		return err
	}
	return nil
}
