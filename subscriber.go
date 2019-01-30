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
		Events:  make(map[string]SubEventInfo),
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
