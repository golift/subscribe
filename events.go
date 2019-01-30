package subscribe

/************************
 *   Events Methods   *
 ************************/

// GetEvents returns all the configured events.
func (s *Subscribe) GetEvents() Events {
	s.RLock()
	defer s.RUnlock()
	return s.Events
}

// GetEvent returns the rules for an event.
// Rules can be used by the user for whatever they want.
func (s *Subscribe) GetEvent(name string) (map[string]string, error) {
	s.RLock()
	defer s.RUnlock()
	if rules, ok := s.Events[name]; ok {
		return rules, nil
	}
	return nil, ErrorEventNotFound
}

// UpdateEvent adds or updates an event.
func (s *Subscribe) UpdateEvent(name string, rules map[string]string) bool {
	s.Lock()
	defer s.Unlock()
	if rules == nil {
		rules = make(map[string]string)
	}
	if _, ok := s.Events[name]; !ok {
		s.Events[name] = rules
		return true
	}
	for ruleName, rule := range rules {
		if rule == "" {
			delete(s.Events[name], ruleName)
		} else {
			s.Events[name][ruleName] = rule
		}
	}
	return false
}

// RemoveEvent obliterates an event and all subsciptions for it.
func (s *Subscribe) RemoveEvent(name string) (removed int) {
	s.Lock()
	delete(s.Events, name)
	s.Unlock()
	for i := range s.Subscribers {
		if _, ok := s.Subscribers[i].Events[name]; ok {
			s.Subscribers[i].Lock()
			delete(s.Subscribers[i].Events, name)
			s.Subscribers[i].Unlock()
			removed++
		}
	}
	return
}
