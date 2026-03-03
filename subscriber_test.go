package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSub(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	assertions.Len(sub.Subscribers, 1, "there must be one subscriber")
	assertions.True(sub.Subscribers[0].Admin, "admin must be true")
	assertions.False(sub.Subscribers[0].Ignored, "ignore must be false")

	// Update values for existing contact.
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	assertions.Len(sub.Subscribers, 1, "there must still be one subscriber")
	assertions.False(sub.Subscribers[0].Admin, "admin must be changed to false")
	assertions.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	assertions.Equal("myContacNameTest", sub.Subscribers[0].Contact, "contact value is incorrect")
	assertions.Equal("apiValueHere", sub.Subscribers[0].API, "api value is incorrect")

	// Add another contact.
	sub.CreateSub("myContacName2Test", "apiValueHere", false, true)
	assertions.Len(sub.Subscribers, 2, "there must be two subscribers")
	assertions.NotNil(sub.Subscribers[1].Events, "events map must not be nil")
}

func TestGetSubscriber(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	_, err := sub.GetSubscriber("im not here", "fake")
	assertions.Equal(ErrSubscriberNotFound, err, "must have a subscriber not found error")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)

	_, err = sub.GetSubscriber("myContacNameTest", "apiValueHere")
	assertions.NoError(err, "must not produce an error getting existing subscriber")
}

func TestAdmin(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetAdmins()
	assertions.Empty(subs, "there must be zero admin since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetAdmins()
	assertions.Len(subs, 1, "there must be one admin")
	assertions.Equal("myContacNameTest", subs[0].Contact)
}

func TestIgnore(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetIgnored()
	assertions.Empty(subs, "there must be zero ignored users since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	sub.CreateSub("myContacNameTest1", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetIgnored()
	assertions.Len(subs, 1, "there must be one ignored user")
	assertions.Equal("myContacNameTest", subs[0].Contact)
}

func TestCreateSubWithID(t *testing.T) {
	t.Parallel()

	sub := &Subscribe{Events: new(Events)}
	assert.Nil(t, sub.CreateSubWithID(0, "contact", "api", true, false))

	first := sub.CreateSubWithID(10, "contact", "api", true, false)
	require.NotNil(t, first)
	assert.EqualValues(t, 10, first.ID)
	assert.True(t, first.Admin)
	assert.False(t, first.Ignored)
	assert.Len(t, sub.Subscribers, 1)

	second := sub.CreateSubWithID(10, "contact-new", "api", false, true)
	require.NotNil(t, second)
	assert.Same(t, first, second)
	assert.False(t, second.Admin)
	assert.True(t, second.Ignored)
	assert.Equal(t, "contact", second.Contact)
	assert.Len(t, sub.Subscribers, 1)
}

func TestGetSubscriberByID(t *testing.T) {
	t.Parallel()

	sub := &Subscribe{Events: new(Events)}
	_, err := sub.GetSubscriberByID(0, "api")
	assert.Equal(t, ErrSubscriberNotFound, err)

	sub.CreateSubWithID(99, "contact", "api", true, false)
	got, err := sub.GetSubscriberByID(99, "api")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.EqualValues(t, 99, got.ID)

	_, err = sub.GetSubscriberByID(99, "api2")
	assert.Equal(t, ErrSubscriberNotFound, err)
}
