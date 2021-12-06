package subscribe

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckAPI(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	s := &Subscribe{Events: new(Events)}
	a.True(s.checkAPI("test_string"), "an empty slice must always return true")

	s.EnableAPIs = []string{"event", "test_string"}
	a.True(s.checkAPI("test_string://event"), "test_string is an allowed api prefix")

	s.EnableAPIs = []string{"event", "any"}
	a.True(s.checkAPI("test_string"), "any as a slice value must return true")

	s.EnableAPIs = []string{"event", "all"}
	a.True(s.checkAPI("test_string"), "all as a slice value must return true")

	s.EnableAPIs = []string{"event", "test_string"}
	a.True(s.checkAPI("test_string"), "test_string is an allowed api")

	s.EnableAPIs = []string{"event", "test_string2"}
	a.False(s.checkAPI("test_string"), "test_string is not an allowed api")
}

func TestUnSubscribe(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Add 1 subscriber and 3 subscriptions.
	s := sub.CreateSub("myContacNameTest", "apiValueHere", true, true)
	a.Nil(s.Subscribe("event_name"))
	a.Nil(s.Subscribe("event_name2"))
	a.Nil(s.Subscribe("event_name3"))

	// Make sure we can't add the same event twice.
	a.EqualValues(ErrEventExists, s.Subscribe("event_name3"), "duplicate event allowed")

	// Remove a subscription.
	s.Events.Remove("event_name3")
	a.EqualValues(2, len(sub.Subscribers[0].Events.Map), "there must be two subscriptions remaining")

	// Remove another.
	s.Events.Remove("event_name2")
	a.EqualValues(1, len(sub.Subscribers[0].Events.Map), "there must be one subscription remaining")
	s.Events.Remove("event_name_not_here")
}

func TestPause(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	s := sub.CreateSub("contact", "api", true, false)
	a.Nil(s.Subscribe("eventName"))

	// Make sure pausing a missing event returns the proper error.
	a.EqualValues(ErrEventNotFound, s.Events.Pause("fake event", 0))

	// Testing a real unpause.
	a.Nil(s.Events.Pause("eventName", 0))
	a.WithinDuration(time.Now(), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)

	// Testing a real pause.
	a.Nil(s.Events.Pause("eventName", 3600*time.Second))
	a.WithinDuration(time.Now().Add(3600*time.Second), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)
}

func TestIsPaused(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	s := sub.CreateSub("contact", "api", true, false)

	// Go back and fourth a few times.
	a.Nil(s.Subscribe("eventName"))
	a.Nil(s.Events.Pause("eventName", 0))
	a.False(s.Events.IsPaused("eventName"))
	a.Nil(s.Events.Pause("eventName", 10*time.Second))
	a.True(s.Events.IsPaused("eventName"))
	a.Nil(s.Events.UnPause("eventName"))
	a.False(s.Events.IsPaused("eventName"))

	// Missing event is always paused.
	a.True(s.Events.IsPaused("missingEvent"))
}

func TestSubscriptions(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	s := sub.CreateSub("contact", "api", true, false)
	events := []string{"eventName", "eventName1", "eventName3", "eventName5"}

	sort.Strings(events)

	for _, e := range events {
		a.Nil(s.Subscribe(e))
	}

	a.Equal(events, s.Events.Names(), "wrong subscriptions provided")
}

func TestGetSubscribers(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subs := sub.GetSubscribers("evn")
	a.EqualValues(0, len(subs), "there must be no subscribers")

	// Add 1 subscriber and 3 subscriptions.
	s := sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	a.Nil(s.Subscribe("event_name"))
	a.Nil(s.Subscribe("event_name2"))
	a.Nil(s.Subscribe("event_name3"))

	// Add 1 more subscriber and 3 more subscriptions, 2 paused.
	s = sub.CreateSub("myContacNameTest2", "apiValueHere", true, false)
	a.Nil(s.Subscribe("event_name"))
	a.Nil(s.Subscribe("event_name2"))
	a.Nil(s.Subscribe("event_name3"))
	a.Nil(s.Events.Pause("event_name2", 10*time.Second))
	a.Nil(s.Events.Pause("event_name3", 10*time.Minute))

	// Add another ignore subscriber with 1 subscription.
	s = sub.CreateSub("myContacNameTest3", "apiValueHere", true, true)
	a.Nil(s.Subscribe("event_name"))

	// Test that ignore keeps the ignored subscriber out.
	a.EqualValues(2, len(sub.GetSubscribers("event_name")), "there must be 2 subscribers")

	// Test that resume time keeps a subscriber out.
	a.EqualValues(1, len(sub.GetSubscribers("event_name2")), "there must be 1 subscriber")
	a.EqualValues(1, len(sub.GetSubscribers("event_name3")), "there must be 1 subscriber")
}
