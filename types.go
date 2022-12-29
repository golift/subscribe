package subscribe

import (
	"fmt"
	"sync"
	"time"
)

var (
	// ErrSubscriberNotFound is returned any time a requested subscriber does not exist.
	ErrSubscriberNotFound = fmt.Errorf("subscriber not found")
	// ErrEventNotFound is returned when a requested event has not been created.
	ErrEventNotFound = fmt.Errorf("event not found")
	// ErrEventExists is returned when a new event with an existing name is created.
	ErrEventExists = fmt.Errorf("event already exists")
)

// Rules contains the pause time and rules for a subscriber's event subscription.
// Rules are unused by the library and available for consumers.
type Rules struct {
	Pause time.Time `json:"pause"`
	D     map[string]time.Duration
	I     map[string]int
	S     map[string]string
	T     map[string]time.Time
}

// Subscriber describes the contact info and subscriptions for a person.
type Subscriber struct {
	// ID is optional. If it provided, this is used as the _match_.
	ID int64 `json:"id"`
	// Meta is optional. This library does not use this value.
	Meta map[string]interface{} `json:"meta"`
	// API is the type of API the subscriber is subscribed with. Used to filter results.
	API string `json:"api"`
	// Contact is the contact info used in the API to send the subscriber a notification.
	// If ID is not present this value is used as the _match_.
	Contact string `json:"contact"`
	// Events is a list of events the subscriber is subscribed to, including a cooldown/pause time.
	Events *Events `json:"events"`
	// This is just extra data that can be used to make the user special.
	Admin bool `json:"isAdmin"`
	// Ignored will exclude a user from GetSubscribers().
	Ignored bool `json:"ignored"`
}

// Events represents the map of tracked global Events.
// This is an arbitrary list that can be used to filter
// notifications in a consuming application.
type Events struct {
	// Map is the events/rules map. Use the provided methods to interact with it.
	Map map[string]*Rules `json:"eventsMap"`
	// sync.RWMutex locks and unlocks the Events map
	sync.RWMutex
}

// Subscribe is the data needed to initialize this module.
type Subscribe struct {
	// EnableAPIs sets the allowed APIs. Only subscriptions that have an API
	// with a prefix in this list will return from the GetSubscribers() method.
	EnableAPIs []string `json:"enabledApis"` // imessage, skype, pushover, email, slack, growl, all, any
	// stateFile is the db location, like: /usr/local/var/lib/motifini/subscribers.json
	stateFile string
	// Events stores a list of arbitrary events. Use the included methods to interact with it.
	// This does not affect GetSubscribers(). Use the data here as a filter in your app.
	Events *Events `json:"events"`
	// Subscribers is a list of all Subscribers.
	Subscribers []*Subscriber `json:"subscribers"`
}
