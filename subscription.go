package subscribe

import (
	"strings"
	"time"
)

/****************************
 *   Subscription Methods   *
 ****************************/

// Subscribe adds an event subscription to a subscriber.
// Returns an error only if the event subscription already exists.
func (s *Subscriber) Subscribe(event string) error {
	return s.Events.New(event, &Rules{Pause: time.Now()})
}

// GetSubscribers returns a list of valid event subscribers.
// This is the main method that should be triggered when an event occurs.
// Call this method when your event fires, collect the subscribers and send
// them notifications in your app. Subscribers can be people. Or functions.
func (s *Subscribe) GetSubscribers(eventName string) (subscribers []*Subscriber) {
	for _, sub := range s.Subscribers {
		if !sub.Ignored && s.checkAPI(sub.API) && !sub.Events.IsPaused(eventName) {
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

// EventRemove obliterates an event and all subsciptions for it.
func (s *Subscribe) EventRemove(event string) {
	s.Events.Remove(event)
	for _, sub := range s.Subscribers {
		sub.Events.Remove(event)
	}
}
