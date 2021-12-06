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
		Events: &Events{
			Map: make(map[string]*Rules),
		},
	})

	return s.Subscribers[len(s.Subscribers)-1]
}

func (s *Subscribe) CreateSubWithID(id int64, contact, api string, admin, ignore bool) *Subscriber {
	if id == 0 {
		return nil
	}

	for i := range s.Subscribers {
		if id == s.Subscribers[i].ID && api == s.Subscribers[i].API {
			s.Subscribers[i].Admin = admin
			s.Subscribers[i].Ignored = ignore
			// Already exists, return it.
			return s.Subscribers[i]
		}
	}

	sub := &Subscriber{
		ID:      id,
		Contact: contact,
		API:     api,
		Admin:   admin,
		Ignored: ignore,
		Events: &Events{
			Map: make(map[string]*Rules),
		},
	}
	s.Subscribers = append(s.Subscribers, sub)

	return sub
}

/* Convenience methods to access specific types of subscribers. */

// GetSubscriber gets a subscriber based on their contact info.
func (s *Subscribe) GetSubscriber(contact, api string) (*Subscriber, error) {
	for _, sub := range s.Subscribers {
		if sub.Contact == contact && sub.API == api {
			return sub, nil
		}
	}

	return nil, ErrSubscriberNotFound
}

// GetSubscriberByID gets a subscriber based on their unique ID.
func (s *Subscribe) GetSubscriberByID(id int64, api string) (*Subscriber, error) {
	if id == 0 {
		return nil, ErrSubscriberNotFound
	}

	for _, sub := range s.Subscribers {
		if sub.ID == id && sub.API == api {
			return sub, nil
		}
	}

	return nil, ErrSubscriberNotFound
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
