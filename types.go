package subscribe

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrorSubscriberNotFound is returned any time a requested subscriber does not exist.
	ErrorSubscriberNotFound = errors.New("subscriber not found")
	// ErrorEventNotFound is returned when a requested event has not been created.
	ErrorEventNotFound = errors.New("event not found")
	// ErrorEventExists is returned when a new event with an existing name is created.
	ErrorEventExists = errors.New("event already exists")
)

// subEventInfo contains the pause time and rules for a subscriber's event subscription.
// Rules are unused by the library and available for consumers.
type subEventInfo struct {
	Pause time.Time `json:"pause"`
	Rules []string  `json:"rules"`
}

// subEvents represent's a subscriber's list of event subscriptions.
// Each subscription has a unique name and some meta data attached.
type subEvents struct {
	// Map is the events/info map. Use the provided methods to interact with it.
	Map map[string]subEventInfo `json:"events_map"`
	// sync.RWMutex Locks/Unlocks Events map
	sync.RWMutex
}

// Subscriber describes the contact info and subscriptions for a person.
type Subscriber struct {
	// API is the type of API the subscriber is subscribed with. Used to filter results.
	API string `json:"api"`
	// Contact is the contact info used in the API to send the subscriber a notification.
	Contact string `json:"contact"`
	// Events is a list of events the subscriber is subscribed to, including a cooldown/pause time.
	Events *subEvents `json:"events"`
	// This is just extra data that can be used to make the user special.
	Admin bool `json:"is_admin"`
	// Ignored will exclude a user from GetSubscribers().
	Ignored bool `json:"ignored"`
}

// Rules is arbitrary data that can be stored with an event.
type Rules map[string]string

// events represents the map of tracked global events.
// This is an arbitrary list that can be used to filter
// notifications in a consuming application.
type events struct {
	// Map is the events/rules map. Use the provided methods to interact with it.
	Map map[string]Rules `json:"events_map"`
	// sync.RWMutex locks and unlocks the Events map
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
	// This does not affect GetSubscribers(). Use the data here as a filter in your app.
	Events *events `json:"events"`
	// Subscribers is a list of all Subscribers.
	Subscribers []*Subscriber `json:"subscribers"`
}
