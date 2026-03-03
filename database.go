// Package subscribe provides a subscription management system.
package subscribe

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"time"
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
	s.mu.RLock()
	stateFile := s.stateFile
	s.mu.RUnlock()

	if stateFile == "" {
		return nil
	}

	// #nosec G304 -- state file path is user-configured on purpose.
	buf, err := os.ReadFile(stateFile)
	if os.IsNotExist(err) {
		return s.StateFileSave()
	}

	if err != nil {
		return fmt.Errorf("failed reading state file: %w", err)
	}

	loaded := new(Subscribe)

	err = json.Unmarshal(buf, loaded)
	if err != nil {
		return fmt.Errorf("failed decoding state file: %w", err)
	}

	normalizeLoadedState(loaded)

	s.mu.Lock()
	s.EnableAPIs = loaded.EnableAPIs
	s.Events = loaded.Events
	s.Subscribers = loaded.Subscribers
	s.mu.Unlock()

	return nil
}

// StateGetJSON returns the state data in json format.
func (s *Subscribe) StateGetJSON() (string, error) {
	snapshot := s.snapshot()

	b, err := json.Marshal(snapshot)

	return string(b), err
}

// StateFileSave writes out the state file.
func (s *Subscribe) StateFileSave() error {
	const stateFileMode = 0o600

	s.mu.RLock()
	stateFile := s.stateFile
	s.mu.RUnlock()

	if stateFile == "" {
		return nil
	}

	snapshot := s.snapshot()

	buf, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	err = os.WriteFile(stateFile, buf, stateFileMode)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// StateFileRelocate writes the state file to a new location.
func (s *Subscribe) StateFileRelocate(newPath string) error {
	s.mu.Lock()
	oldPath := s.stateFile
	s.stateFile = newPath
	s.mu.Unlock()

	err := s.StateFileLoad()
	if err != nil {
		s.mu.Lock()
		s.stateFile = oldPath
		s.mu.Unlock()
	}

	return err
}

func normalizeLoadedState(loaded *Subscribe) {
	if loaded.EnableAPIs == nil {
		loaded.EnableAPIs = make([]string, 0)
	}

	if loaded.Events == nil {
		loaded.Events = &Events{Map: make(map[string]*Rules)}
	} else {
		normalizeEvents(loaded.Events)
	}

	if loaded.Subscribers == nil {
		loaded.Subscribers = make([]*Subscriber, 0)
	}

	for _, sub := range loaded.Subscribers {
		if sub == nil {
			continue
		}

		if sub.Events == nil {
			sub.Events = &Events{Map: make(map[string]*Rules)}

			continue
		}

		normalizeEvents(sub.Events)
	}
}

func normalizeEvents(events *Events) {
	if events.Map == nil {
		events.Map = make(map[string]*Rules)
	}

	for key, rules := range events.Map {
		if rules == nil {
			events.Map[key] = &Rules{
				D: make(map[string]time.Duration),
				I: make(map[string]int),
				S: make(map[string]string),
				T: make(map[string]time.Time),
			}

			continue
		}

		if rules.D == nil {
			rules.D = make(map[string]time.Duration)
		}

		if rules.I == nil {
			rules.I = make(map[string]int)
		}

		if rules.S == nil {
			rules.S = make(map[string]string)
		}

		if rules.T == nil {
			rules.T = make(map[string]time.Time)
		}
	}
}

func (s *Subscribe) snapshot() *Subscribe {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := &Subscribe{
		EnableAPIs:  append(make([]string, 0, len(s.EnableAPIs)), s.EnableAPIs...),
		Events:      snapshotEvents(s.Events),
		Subscribers: make([]*Subscriber, 0, len(s.Subscribers)),
	}

	for _, sub := range s.Subscribers {
		out.Subscribers = append(out.Subscribers, snapshotSubscriber(sub))
	}

	return out
}

func snapshotSubscriber(sub *Subscriber) *Subscriber {
	if sub == nil {
		return nil
	}

	out := &Subscriber{
		ID:      sub.ID,
		API:     sub.API,
		Contact: sub.Contact,
		Events:  snapshotEvents(sub.Events),
		Admin:   sub.Admin,
		Ignored: sub.Ignored,
	}

	if sub.Meta != nil {
		out.Meta = make(map[string]any, len(sub.Meta))
		maps.Copy(out.Meta, sub.Meta)
	}

	return out
}

func snapshotEvents(events *Events) *Events {
	if events == nil {
		return &Events{Map: make(map[string]*Rules)}
	}

	events.mu.RLock()
	defer events.mu.RUnlock()

	out := &Events{Map: make(map[string]*Rules, len(events.Map))}
	for event, rules := range events.Map {
		out.Map[event] = cloneRules(rules)
	}

	return out
}
