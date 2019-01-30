package subscribe

/* TODO: a few new methods require tests. */
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSub(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	a.EqualValues(1, len(sub.Subscribers), "there must be one subscriber")
	a.True(sub.Subscribers[0].Admin, "admin must be true")
	a.False(sub.Subscribers[0].Ignored, "ignore must be false")
	// Update values for existing contact.
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	a.EqualValues(1, len(sub.Subscribers), "there must still be one subscriber")
	a.False(sub.Subscribers[0].Admin, "admin must be changed to false")
	a.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	a.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	a.EqualValues(sub.Subscribers[0].Contact, "myContacNameTest", "contact value is incorrect")
	a.EqualValues(sub.Subscribers[0].API, "apiValueHere", "api value is incorrect")
	// Add another contact.
	sub.CreateSub("myContacName2Test", "apiValueHere", false, true)
	a.EqualValues(2, len(sub.Subscribers), "there must be two subscribers")
	a.NotNil(sub.Subscribers[1].Events, "events map must not be nil")
}

func TestGetSubscriber(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	// Test missing subscriber
	_, err := sub.GetSubscriber("im not here", "fake")
	a.EqualValues(ErrorSubscriberNotFound, err, "must have a subscriber not found error")
	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	_, err = sub.GetSubscriber("myContacNameTest", "apiValueHere")
	a.Nil(err, "must not produce an error getting existing subscriber")
}

func TestAdmin(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	// Test missing subscriber
	subs := sub.GetAdmins()
	a.EqualValues(0, len(subs), "there must be zero admin since none were added")
	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)
	subs = sub.GetAdmins()
	a.EqualValues(1, len(subs), "there must be one admin")
}

func TestIgnore(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	// Test missing subscriber
	subs := sub.GetIgnored()
	a.EqualValues(0, len(subs), "there must be zero ignored users since none were added")
	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	sub.CreateSub("myContacNameTest1", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)
	subs = sub.GetIgnored()
	a.EqualValues(1, len(subs), "there must be one ignored user")
}

func TestGetAllSubscribers(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub := &Subscribe{Events: make(Events)}
	// Test missing subscriber
	subs := sub.GetAllSubscribers()
	a.EqualValues(0, len(subs), "there must be zero subs since none were added")
	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	subs = sub.GetAllSubscribers()
	a.EqualValues(1, len(subs), "there must be one sub")
	sub.CreateSub("myContacNameTest2", "apiValueHere2", true, false)
	subs = sub.GetAllSubscribers()
	a.EqualValues(2, len(subs), "there must be two subs")
}
