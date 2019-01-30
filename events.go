package subscribe

/************************
 *    Events Methods    *
 ************************/

// Names returns all the configured event namee.
func (e *events) Names() []string {
	e.RLock()
	defer e.RUnlock()
	names := []string{}
	for name := range e.Map {
		names = append(names, name)
	}
	return names
}

// Get returns the rules for an event.
// Rules can be used by the user for whatever they want.
func (e *events) Get(name string) (Rules, error) {
	e.RLock()
	defer e.RUnlock()
	if rules, ok := e.Map[name]; ok {
		return rules, nil
	}
	return nil, ErrorEventNotFound
}

// Update adds or updates an event.
func (e *events) Update(name string, rules Rules) bool {
	e.Lock()
	defer e.Unlock()
	if rules == nil {
		rules = make(Rules)
	}
	if _, ok := e.Map[name]; !ok {
		e.Map[name] = rules
		return true
	}
	for ruleName, rule := range rules {
		if rule == "" {
			delete(e.Map[name], ruleName)
		} else {
			e.Map[name][ruleName] = rule
		}
	}
	return false
}

// Remove deletes an event, and orphans any subscriptions.
func (e *events) Remove(name string) {
	e.Lock()
	delete(e.Map, name)
	e.Unlock()
}

// EventRemove obliterates an event and all subsciptions for it.
func (s *Subscribe) EventRemove(name string) (removed int) {
	s.Events.Remove(name)
	for i := range s.Subscribers {
		s.Subscribers[i].Events.Lock()
		if _, ok := s.Subscribers[i].Events.Map[name]; ok {
			delete(s.Subscribers[i].Events.Map, name)
			removed++
		}
		s.Subscribers[i].Events.Unlock()
	}
	return
}
