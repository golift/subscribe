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

// Update adds or updates an event. Returns true for new events.
// Returns false if an existing event's rules were updated.
// Empty rules are removed. Existing rules are replaced.
func (e *events) Update(name string, rules Rules) (new bool) {
	e.Lock()
	defer e.Unlock()
	if rules == nil {
		rules = make(Rules)
	}
	if _, ok := e.Map[name]; !ok {
		e.Map[name] = rules
		new = true
	}
	for ruleName, rule := range rules {
		if rule == "" {
			delete(e.Map[name], ruleName)
		} else {
			e.Map[name][ruleName] = rule
		}
	}
	return
}

// Remove deletes an event, and orphans any subscriptions.
func (e *events) Remove(name string) {
	e.Lock()
	defer e.Unlock()
	delete(e.Map, name)
}

// EventRemove obliterates an event and all subsciptions for it.
// Returns the number of subscriptions removed.
func (s *Subscribe) EventRemove(eventName string) (removed int) {
	s.Events.Remove(eventName)
	for _, sub := range s.Subscribers {
		if sub.UnSubscribe(eventName) == nil {
			removed++
		}
	}
	return
}
