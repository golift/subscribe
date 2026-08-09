package subscribe

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrSubscriberNotFound is returned any time a requested subscriber does not exist.
	ErrSubscriberNotFound = errors.New("subscriber not found")
	// ErrEventNotFound is returned when a requested event has not been created.
	ErrEventNotFound = errors.New("event not found")
	// ErrEventExists is returned when a new event with an existing name is created.
	ErrEventExists = errors.New("event already exists")
)

// Rules contains the pause time and rules for a subscriber's event subscription.
// Rules are unused by the library and available for consumers.
type Rules struct {
	Pause time.Time                `json:"pause"`
	D     map[string]time.Duration `json:"durations"`
	I     map[string]int           `json:"integers"`
	S     map[string]string        `json:"strings"`
	T     map[string]time.Time     `json:"times"`
}

// Subscriber describes the contact info and subscriptions for a person.
//
// Meta, Contact, Admin and Ignored change over a subscriber's life, and the
// library reads them from whichever goroutine calls GetSubscribers() and
// friends. Read and write those four through the accessor methods — GetMeta,
// SetContact, IsAdmin, SetIgnored and so on — which hold the lock below. The
// fields stay exported for JSON and for building a record with a struct
// literal; touching them directly on a record the library already holds races
// with the state file save.
type Subscriber struct {
	// ID is optional. If it provided, this is used as the _match_.
	ID int64 `json:"id"`
	// Meta is optional. This library does not use this value.
	// Use GetMeta/SetMeta/DeleteMeta.
	Meta map[string]any `json:"meta"`
	// API is the type of API the subscriber is subscribed with. Used to filter results.
	API string `json:"api"`
	// Contact is the contact info used in the API to send the subscriber a notification.
	// If ID is not present this value is used as the _match_.
	// Use GetContact/SetContact.
	Contact string `json:"contact"`
	// Events is a list of events the subscriber is subscribed to, including a cooldown/pause time.
	Events *Events `json:"events"`
	// This is just extra data that can be used to make the user special.
	// Use IsAdmin/SetAdmin.
	Admin bool `json:"isAdmin"`
	// Ignored will exclude a user from GetSubscribers().
	// Use IsIgnored/SetIgnored.
	Ignored bool `json:"ignored"`
	// FirstSeen is when this subscriber record was first created. Zero if unknown (legacy rows).
	FirstSeen time.Time `json:"firstSeen,omitzero"`
	// mu protects Meta, Contact, Admin and Ignored.
	mu sync.RWMutex
}

// Events represents the map of tracked global Events.
// This is an arbitrary list that can be used to filter
// notifications in a consuming application.
type Events struct {
	// Map is the events/rules map. Use the provided methods to interact with it.
	Map map[string]*Rules `json:"eventsMap"`
	// sync.mu locks and unlocks the Events map
	mu sync.RWMutex
}

// Subscribe is the data needed to initialize this module.
//
// Three locks guard the tree of data below, and code that needs more than one
// takes them outermost first: Subscribe.mu, then Subscriber.mu, then Events.mu.
// The state file save walks all three in that order.
type Subscribe struct {
	// EnableAPIs sets the allowed APIs. Only subscriptions that have an API
	// with a prefix in this list will return from the GetSubscribers() method.
	EnableAPIs []string `json:"enabledApis"` // imessage, skype, pushover, email, slack, growl, all, any
	// mu protects mutable Subscribe fields.
	mu sync.RWMutex
	// stateFile is the db location, like: /usr/local/var/lib/motifini/subscribers.json
	stateFile string
	// Events stores a list of arbitrary events. Use the included methods to interact with it.
	// This does not affect GetSubscribers(). Use the data here as a filter in your app.
	Events *Events `json:"events"`
	// Subscribers is a list of all Subscribers.
	Subscribers []*Subscriber `json:"subscribers"`
}
