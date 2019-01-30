package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFile3 = "/tmp/this_is_a_testfile_for_subtscribe_test3.go.json"

func TestGetEvents(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub, err := GetDB(testFile3)
	a.NotNil(sub.GetEvents(), "the events map must not be nil")
	a.Nil(err, "getting a db must produce no error")
	a.EqualValues(0, len(sub.GetEvents()), "event count must be 0 since none have been added")
	sub.UpdateEvent("event_test", nil)
	a.EqualValues(1, len(sub.GetEvents()), "event count must be 1 since 1 was added")
}

func TestGetEvent(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub, err := GetDB("")
	a.NotNil(sub.GetEvents(), "the events map must not be nil")
	a.Nil(err, "getting a db must produce no error")
	sub.UpdateEvent("event_test", nil)
	evn, err := sub.GetEvent("event_test")
	a.Nil(err, "there must be no error getting the events that was created")
	a.NotNil(evn, "the event rules map must not be nil")
	a.EqualValues(1, len(sub.GetEvents()), "event count must be 1 since 1 was added")
	evn, err = sub.GetEvent("missing_event")
	a.NotNil(err, "the event is missing and must produce an error")
	a.Nil(evn, "the event rules map must be nil when the event is missing")
}

func TestUpdateEvent(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	sub.UpdateEvent("event_test", nil)
	a.NotNil(sub.Events["event_test"], "the event rules map must not be nil")
	a.EqualValues(0, len(sub.Events["event_test"]), "the event rules map must have zero length")

	// Add 1 rule
	sub.UpdateEvent("event_test", map[string]string{"rule_name": "bar"})
	a.EqualValues(1, len(sub.Events["event_test"]), "the event rules map must have length of 1")
	a.EqualValues("bar", sub.Events["event_test"]["rule_name"], "the rule has the wrong value")
	// Update the same rule.
	sub.UpdateEvent("event_test", map[string]string{"rule_name": "bar2"})
	a.EqualValues(1, len(sub.Events["event_test"]), "the event rules map must have length of 1")
	a.EqualValues("bar2", sub.Events["event_test"]["rule_name"], "the rule did not update")
	// Add a enw rule.
	sub.UpdateEvent("event_test", map[string]string{"rule_name2": "some value"})
	a.EqualValues(2, len(sub.Events["event_test"]), "the event rules map must have length of 1")
	a.EqualValues("some value", sub.Events["event_test"]["rule_name2"], "the rule has the wrong value")
	// Delete a rule.
	sub.UpdateEvent("event_test", map[string]string{"rule_name": ""})
	a.EqualValues(1, len(sub.Events["event_test"]), "the event rules map must have length of 1")
	a.EqualValues("some value", sub.Events["event_test"]["rule_name2"], "the second rule has the wrong value")
}

func TestRemoveEvent(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	a.EqualValues(0, sub.RemoveEvent("no_event"), "event had no subscribers and must not produce any deletions")
	// Make two events to remove.
	sub.Events["some_event"] = nil
	sub.Events["some_event2"] = nil
	// Subscribe a user to one of them.
	s := sub.CreateSub("test_contact", "api", true, false)
	a.Nil(s.Subscribe("some_event2"))
	a.EqualValues(1, sub.RemoveEvent("some_event2"), "event had 1 subscriber")
	a.EqualValues(1, len(sub.GetEvents()), "the event must be deleted")
	a.EqualValues(0, sub.RemoveEvent("some_event"), "event had no subscribers and must not produce any deletions")
	a.EqualValues(0, len(sub.GetEvents()), "the event must be deleted")
}
