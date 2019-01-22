package subscribe

/* Subscriptions Library!
    Reasonably Generic and fully tested. May work in your application!
		Check out the interfaces in types.go to get an idea how it works.
*/
import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// GetDB returns an interface to manage events
func GetDB(StateFile string) (*Subscribe, error) {
	s := &Subscribe{
		stateFile:   StateFile,
		EnableAPIs:  make([]string, 0),
		Events:      make(map[string]map[string]string),
		Subscribers: make([]*Subscriber, 0),
	}
	return s, s.LoadStateFile()
}

// LoadStateFile data from a json file.
func (s *Subscribe) LoadStateFile() error {
	if s.stateFile == "" {
		return nil
	}
	if buf, err := ioutil.ReadFile(s.stateFile); os.IsNotExist(err) {
		return s.SaveStateFile()
	} else if err != nil {
		return err
	} else if err := json.Unmarshal(buf, s); err != nil {
		return err
	}
	return nil
}

// GetStateJSON returns the state data in json format.
func (s *Subscribe) GetStateJSON() (string, error) {
	s.RLock()
	defer s.RUnlock()
	b, err := json.Marshal(s)
	return string(b), err
}

// SaveStateFile writes out the state file.
func (s *Subscribe) SaveStateFile() error {
	if s.stateFile == "" {
		return nil
	}
	s.RLock()
	defer s.RUnlock()
	if buf, err := json.Marshal(s); err != nil {
		return err
	} else if err := ioutil.WriteFile(s.stateFile, buf, 0640); err != nil {
		return err
	}
	return nil
}

/************************
 *   Events Methods   *
 ************************/

// GetEvents returns all the configured events.
func (s *Subscribe) GetEvents() map[string]map[string]string {
	s.RLock()
	defer s.RUnlock()
	return s.Events
}

// GetEvent returns the rules for an event.
// Rules can be used by the user for whatever they want.
func (s *Subscribe) GetEvent(name string) (map[string]string, error) {
	s.RLock()
	defer s.RUnlock()
	if rules, ok := s.Events[name]; ok {
		return rules, nil
	}
	return nil, ErrorEventNotFound
}

// UpdateEvent adds or updates an event.
func (s *Subscribe) UpdateEvent(name string, rules map[string]string) bool {
	s.Lock()
	defer s.Unlock()
	if rules == nil {
		rules = make(map[string]string)
	}
	if _, ok := s.Events[name]; !ok {
		s.Events[name] = rules
		return true
	}
	for ruleName, rule := range rules {
		if rule == "" {
			delete(s.Events[name], ruleName)
		} else {
			s.Events[name][ruleName] = rule
		}
	}
	return false
}

// RemoveEvent obliterates an event and all subsciptions for it.
func (s *Subscribe) RemoveEvent(name string) (removed int) {
	s.Lock()
	delete(s.Events, name)
	s.Unlock()
	for i := range s.Subscribers {
		if _, ok := s.Subscribers[i].Events[name]; ok {
			s.Subscribers[i].Lock()
			delete(s.Subscribers[i].Events, name)
			s.Subscribers[i].Unlock()
			removed++
		}
	}
	return
}

/**************************
 *   Subscriber Methods   *
 **************************/

// CreateSub creates or updates a subscriber.
func (s *Subscribe) CreateSub(contact, api string, admin, ignore bool) *Subscriber {
	for i := range s.Subscribers {
		if contact == s.Subscribers[i].Contact && api == s.Subscribers[i].API {
			s.Subscribers[i].Admin = admin
			s.Subscribers[i].Ignored = ignore
			// Already exists, return it.
			return s.Subscribers[i]
		}
	}

	s.Subscribers = append(s.Subscribers, &Subscriber{
		Contact: contact,
		API:     api,
		Admin:   admin,
		Ignored: ignore,
		Events:  make(map[string]time.Time),
	})
	return s.Subscribers[len(s.Subscribers)-1:][0]
}

// GetSubscriber gets a subscriber based on their contact info.
func (s *Subscribe) GetSubscriber(contact, api string) (*Subscriber, error) {
	sub := &Subscriber{}
	for i := range s.Subscribers {
		if s.Subscribers[i].Contact == contact && s.Subscribers[i].API == api {
			return s.Subscribers[i], nil
		}
	}
	return sub, ErrorSubscriberNotFound
}

// GetAdmins returns a list of subscribed admins.
func (s *Subscribe) GetAdmins() (subs []*Subscriber) {
	for i := range s.Subscribers {
		if s.Subscribers[i].Admin {
			subs = append(subs, s.Subscribers[i])
		}
	}
	return
}

// GetIgnored returns a list of ignored subscribers.
func (s *Subscribe) GetIgnored() (subs []*Subscriber) {
	for i := range s.Subscribers {
		if s.Subscribers[i].Ignored {
			subs = append(subs, s.Subscribers[i])
		}
	}
	return
}

// GetAllSubscribers returns a list of all subscribers.
func (s *Subscribe) GetAllSubscribers() (subs []*Subscriber) {
	for i := range s.Subscribers {
		subs = append(subs, s.Subscribers[i])
	}
	return
}

/****************************
 *   Subscription Methods   *
 ****************************/

// Subscribe adds an event subscription to a subscriber.
func (s *Subscriber) Subscribe(eventName string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Events[eventName]; ok {
		return ErrorEventExists
	}
	s.Events[eventName] = time.Now()
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
	_, ok := s.Events[eventName]
	if !ok {
		return ErrorEventNotFound
	}
	s.Events[eventName] = time.Now().Add(duration)
	return nil
}

// Subscriptions returns a subscriber's event subscriptions.
func (s *Subscriber) Subscriptions() (events map[string]time.Time) {
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
			if event == eventName && evnData.Before(time.Now()) && checkAPI(s.Subscribers[i].API, s.EnableAPIs) {
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
