package subscribe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckAPI(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.True(checkAPI("test_string", []string{}), "an empty slice must always return true")
	a.True(checkAPI("test_string://event", []string{"event", "test_string"}), "test_string is an allowed api prefix")
	a.True(checkAPI("test_string", []string{"event", "any"}), "any as a slice value must return true")
	a.True(checkAPI("test_string", []string{"event", "all"}), "all as a slice value must return true")
	a.True(checkAPI("test_string", []string{"event", "test_string"}), "test_string is an allowed api")
	a.False(checkAPI("test_string", []string{"event", "test_string2"}), "test_string is not an allowed api")
}

func TestUnSubscribe(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	// Add 1 subscriber and 3 subscriptions.
	s := sub.CreateSub("myContacNameTest", "apiValueHere", true, true)
	a.Nil(s.Subscribe("event_name"))
	a.Nil(s.Subscribe("event_name2"))
	a.Nil(s.Subscribe("event_name3"))
	// Make sure we can't add the same event twice.
	a.EqualValues(ErrorEventExists, s.Subscribe("event_name3"), "duplicate event allowed")
	// Remove a subscription.
	a.Nil(s.UnSubscribe("event_name3"))
	a.EqualValues(2, len(sub.Subscribers[0].Events), "there must be two subscriptions remaining")
	// Remove another.
	a.Nil(s.UnSubscribe("event_name2"))
	a.EqualValues(1, len(sub.Subscribers[0].Events), "there must be one subscription remaining")
	// Make sure we get accurate error when removing a missing event subscription.
	a.EqualValues(ErrorEventNotFound, s.UnSubscribe("event_name_not_here"))
}

func TestPause(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	s := sub.CreateSub("contact", "api", true, false)
	a.Nil(s.Subscribe("eventName"))
	// Make sure pausing a missing event returns the proper error.
	a.EqualValues(ErrorEventNotFound, s.Pause("fake event", 0))
	// Testing a real unpause.
	a.Nil(s.Pause("eventName", 0))
	a.WithinDuration(time.Now(), sub.Subscribers[0].Events["eventName"].Pause, 1*time.Second)
	// Testing a real pause.
	a.Nil(s.Pause("eventName", 3600*time.Second))
	a.WithinDuration(time.Now().Add(3600*time.Second), sub.Subscribers[0].Events["eventName"].Pause, 1*time.Second)
}

func TestGetSubscribers(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
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
	a.Nil(s.Pause("event_name2", 10*time.Second))
	a.Nil(s.Pause("event_name3", 10*time.Minute))
	// Add another ignore subscriber with 1 subscription.
	s = sub.CreateSub("myContacNameTest3", "apiValueHere", true, true)
	a.Nil(s.Subscribe("event_name"))
	// Test that ignore keeps the ignored subscriber out.
	a.EqualValues(2, len(sub.GetSubscribers("event_name")), "there must be 2 subscribers")
	// Test that resume time keeps a subscriber out.
	a.EqualValues(1, len(sub.GetSubscribers("event_name2")), "there must be 1 subscriber")
	a.EqualValues(1, len(sub.GetSubscribers("event_name3")), "there must be 1 subscriber")
}
