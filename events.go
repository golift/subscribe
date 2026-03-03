package subscribe

import (
	"maps"
	"sort"
	"strings"
	"time"
)

/************************
 *    Events Methods    *
 ************************/

// Names returns all the configured event names.
func (e *Events) Names() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	names := make([]string, 0, len(e.Map))

	for name := range e.Map {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Len returns the number of configured events.
func (e *Events) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return len(e.Map)
}

// Name finds an event case insensitively.
func (e *Events) Name(event string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, ok := e.Map[event]; ok {
		return event
	}

	for k := range e.Map {
		if strings.EqualFold(k, event) {
			return k
		}
	}

	return ""
}

// Exists returns true if an event exists.
func (e *Events) Exists(event string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, ok := e.Map[event]; ok {
		return true
	}

	return false
}

// New adds an event.
func (e *Events) New(event string, rules *Rules) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; ok {
		return ErrEventExists
	}

	e.Map[event] = cloneRules(rules)

	return nil
}

// UnPause resumes a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (e *Events) UnPause(event string) error {
	return e.Pause(event, 0)
}

// Pause (or unpause with 0 duration) a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (e *Events) Pause(event string, duration time.Duration) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return ErrEventNotFound
	}

	e.Map[event].Pause = time.Now().Add(duration)

	return nil
}

// IsPaused returns true if the event's notifications are paused.
// Returns true if the event subscription does not exist.
func (e *Events) IsPaused(event string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	info, ok := e.Map[event]
	if !ok {
		return true
	}

	return info.Pause.After(time.Now())
}

// PauseTime returns the pause time for an event.
func (e *Events) PauseTime(event string) time.Time {
	e.mu.RLock()
	defer e.mu.RUnlock()

	info, ok := e.Map[event]
	if !ok {
		return time.Time{}
	}

	return info.Pause
}

// Remove deletes an event.
func (e *Events) Remove(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.Map, event)
}

// RuleGetD returns a Duration rule.
func (e *Events) RuleGetD(event, rule string) (time.Duration, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules, found := e.Map[event]
	if !found || rules == nil {
		return 0, false
	}

	val, found := rules.D[rule]

	return val, found
}

// RuleGetI returns an integer rule.
func (e *Events) RuleGetI(event, rule string) (int, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules, found := e.Map[event]
	if !found || rules == nil {
		return 0, false
	}

	val, found := rules.I[rule]

	return val, found
}

// RuleGetS returns a string rule.
func (e *Events) RuleGetS(event, rule string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules, found := e.Map[event]
	if !found || rules == nil {
		return "", false
	}

	val, found := rules.S[rule]

	return val, found
}

// RuleGetT returns a Time rule.
func (e *Events) RuleGetT(event, rule string) (time.Time, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules, found := e.Map[event]
	if !found || rules == nil {
		return time.Now(), false
	}

	val, found := rules.T[rule]

	return val, found
}

// RuleSetD updates or sets a Duration rule.
func (e *Events) RuleSetD(event, rule string, val time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return
	}

	if e.Map[event].D == nil {
		e.Map[event].D = make(map[string]time.Duration)
	}

	e.Map[event].D[rule] = val
}

// RuleSetI updates or sets an integer rule.
func (e *Events) RuleSetI(event, rule string, val int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return
	}

	if e.Map[event].I == nil {
		e.Map[event].I = make(map[string]int)
	}

	e.Map[event].I[rule] = val
}

// RuleSetS updates or sets a string rule.
func (e *Events) RuleSetS(event, rule, val string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return
	}

	if e.Map[event].S == nil {
		e.Map[event].S = make(map[string]string)
	}

	e.Map[event].S[rule] = val
}

// RuleSetT updates or sets a Time rule.
func (e *Events) RuleSetT(event, rule string, val time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return
	}

	if e.Map[event].T == nil {
		e.Map[event].T = make(map[string]time.Time)
	}

	e.Map[event].T[rule] = val
}

// RuleDelD deletes a Duration rule.
func (e *Events) RuleDelD(event, rule string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok || e.Map[event].D == nil {
		return
	}

	delete(e.Map[event].D, rule)
}

// RuleDelI deletes an integer rule.
func (e *Events) RuleDelI(event, rule string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok || e.Map[event].I == nil {
		return
	}

	delete(e.Map[event].I, rule)
}

// RuleDelS deletes a string rule.
func (e *Events) RuleDelS(event, rule string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok || e.Map[event].S == nil {
		return
	}

	delete(e.Map[event].S, rule)
}

// RuleDelT deletes a Time rule.
func (e *Events) RuleDelT(event, rule string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok || e.Map[event].T == nil {
		return
	}

	delete(e.Map[event].T, rule)
}

// RuleDelAll deletes rules of any type with a specific name.
func (e *Events) RuleDelAll(event, rule string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[event]; !ok {
		return
	}

	delete(e.Map[event].D, rule)
	delete(e.Map[event].I, rule)
	delete(e.Map[event].S, rule)
	delete(e.Map[event].T, rule)
}

func cloneRules(rules *Rules) *Rules {
	if rules == nil {
		return &Rules{
			D: make(map[string]time.Duration),
			I: make(map[string]int),
			S: make(map[string]string),
			T: make(map[string]time.Time),
		}
	}

	cloned := &Rules{
		Pause: rules.Pause,
		D:     make(map[string]time.Duration, len(rules.D)),
		I:     make(map[string]int, len(rules.I)),
		S:     make(map[string]string, len(rules.S)),
		T:     make(map[string]time.Time, len(rules.T)),
	}

	maps.Copy(cloned.D, rules.D)
	maps.Copy(cloned.I, rules.I)
	maps.Copy(cloned.S, rules.S)
	maps.Copy(cloned.T, rules.T)

	return cloned
}
