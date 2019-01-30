package subscribe

import (
	"strings"
	"time"
)

/****************************
 *   Subscription Methods   *
 ****************************/

// Subscribe adds an event subscription to a subscriber.
func (s *Subscriber) Subscribe(eventName string) error {
	s.Lock()
	defer s.Unlock()
	info, ok := s.Events[eventName]
	if ok {
		return ErrorEventExists
	}
	info.Pause = time.Now()
	s.Events[eventName] = info
	return nil
}

// UnSubscribe a subscriber from an event subscription.
func (s *Subscriber) UnSubscribe(eventName string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Events[eventName]; !ok {
		return ErrorEventNotFound
	}
	delete(s.Events, eventName)
	return nil
}

// Pause (or unpause with 0 duration) a subscriber's event subscription.
func (s *Subscriber) Pause(eventName string, duration time.Duration) error {
	s.Lock()
	defer s.Unlock()
	info, ok := s.Events[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	info.Pause = time.Now().Add(duration)
	s.Events[eventName] = info
	return nil
}

// Subscriptions returns a subscriber's event subscriptions.
func (s *Subscriber) Subscriptions() (events map[string]SubEventInfo) {
	s.Lock()
	defer s.Unlock()
	events = s.Events
	return
} /* not tested */

// GetSubscribers returns a list of valid event subscribers.
func (s *Subscribe) GetSubscribers(eventName string) (subscribers []*Subscriber) {
	for i := range s.Subscribers {
		if s.Subscribers[i].Ignored {
			continue
		}
		for event, evnData := range s.Subscribers[i].Events {
			if event == eventName && evnData.Pause.Before(time.Now()) && checkAPI(s.Subscribers[i].API, s.EnableAPIs) {
				subscribers = append(subscribers, s.Subscribers[i])
			}
		}
	}
	return
}

// checkAPI just looks for a string in a slice of strings with a twist.
func checkAPI(s string, slice []string) bool {
	if len(slice) < 1 {
		return true
	}
	for _, v := range slice {
		if v == s || strings.HasPrefix(s, v) || v == "all" || v == "any" {
			return true
		}
	}
	return false
}
