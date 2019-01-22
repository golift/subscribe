package subscribe

import (
	"sync"
	"time"
)

// Error enables constant errors.
type Error string

// Error allows a string to satisfy the error type.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrorSubscriberNotFound is returned any time a requested subscriber does not exist.
	ErrorSubscriberNotFound = Error("subscriber not found")
	// ErrorEventNotFound is returned when a requested event has not been created.
	ErrorEventNotFound = Error("event not found")
	// ErrorEventExists is returned when a new event with an existing name is created.
	ErrorEventExists = Error("event already exists")
)

type subscriberEvents map[string]time.Time
type subscribeEvents map[string]map[string]string

// Subscriber describes the contact info and subscriptions for a person.
type Subscriber struct {
	// API is the type of API the subscriber is subscribed with.
	API string `json:"api"`
	// Contact is the contact info used in the API to send the subscriber a notification.
	Contact string `json:"contact"`
	// Events is a list of events the subscriber is subscribed to, including a cooldown/pause time.
	Events subscriberEvents `json:"events"`
	// This is just extra data that can be used to make the user special.
	Admin bool `json:"is_admin"`
	// Ignored will exclude a user from GetSubscribers().
	Ignored bool `json:"ignored"`
	// sync.RWMutex Locks/UnlocksE vents map
	sync.RWMutex
}

// Subscribe is the data needed to initialize this module.
type Subscribe struct {
	// EnableAPIs sets the allowed APIs. Only subscriptions that have an API
	// with a prefix in this list will return from the GetSubscribers() method.
	EnableAPIs []string `json:"enabled_apis"` // imessage, skype, pushover, email, slack, growl, all, any
	// stateFile is the db location, like: /usr/local/var/lib/motifini/subscribers.json
	stateFile string
	// Events stores a list of arbitrary events. Use the included methods to interact with it.
	// This does not affect GetSubscribers().
	Events subscribeEvents `json:"events"`
	// Subscribers is a list of all Subscribers.
	Subscribers []*Subscriber `json:"subscribers"`
	// sync.RWMutex locks and unlocks the Events map
	sync.RWMutex
}
