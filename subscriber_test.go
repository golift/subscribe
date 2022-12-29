package subscribe

/* XXX: a few new methods require tests. */
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSub(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	assert.EqualValues(1, len(sub.Subscribers), "there must be one subscriber")
	assert.True(sub.Subscribers[0].Admin, "admin must be true")
	assert.False(sub.Subscribers[0].Ignored, "ignore must be false")

	// Update values for existing contact.
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	assert.EqualValues(1, len(sub.Subscribers), "there must still be one subscriber")
	assert.False(sub.Subscribers[0].Admin, "admin must be changed to false")
	assert.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	assert.True(sub.Subscribers[0].Ignored, "ignore must be changed to true")
	assert.EqualValues(sub.Subscribers[0].Contact, "myContacNameTest", "contact value is incorrect")
	assert.EqualValues(sub.Subscribers[0].API, "apiValueHere", "api value is incorrect")

	// Add another contact.
	sub.CreateSub("myContacName2Test", "apiValueHere", false, true)
	assert.EqualValues(2, len(sub.Subscribers), "there must be two subscribers")
	assert.NotNil(sub.Subscribers[1].Events, "events map must not be nil")
}

func TestGetSubscriber(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	_, err := sub.GetSubscriber("im not here", "fake")
	assert.EqualValues(ErrSubscriberNotFound, err, "must have a subscriber not found error")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)

	_, err = sub.GetSubscriber("myContacNameTest", "apiValueHere")
	assert.Nil(err, "must not produce an error getting existing subscriber")
}

func TestAdmin(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetAdmins()
	assert.EqualValues(0, len(subs), "there must be zero admin since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetAdmins()
	assert.EqualValues(1, len(subs), "there must be one admin")
}

func TestIgnore(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sub := &Subscribe{Events: new(Events)}

	// Test missing subscriber
	subs := sub.GetIgnored()
	assert.EqualValues(0, len(subs), "there must be zero ignored users since none were added")

	// Test getting real subscriber
	sub.CreateSub("myContacNameTest", "apiValueHere", false, true)
	sub.CreateSub("myContacNameTest1", "apiValueHere", true, false)
	sub.CreateSub("myContacNameTest2", "apiValueHere", false, false)
	sub.CreateSub("myContacNameTest3", "apiValueHere", false, false)

	subs = sub.GetIgnored()
	assert.EqualValues(1, len(subs), "there must be one ignored user")
}
