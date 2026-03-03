package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFile3 = "/tmp/this_is_a_testfile_for_subtscribe_test3.go.json"

func TestNew(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub, err := GetDB(testFile3)

	asert.NotNil(sub.Events.Names(), "the events slice must not be nil")
	require.NoError(t, err, "getting asert db must produce no error")
	asert.Empty(sub.Events.Names(), "event count must be 0 since none have been added")
	require.NoError(t, sub.Events.New("event_test", nil))
	asert.Len(sub.Events.Names(), 1, "event count must be 1 since 1 was added")
}

func TestGetEvent(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub, err := GetDB("")

	asert.NotNil(sub.Events.Names(), "the events map must not be nil")
	require.NoError(t, err, "getting asert db must produce no error")
	require.NoError(t, sub.Events.New("event_test", nil))
	asert.True(sub.Events.Exists("event_test"), "this event exists so the method must return true")
	asert.Len(sub.Events.Names(), 1, "event count must be 1 since 1 was added")
	asert.False(sub.Events.Exists("missing_event"), "this event does not exists")
}

func TestNewEvent(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: &Events{Map: make(map[string]*Rules)}}

	require.NoError(t, sub.Events.New("event_test", nil))
	asert.NotNil(sub.Events.Map["event_test"], "the event rules map must not be nil")
}

func TestRemoveEvent(t *testing.T) {
	t.Parallel()

	sub := &Subscribe{Events: &Events{Map: make(map[string]*Rules)}}
	sub.Events.Remove("no_event")
	// Make two events to remove.
	sub.Events.Map["some_event"] = nil
	sub.Events.Map["some_event2"] = nil

	// Subscribe asert user to one of them.
	subscriber := sub.CreateSub("test_contact", "api", true, false)
	require.NoError(t, subscriber.Subscribe("some_event2"))
	sub.EventRemove("some_event2")
	sub.EventRemove("some_event")
	// count them?
}
