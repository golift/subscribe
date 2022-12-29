package subscribe

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckAPI(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	sub := &Subscribe{Events: new(Events)}
	assert.True(sub.checkAPI("test_string"), "an empty slice must always return true")

	sub.EnableAPIs = []string{"event", "test_string"}
	assert.True(sub.checkAPI("test_string://event"), "test_string is an allowed api prefix")

	sub.EnableAPIs = []string{"event", "any"}
	assert.True(sub.checkAPI("test_string"), "any as a slice value must return true")

	sub.EnableAPIs = []string{"event", "all"}
	assert.True(sub.checkAPI("test_string"), "all as a slice value must return true")

	sub.EnableAPIs = []string{"event", "test_string"}
	assert.True(sub.checkAPI("test_string"), "test_string is an allowed api")

	sub.EnableAPIs = []string{"event", "test_string2"}
	assert.False(sub.checkAPI("test_string"), "test_string is not an allowed api")
}

func TestUnSubscribe(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Add 1 user and 3 subscriptions.
	user := sub.CreateSub("myContacNameTest", "apiValueHere", true, true)
	assert.Nil(user.Subscribe("event_name"))
	assert.Nil(user.Subscribe("event_name2"))
	assert.Nil(user.Subscribe("event_name3"))

	// Make sure we can't add the same event twice.
	assert.EqualValues(ErrEventExists, user.Subscribe("event_name3"), "duplicate event allowed")

	// Remove a subscription.
	user.Events.Remove("event_name3")
	assert.EqualValues(2, len(sub.Subscribers[0].Events.Map), "there must be two subscriptions remaining")

	// Remove another.
	user.Events.Remove("event_name2")
	assert.EqualValues(1, len(sub.Subscribers[0].Events.Map), "there must be one subscription remaining")
	user.Events.Remove("event_name_not_here")
}

func TestPause(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	user := sub.CreateSub("contact", "api", true, false)
	assert.Nil(user.Subscribe("eventName"))

	// Make sure pausing a missing event returns the proper error.
	assert.EqualValues(ErrEventNotFound, user.Events.Pause("fake event", 0))

	// Testing a real unpause.
	assert.Nil(user.Events.Pause("eventName", 0))
	assert.WithinDuration(time.Now(), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)

	// Testing a real pause.
	assert.Nil(user.Events.Pause("eventName", 3600*time.Second))
	assert.WithinDuration(time.Now().Add(3600*time.Second),
		sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)
}

func TestIsPaused(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	user := sub.CreateSub("contact", "api", true, false)

	// Go back and fourth a few times.
	assert.Nil(user.Subscribe("eventName"))
	assert.Nil(user.Events.Pause("eventName", 0))
	assert.False(user.Events.IsPaused("eventName"))
	assert.Nil(user.Events.Pause("eventName", 10*time.Second))
	assert.True(user.Events.IsPaused("eventName"))
	assert.Nil(user.Events.UnPause("eventName"))
	assert.False(user.Events.IsPaused("eventName"))

	// Missing event is always paused.
	assert.True(user.Events.IsPaused("missingEvent"))
}

func TestSubscriptions(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	user := sub.CreateSub("contact", "api", true, false)
	events := []string{"eventName", "eventName1", "eventName3", "eventName5"}

	sort.Strings(events)

	for _, e := range events {
		assert.Nil(user.Subscribe(e))
	}

	assert.Equal(events, user.Events.Names(), "wrong subscriptions provided")
}

func TestGetSubscribers(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subs := sub.GetSubscribers("evn")
	assert.EqualValues(0, len(subs), "there must be no subscribers")

	// Add 1 subscriber and 3 subscriptions.
	user := sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	assert.Nil(user.Subscribe("event_name"))
	assert.Nil(user.Subscribe("event_name2"))
	assert.Nil(user.Subscribe("event_name3"))

	// Add 1 more subscriber and 3 more subscriptions, 2 paused.
	user = sub.CreateSub("myContacNameTest2", "apiValueHere", true, false)
	assert.Nil(user.Subscribe("event_name"))
	assert.Nil(user.Subscribe("event_name2"))
	assert.Nil(user.Subscribe("event_name3"))
	assert.Nil(user.Events.Pause("event_name2", 10*time.Second))
	assert.Nil(user.Events.Pause("event_name3", 10*time.Minute))

	// Add another ignore subscriber with 1 subscription.
	user = sub.CreateSub("myContacNameTest3", "apiValueHere", true, true)
	assert.Nil(user.Subscribe("event_name"))

	// Test that ignore keeps the ignored subscriber out.
	assert.EqualValues(2, len(sub.GetSubscribers("event_name")), "there must be 2 subscribers")

	// Test that resume time keeps a subscriber out.
	assert.EqualValues(1, len(sub.GetSubscribers("event_name2")), "there must be 1 subscriber")
	assert.EqualValues(1, len(sub.GetSubscribers("event_name3")), "there must be 1 subscriber")
}
