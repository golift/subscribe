package subscribe

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
		Events: &subEvents{
			Map: make(map[string]subEventInfo),
		},
	})
	return s.Subscribers[len(s.Subscribers)-1]
}

/* Convenience methods to access specific types of subscribers. */

// GetSubscriber gets a subscriber based on their contact info.
func (s *Subscribe) GetSubscriber(contact, api string) (*Subscriber, error) {
	for _, sub := range s.Subscribers {
		if sub.Contact == contact && sub.API == api {
			return sub, nil
		}
	}
	return nil, ErrorSubscriberNotFound
}

// GetAdmins returns a list of subscribed admins.
func (s *Subscribe) GetAdmins() (subs []*Subscriber) {
	for _, sub := range s.Subscribers {
		if sub.Admin {
			subs = append(subs, sub)
		}
	}
	return
}

// GetIgnored returns a list of ignored subscribers.
func (s *Subscribe) GetIgnored() (subs []*Subscriber) {
	for _, sub := range s.Subscribers {
		if sub.Ignored {
			subs = append(subs, sub)
		}
	}
	return
}
