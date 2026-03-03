package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSub(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	asert.Len(sub.Subscribers, 1, "there must be one subscriber")
	asert.True(sub.Subscribers[0].Admin, "admin must be true")
	asert.False(sub.Subscribers[0].Ignored, "ignore must be false")

	// Update values for existing contact.
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	asert.Len(sub.Subscribers, 1, "there must still be one subscriber")
	asert.False(sub.Subscribers[0].Admin, "admin must be changed to false")
	asert.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	asert.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	asert.Equal("myContacNameTest", sub.Subscribers[0].Contact, "contact value is incorrect")
	asert.Equal("apiValueHere", sub.Subscribers[0].API, "api value is incorrect")

	// Add another contact.
	sub.CreateSub("myContacName2Test", "apiValueHere", false, true)
	asert.Len(sub.Subscribers, 2, "there must be two subscribers")
	asert.NotNil(sub.Subscribers[1].Events, "events map must not be nil")
}

func TestGetSubscriber(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	_, err := sub.GetSubscriber("im not here", "fake")
	asert.Equal(ErrSubscriberNotFound, err, "must have asert subscriber not found error")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)

	_, err = sub.GetSubscriber("myContacNameTest", "apiValueHere")
	asert.NoError(err, "must not produce an error getting existing subscriber")
}

func TestAdmin(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetAdmins()
	asert.Empty(subs, "there must be zero admin since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetAdmins()
	asert.Len(subs, 1, "there must be one admin")
}

func TestIgnore(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetIgnored()
	asert.Empty(subs, "there must be zero ignored users since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	sub.CreateSub("myContacNameTest1", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetIgnored()
	asert.Len(subs, 1, "there must be one ignored user")
}
