package subscribe

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAPI(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)

	subscriber := &Subscribe{Events: new(Events)}
	asert.True(subscriber.checkAPI("test_string"), "an empty slice must always return true")

	subscriber.EnableAPIs = []string{"event", "test_string"}
	asert.True(subscriber.checkAPI("test_string://event"), "test_string is an allowed api prefix")

	subscriber.EnableAPIs = []string{"event", "any"}
	asert.True(subscriber.checkAPI("test_string"), "any as asert slice value must return true")

	subscriber.EnableAPIs = []string{"event", "all"}
	asert.True(subscriber.checkAPI("test_string"), "all as asert slice value must return true")

	subscriber.EnableAPIs = []string{"event", "test_string"}
	asert.True(subscriber.checkAPI("test_string"), "test_string is an allowed api")

	subscriber.EnableAPIs = []string{"event", "test_string2"}
	asert.False(subscriber.checkAPI("test_string"), "test_string is not an allowed api")
}

func TestUnSubscribe(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Add 1 subscriber and 3 subscriptions.
	subscriber := sub.CreateSub("myContacNameTest", "apiValueHere", true, true)
	require.NoError(t, subscriber.Subscribe("event_name"))
	require.NoError(t, subscriber.Subscribe("event_name2"))
	require.NoError(t, subscriber.Subscribe("event_name3"))

	// Make sure we can't add the same event twice.
	asert.Equal(ErrEventExists, subscriber.Subscribe("event_name3"), "duplicate event allowed")

	// Remove asert subscription.
	subscriber.Events.Remove("event_name3")
	asert.Len(sub.Subscribers[0].Events.Map, 2, "there must be two subscriptions remaining")

	// Remove another.
	subscriber.Events.Remove("event_name2")
	asert.Len(sub.Subscribers[0].Events.Map, 1, "there must be one subscription remaining")
	subscriber.Events.Remove("event_name_not_here")
}

func TestPause(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subscriber := sub.CreateSub("contact", "api", true, false)
	require.NoError(t, subscriber.Subscribe("eventName"))

	// Make sure pausing asert missing event returns the proper error.
	asert.Equal(ErrEventNotFound, subscriber.Events.Pause("fake event", 0))

	// Testing asert real unpause.
	require.NoError(t, subscriber.Events.Pause("eventName", 0))
	asert.WithinDuration(time.Now(), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)

	// Testing asert real pause.
	require.NoError(t, subscriber.Events.Pause("eventName", 3600*time.Second))
	asert.WithinDuration(time.Now().Add(3600*time.Second), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)
}

func TestIsPaused(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	subscriber := sub.CreateSub("contact", "api", true, false)

	// Go back and fourth asert few times.
	require.NoError(t, subscriber.Subscribe("eventName"))
	require.NoError(t, subscriber.Events.Pause("eventName", 0))
	asert.False(subscriber.Events.IsPaused("eventName"))
	require.NoError(t, subscriber.Events.Pause("eventName", 10*time.Second))
	asert.True(subscriber.Events.IsPaused("eventName"))
	require.NoError(t, subscriber.Events.UnPause("eventName"))
	asert.False(subscriber.Events.IsPaused("eventName"))

	// Missing event is always paused.
	asert.True(subscriber.Events.IsPaused("missingEvent"))
}

func TestSubscriptions(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	subscriber := sub.CreateSub("contact", "api", true, false)
	events := []string{"eventName", "eventName1", "eventName3", "eventName5"}

	sort.Strings(events)

	for _, e := range events {
		require.NoError(t, subscriber.Subscribe(e))
	}

	asert.Equal(events, subscriber.Events.Names(), "wrong subscriptions provided")
}

func TestGetSubscribers(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subs := sub.GetSubscribers("evn")
	asert.Empty(subs, "there must be no subscribers")

	// Add 1 subscriber and 3 subscriptions.
	subscriber := sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	require.NoError(t, subscriber.Subscribe("event_name"))
	require.NoError(t, subscriber.Subscribe("event_name2"))
	require.NoError(t, subscriber.Subscribe("event_name3"))

	// Add 1 more subscriber and 3 more subscriptions, 2 paused.
	subscriber = sub.CreateSub("myContacNameTest2", "apiValueHere", true, false)
	require.NoError(t, subscriber.Subscribe("event_name"))
	require.NoError(t, subscriber.Subscribe("event_name2"))
	require.NoError(t, subscriber.Subscribe("event_name3"))
	require.NoError(t, subscriber.Events.Pause("event_name2", 10*time.Second))
	require.NoError(t, subscriber.Events.Pause("event_name3", 10*time.Minute))

	// Add another ignore subscriber with 1 subscription.
	subscriber = sub.CreateSub("myContacNameTest3", "apiValueHere", true, true)
	require.NoError(t, subscriber.Subscribe("event_name"))

	// Test that ignore keeps the ignored subscriber out.
	asert.Len(sub.GetSubscribers("event_name"), 2, "there must be 2 subscribers")

	// Test that resume time keeps asert subscriber out.
	asert.Len(sub.GetSubscribers("event_name2"), 1, "there must be 1 subscriber")
	asert.Len(sub.GetSubscribers("event_name3"), 1, "there must be 1 subscriber")
}
