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

	assertions := assert.New(t)

	subscriber := &Subscribe{Events: new(Events)}
	assertions.True(subscriber.checkAPI("test_string"), "an empty slice must always return true")

	subscriber.EnableAPIs = []string{"event", "test_string"}
	assertions.True(subscriber.checkAPI("test_string://event"), "test_string is an allowed api prefix")

	subscriber.EnableAPIs = []string{"event", "any"}
	assertions.True(subscriber.checkAPI("test_string"), "any as an allowed value must return true")

	subscriber.EnableAPIs = []string{"event", "all"}
	assertions.True(subscriber.checkAPI("test_string"), "all as an allowed value must return true")

	subscriber.EnableAPIs = []string{"event", "test_string"}
	assertions.True(subscriber.checkAPI("test_string"), "test_string is an allowed api")

	subscriber.EnableAPIs = []string{"event", "test_string2"}
	assertions.False(subscriber.checkAPI("test_string"), "test_string is not an allowed api")
}

func TestUnSubscribe(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Add 1 subscriber and 3 subscriptions.
	subscriber := sub.CreateSub("myContacNameTest", "apiValueHere", true, true)
	require.NoError(t, subscriber.Subscribe("event_name"))
	require.NoError(t, subscriber.Subscribe("event_name2"))
	require.NoError(t, subscriber.Subscribe("event_name3"))

	// Make sure we can't add the same event twice.
	assertions.Equal(ErrEventExists, subscriber.Subscribe("event_name3"), "duplicate event allowed")

	// Remove a subscription.
	subscriber.Events.Remove("event_name3")
	assertions.Len(sub.Subscribers[0].Events.Map, 2, "there must be two subscriptions remaining")

	// Remove another.
	subscriber.Events.Remove("event_name2")
	assertions.Len(sub.Subscribers[0].Events.Map, 1, "there must be one subscription remaining")
	subscriber.Events.Remove("event_name_not_here")
}

func TestPause(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subscriber := sub.CreateSub("contact", "api", true, false)
	require.NoError(t, subscriber.Subscribe("eventName"))

	// Make sure pausing a missing event returns the proper error.
	assertions.Equal(ErrEventNotFound, subscriber.Events.Pause("fake event", 0))

	// Testing a real unpause.
	require.NoError(t, subscriber.Events.Pause("eventName", 0))
	assertions.WithinDuration(time.Now(), sub.Subscribers[0].Events.Map["eventName"].Pause, 1*time.Second)

	// Testing a real pause.
	require.NoError(t, subscriber.Events.Pause("eventName", 3600*time.Second))
	assertions.WithinDuration(
		time.Now().Add(3600*time.Second),
		sub.Subscribers[0].Events.Map["eventName"].Pause,
		1*time.Second,
	)
}

func TestIsPaused(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	subscriber := sub.CreateSub("contact", "api", true, false)

	// Go back and forth a few times.
	require.NoError(t, subscriber.Subscribe("eventName"))
	require.NoError(t, subscriber.Events.Pause("eventName", 0))
	assertions.False(subscriber.Events.IsPaused("eventName"))
	require.NoError(t, subscriber.Events.Pause("eventName", 10*time.Second))
	assertions.True(subscriber.Events.IsPaused("eventName"))
	require.NoError(t, subscriber.Events.UnPause("eventName"))
	assertions.False(subscriber.Events.IsPaused("eventName"))

	// Missing event is always paused.
	assertions.True(subscriber.Events.IsPaused("missingEvent"))
}

func TestSubscriptions(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}
	subscriber := sub.CreateSub("contact", "api", true, false)
	events := []string{"eventName", "eventName1", "eventName3", "eventName5"}

	sort.Strings(events)

	for _, e := range events {
		require.NoError(t, subscriber.Subscribe(e))
	}

	assertions.Equal(events, subscriber.Events.Names(), "wrong subscriptions provided")
}

func TestGetSubscribers(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	subs := sub.GetSubscribers("evn")
	assertions.Empty(subs, "there must be no subscribers")

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
	subs = sub.GetSubscribers("event_name")
	assertions.Len(subs, 2, "there must be 2 subscribers")
	assertions.ElementsMatch([]string{"myContacNameTest", "myContacNameTest2"}, []string{subs[0].Contact, subs[1].Contact})

	// Test that resume time keeps a subscriber out.
	subs = sub.GetSubscribers("event_name2")
	assertions.Len(subs, 1, "there must be 1 subscriber")
	assertions.Equal("myContacNameTest", subs[0].Contact)

	subs = sub.GetSubscribers("event_name3")
	assertions.Len(subs, 1, "there must be 1 subscriber")
	assertions.Equal("myContacNameTest", subs[0].Contact)
}
