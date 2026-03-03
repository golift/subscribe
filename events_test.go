package subscribe

import (
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	stateFile := filepath.Join(t.TempDir(), "events_state.json")
	sub, err := GetDB(stateFile)

	assertions.NotNil(sub.Events.Names(), "the events slice must not be nil")
	require.NoError(t, err, "getting db must produce no error")
	assertions.Empty(sub.Events.Names(), "event count must be 0 since none have been added")
	require.NoError(t, sub.Events.New("event_test", nil))
	assertions.Len(sub.Events.Names(), 1, "event count must be 1 since 1 was added")
}

func TestGetEvent(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub, err := GetDB("")

	assertions.NotNil(sub.Events.Names(), "the events map must not be nil")
	require.NoError(t, err, "getting db must produce no error")
	require.NoError(t, sub.Events.New("event_test", nil))
	assertions.True(sub.Events.Exists("event_test"), "this event exists so the method must return true")
	assertions.Len(sub.Events.Names(), 1, "event count must be 1 since 1 was added")
	assertions.False(sub.Events.Exists("missing_event"), "this event does not exist")
}

func TestNewEvent(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: &Events{Map: make(map[string]*Rules)}}

	require.NoError(t, sub.Events.New("event_test", nil))
	assertions.NotNil(sub.Events.Map["event_test"], "the event rules map must not be nil")
	assertions.NotNil(sub.Events.Map["event_test"].D, "duration map must be initialized")
	assertions.NotNil(sub.Events.Map["event_test"].I, "integer map must be initialized")
	assertions.NotNil(sub.Events.Map["event_test"].S, "string map must be initialized")
	assertions.NotNil(sub.Events.Map["event_test"].T, "time map must be initialized")
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

	assert.False(t, sub.Events.Exists("some_event"), "global event should be removed")
	assert.False(t, sub.Events.Exists("some_event2"), "global event should be removed")
	assert.False(t, subscriber.Events.Exists("some_event2"), "subscription event should be removed")
}

func TestEventsName(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	require.NoError(t, events.New("CaseSensitive", nil))
	require.NoError(t, events.New("other", nil))

	assert.Equal(t, "CaseSensitive", events.Name("CaseSensitive"))
	assert.Equal(t, "CaseSensitive", events.Name("casesensitive"))
	assert.Empty(t, events.Name("missing"))
}

func TestEventsPauseTimeAndLen(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	assert.Equal(t, 0, events.Len())
	assert.Equal(t, time.Time{}, events.PauseTime("missing"))

	require.NoError(t, events.New("pause_test", nil))
	require.NoError(t, events.Pause("pause_test", 2*time.Minute))
	assert.Equal(t, 1, events.Len())
	assert.WithinDuration(t, time.Now().Add(2*time.Minute), events.PauseTime("pause_test"), 2*time.Second)
}

func TestEventsRuleLifecycleAllTypes(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	require.NoError(t, events.New("event", nil))

	when := time.Now().UTC().Round(time.Second)

	events.RuleSetD("event", "d", 3*time.Minute)
	events.RuleSetI("event", "i", 55)
	events.RuleSetS("event", "s", "value")
	events.RuleSetT("event", "t", when)

	d, found := events.RuleGetD("event", "d")
	require.True(t, found)
	assert.Equal(t, 3*time.Minute, d)

	i, found := events.RuleGetI("event", "i")
	require.True(t, found)
	assert.Equal(t, 55, i)

	s, found := events.RuleGetS("event", "s")
	require.True(t, found)
	assert.Equal(t, "value", s)

	gotTime, found := events.RuleGetT("event", "t")
	require.True(t, found)
	assert.Equal(t, when, gotTime)

	events.RuleDelD("event", "d")
	events.RuleDelI("event", "i")
	events.RuleDelS("event", "s")
	events.RuleDelT("event", "t")

	_, found = events.RuleGetD("event", "d")
	assert.False(t, found)
	_, found = events.RuleGetI("event", "i")
	assert.False(t, found)
	_, found = events.RuleGetS("event", "s")
	assert.False(t, found)
	_, found = events.RuleGetT("event", "t")
	assert.False(t, found)
}

func TestEventsRuleDelAll(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	require.NoError(t, events.New("event", nil))

	events.RuleSetD("event", "shared", time.Second)
	events.RuleSetI("event", "shared", 1)
	events.RuleSetS("event", "shared", "s")
	events.RuleSetT("event", "shared", time.Now())
	events.RuleDelAll("event", "shared")

	_, found := events.RuleGetD("event", "shared")
	assert.False(t, found)
	_, found = events.RuleGetI("event", "shared")
	assert.False(t, found)
	_, found = events.RuleGetS("event", "shared")
	assert.False(t, found)
	_, found = events.RuleGetT("event", "shared")
	assert.False(t, found)
}

func TestEventsRuleGetMissingEvent(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	events.Map["bad"] = nil

	_, found := events.RuleGetD("missing", "r")
	assert.False(t, found)
	_, found = events.RuleGetI("missing", "r")
	assert.False(t, found)
	_, found = events.RuleGetS("missing", "r")
	assert.False(t, found)
	_, found = events.RuleGetT("missing", "r")
	assert.False(t, found)

	_, found = events.RuleGetD("bad", "r")
	assert.False(t, found)
	_, found = events.RuleGetI("bad", "r")
	assert.False(t, found)
	_, found = events.RuleGetS("bad", "r")
	assert.False(t, found)
	_, found = events.RuleGetT("bad", "r")
	assert.False(t, found)
}

func TestEventsNewClonesRules(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	external := &Rules{
		D: map[string]time.Duration{"a": time.Second},
		I: map[string]int{"b": 2},
		S: map[string]string{"c": "value"},
		T: map[string]time.Time{"d": time.Now()},
	}
	require.NoError(t, events.New("event", external))

	external.D["a"] = 4 * time.Second
	external.I["b"] = 42
	external.S["c"] = "changed"
	external.T["d"] = time.Now().Add(time.Hour)

	d, _ := events.RuleGetD("event", "a")
	i, _ := events.RuleGetI("event", "b")
	s, _ := events.RuleGetS("event", "c")
	ts, _ := events.RuleGetT("event", "d")

	assert.Equal(t, time.Second, d)
	assert.Equal(t, 2, i)
	assert.Equal(t, "value", s)
	assert.False(t, ts.IsZero())
}

func TestEventsNamesSorted(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	require.NoError(t, events.New("c", nil))
	require.NoError(t, events.New("a", nil))
	require.NoError(t, events.New("b", nil))

	names := events.Names()
	expected := []string{"a", "b", "c"}
	sort.Strings(expected)
	assert.Equal(t, expected, names)
}
