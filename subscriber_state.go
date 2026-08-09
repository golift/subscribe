package subscribe

import "maps"

/********************************
 *   Subscriber State Accessors  *
 ********************************/

// The accessors below guard the Subscriber fields that change after the record
// is created: Meta, Contact, Admin and Ignored. A consuming application usually
// flips them from a different goroutine than the one that reads them — an admin
// renaming or ignoring somebody while a notification fans out, say — so the
// library and its consumers share one lock per record.
//
// Every accessor tolerates a nil receiver, which lets callers skip a nil check
// on lookups that may not have found anybody.

// IsAdmin reports whether the subscriber is flagged as an admin.
func (s *Subscriber) IsAdmin() bool {
	if s == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Admin
}

// SetAdmin grants or revokes the admin flag.
func (s *Subscriber) SetAdmin(admin bool) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Admin = admin
}

// IsIgnored reports whether the subscriber is excluded from GetSubscribers().
func (s *Subscriber) IsIgnored() bool {
	if s == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Ignored
}

// SetIgnored excludes the subscriber from GetSubscribers(), or stops doing so.
func (s *Subscriber) SetIgnored(ignored bool) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Ignored = ignored
}

// GetContact returns the contact info used to reach the subscriber.
func (s *Subscriber) GetContact() string {
	if s == nil {
		return ""
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Contact
}

// SetContact replaces the contact info used to reach the subscriber.
func (s *Subscriber) SetContact(contact string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Contact = contact
}

// GetMeta returns one value from the subscriber's Meta map, and whether it was
// present.
func (s *Subscriber) GetMeta(key string) (any, bool) {
	if s == nil {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.Meta[key]

	return value, ok
}

// SetMeta stores one value in the subscriber's Meta map, creating the map when
// it does not exist yet.
func (s *Subscriber) SetMeta(key string, value any) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Meta == nil {
		s.Meta = make(map[string]any)
	}

	s.Meta[key] = value
}

// DeleteMeta drops one value from the subscriber's Meta map.
func (s *Subscriber) DeleteMeta(key string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Meta, key)
}

// GetAllMeta returns a copy of the subscriber's Meta map. Ranging over Meta
// directly races anything that writes it; range over this instead.
func (s *Subscriber) GetAllMeta() map[string]any {
	if s == nil {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Meta == nil {
		return nil
	}

	out := make(map[string]any, len(s.Meta))
	maps.Copy(out, s.Meta)

	return out
}
