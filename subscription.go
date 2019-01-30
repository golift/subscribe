package subscribe

import (
	"sort"
	"strings"
	"time"
)

/****************************
 *   Subscription Methods   *
 ****************************/

// Subscribe adds an event subscription to a subscriber.
// Returns an error only if the event subscription already exists.
func (s *Subscriber) Subscribe(eventName string) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	info, ok := s.Events.Map[eventName]
	if ok {
		return ErrorEventExists
	}
	info.Pause = time.Now()
	s.Events.Map[eventName] = info
	return nil
}

// UnSubscribe a subscriber from an event subscription.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) UnSubscribe(eventName string) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	if _, ok := s.Events.Map[eventName]; !ok {
		return ErrorEventNotFound
	}
	delete(s.Events.Map, eventName)
	return nil
}

// UnPause resumes a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) UnPause(eventName string) error {
	return s.Pause(eventName, 0)
}

// Pause (or unpause with 0 duration) a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) Pause(eventName string, duration time.Duration) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	info, ok := s.Events.Map[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	info.Pause = time.Now().Add(duration)
	s.Events.Map[eventName] = info
	return nil
}

// IsPaused returns true if the event's notifications are pasued.
// Returns true if the event subscription does not exist.
func (s *Subscriber) IsPaused(eventName string) bool {
	s.Events.RLock()
	defer s.Events.RUnlock()
	info, ok := s.Events.Map[eventName]
	if !ok {
		return true
	}
	return info.Pause.After(time.Now())
}

// RuleExists returns true if an event rule exists.
// Returns false if the event subscription or rule do not exist.
func (s *Subscriber) RuleExists(eventName string, ruleName string) bool {
	s.Events.RLock()
	defer s.Events.RUnlock()
	if info, ok := s.Events.Map[eventName]; ok {
		for _, r := range info.Rules {
			if r == ruleName {
				return true
			}
		}
	}
	return false
}

// RulesGet returns a subscriber's event subscription rules.
// Returns an empty slice if the event subscription does not exist (or has no rules).
func (s *Subscriber) RulesGet(eventName string) []string {
	s.Events.RLock()
	defer s.Events.RUnlock()
	info, ok := s.Events.Map[eventName]
	if !ok {
		return []string{}
	}
	return info.Rules
}

// RulesReplace replaces a subscriber's event subscription rules with new rules.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) RulesReplace(eventName string, newRules []string) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	info, ok := s.Events.Map[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	info.Rules = newRules
	s.Events.Map[eventName] = info
	return nil
}

// RulesAdd appends new rule(s) to a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) RulesAdd(eventName string, appendRules []string) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	info, ok := s.Events.Map[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	info.Rules = append(info.Rules, appendRules...)
	s.Events.Map[eventName] = info
	return nil
}

// RulesRemove removes a rule from a subscriber's event subscription.
// Returns an error only if the event subscription is not found.
func (s *Subscriber) RulesRemove(eventName string, rule string) error {
	s.Events.Lock()
	defer s.Events.Unlock()
	var newRules []string
	info, ok := s.Events.Map[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	for _, r := range info.Rules {
		if r != rule {
			newRules = append(newRules, r)
		}
	}
	info.Rules = newRules
	s.Events.Map[eventName] = info
	return nil
}

// Subscriptions returns a subscriber's event subscriptions (names).
// Returns an empty slice if there are no subscriptions. Check the len().
func (s *Subscriber) Subscriptions() []string {
	s.Events.RLock()
	defer s.Events.RUnlock()
	names := []string{}
	for name := range s.Events.Map {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetSubscribers returns a list of valid event subscribers.
// This is the main method that should be triggered when an event occurs.
// Call this method when your event fires, collect the subscribers and send
// them notifications in your app. Subscribers can be people. Or functions.
func (s *Subscribe) GetSubscribers(eventName string) (subscribers []*Subscriber) {
	for _, sub := range s.Subscribers {
		if !sub.Ignored && s.checkAPI(sub.API) && !sub.IsPaused(eventName) {
			subscribers = append(subscribers, sub)
		}
	}
	return
}

// checkAPI just looks for a string in a slice of strings with a twist.
func (s *Subscribe) checkAPI(api string) bool {
	if len(s.EnableAPIs) < 1 {
		return true
	}
	for _, a := range s.EnableAPIs {
		if a == api || strings.HasPrefix(api, a) || a == "all" || a == "any" {
			return true
		}
	}
	return false
}
