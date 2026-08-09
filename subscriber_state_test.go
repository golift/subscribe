package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPI = "api"

func TestSubscriberStateAccessors(t *testing.T) {
	t.Parallel()

	sub := &Subscriber{ID: 1, API: testAPI, Contact: "original"}

	assert.False(t, sub.IsAdmin())
	assert.False(t, sub.IsIgnored())
	assert.Equal(t, "original", sub.GetContact())

	sub.SetAdmin(true)
	sub.SetIgnored(true)
	sub.SetContact("renamed")

	assert.True(t, sub.IsAdmin())
	assert.True(t, sub.IsIgnored())
	assert.Equal(t, "renamed", sub.GetContact())

	sub.SetAdmin(false)
	sub.SetIgnored(false)

	assert.False(t, sub.IsAdmin())
	assert.False(t, sub.IsIgnored())
}

func TestSetContactIfEmpty(t *testing.T) {
	t.Parallel()

	sub := &Subscriber{ID: 1, API: testAPI}

	assert.True(t, sub.SetContactIfEmpty("display name"))
	assert.Equal(t, "display name", sub.GetContact())

	// An existing name wins, including one an admin just chose.
	assert.False(t, sub.SetContactIfEmpty("other name"))
	assert.Equal(t, "display name", sub.GetContact())

	// Nothing to fill with is a no-op, not a wipe.
	assert.False(t, sub.SetContactIfEmpty(""))
	assert.Equal(t, "display name", sub.GetContact())

	var missing *Subscriber

	assert.False(t, missing.SetContactIfEmpty("nobody"))
}

func TestSubscriberMetaAccessors(t *testing.T) {
	t.Parallel()

	sub := &Subscriber{ID: 1, API: testAPI}

	value, found := sub.GetMeta("missing")
	assert.Nil(t, value)
	assert.False(t, found)
	assert.Nil(t, sub.GetAllMeta())

	// SetMeta creates the map when a record arrives without one.
	sub.SetMeta("hasAuth", true)
	sub.SetMeta("count", 7)
	require.NotNil(t, sub.Meta)

	value, found = sub.GetMeta("hasAuth")
	assert.True(t, found)
	assert.Equal(t, true, value)

	all := sub.GetAllMeta()
	assert.Len(t, all, 2)

	// GetAllMeta hands back a copy: editing it must not reach the record.
	all["count"] = 99
	value, _ = sub.GetMeta("count")
	assert.Equal(t, 7, value)

	sub.DeleteMeta("count")
	_, found = sub.GetMeta("count")
	assert.False(t, found)

	// Deleting from a nil map is a no-op, not a panic.
	empty := &Subscriber{}
	empty.DeleteMeta("nope")
	assert.Empty(t, empty.GetAllMeta())
}

// TestSubscriberStateNilReceiver: lookups that find nobody return a nil
// subscriber, so the accessors tolerate one rather than making every caller
// check.
func TestSubscriberStateNilReceiver(t *testing.T) {
	t.Parallel()

	var sub *Subscriber

	assert.False(t, sub.IsAdmin())
	assert.False(t, sub.IsIgnored())
	assert.Empty(t, sub.GetContact())
	assert.Nil(t, sub.GetAllMeta())

	value, found := sub.GetMeta("key")
	assert.Nil(t, value)
	assert.False(t, found)

	// Writes on a nil record are dropped instead of panicking.
	sub.SetAdmin(true)
	sub.SetIgnored(true)
	sub.SetContact("nobody")
	sub.SetMeta("key", "value")
	sub.DeleteMeta("key")
}

// TestSnapshotSubscriberCopiesMeta: the state file snapshot must deep-copy Meta,
// or the marshaler walks a map the application can still be writing.
func TestSnapshotSubscriberCopiesMeta(t *testing.T) {
	t.Parallel()

	sub := &Subscriber{ID: 1, API: testAPI, Contact: "someone"}
	sub.SetMeta("key", "value")

	out := snapshotSubscriber(sub)
	require.NotNil(t, out)
	assert.Equal(t, "someone", out.Contact)
	assert.Equal(t, "value", out.Meta["key"])

	sub.SetMeta("key", "changed")
	assert.Equal(t, "value", out.Meta["key"])

	assert.Nil(t, snapshotSubscriber(nil))
}
