package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFile3 = "/tmp/this_is_a_testfile_for_subtscribe_test3.go.json"

func TestNew(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub, err := GetDB(testFile3)

	assert.NotNil(sub.Events.Names(), "the events slice must not be nil")
	assert.Nil(err, "getting a db must produce no error")
	assert.EqualValues(0, len(sub.Events.Names()), "event count must be 0 since none have been added")
	assert.Nil(sub.Events.New("event_test", nil))
	assert.EqualValues(1, len(sub.Events.Names()), "event count must be 1 since 1 was added")
}

func TestGetEvent(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub, err := GetDB("")

	assert.NotNil(sub.Events.Names(), "the events map must not be nil")
	assert.Nil(err, "getting a db must produce no error")
	assert.Nil(sub.Events.New("event_test", nil))
	assert.True(sub.Events.Exists("event_test"), "this event exists so the method must return true")
	assert.EqualValues(1, len(sub.Events.Names()), "event count must be 1 since 1 was added")
	assert.False(sub.Events.Exists("missing_event"), "this event does not exists")
}

func TestNewEvent(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: &Events{Map: make(map[string]*Rules)}}

	a.Nil(sub.Events.New("event_test", nil))
	a.NotNil(sub.Events.Map["event_test"], "the event rules map must not be nil")
}

func TestRemoveEvent(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: &Events{Map: make(map[string]*Rules)}}

	sub.Events.Remove("no_event")
	// Make two events to remove.
	sub.Events.Map["some_event"] = nil
	sub.Events.Map["some_event2"] = nil

	// Subscribe a user to one of them.
	s := sub.CreateSub("test_contact", "api", true, false) // XXX: count them?
	assert.Nil(s.Subscribe("some_event2"))
	sub.EventRemove("some_event2")
	sub.EventRemove("some_event")
}
