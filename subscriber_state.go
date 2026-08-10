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

// SetContactIfEmpty fills a blank contact and reports whether it did. Useful
// when a chat provider offers a display name that should never clobber a name
// the subscriber or an admin already chose. Get-then-Set cannot do this: a
// rename landing between the two calls would be overwritten.
func (s *Subscriber) SetContactIfEmpty(contact string) bool {
	if s == nil || contact == "" {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Contact != "" {
		return false
	}

	s.Contact = contact

	return true
}

// Values stored in Meta must be treated as immutable. The library copies the
// Meta map — for a caller, and again for the state file — but a copy of a map
// of `any` is one level deep, so a slice, map or pointer stored in it stays
// shared with every copy. Mutating one of those in place sidesteps the lock and
// can race the JSON marshaling in StateFileSave. Replace the value with SetMeta
// instead of editing what is already in there. A deep copy is not an option
// here: Meta holds arbitrary caller types, which the library cannot walk.

// GetMeta returns one value from the subscriber's Meta map, and whether it was
// present. Treat the value as immutable; see the note above.
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
// it does not exist yet. The value is stored as given, not copied, so hand over
// something that will not be mutated afterwards.
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
// directly races anything that writes it; range over this instead. Adding to or
// deleting from the returned map is safe and does not touch the record, but the
// values are shared with it — treat them as immutable, as above.
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
