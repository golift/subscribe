package subscribe

import (
	"errors"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEventsConcurrentAccess(t *testing.T) {
	t.Parallel()

	events := &Events{Map: make(map[string]*Rules)}
	require.NoError(t, events.New("event", nil))

	var waitGroup sync.WaitGroup

	for i := range 32 {
		waitGroup.Go(func() {
			rule := "rule_" + strconv.Itoa(i)
			for j := range 200 {
				events.RuleSetI("event", rule, j)
				events.RuleSetS("event", rule, rule)
				events.RuleSetD("event", rule, time.Duration(j)*time.Millisecond)
				events.RuleSetT("event", rule, time.Now().Add(time.Duration(j)*time.Second))

				err := events.Pause("event", time.Millisecond)
				if err != nil {
					t.Errorf("pause failed: %v", err)

					return
				}

				err = events.UnPause("event")
				if err != nil {
					t.Errorf("unpause failed: %v", err)

					return
				}

				events.IsPaused("event")
				events.PauseTime("event")
				events.RuleGetI("event", rule)
				events.RuleGetS("event", rule)
				events.RuleGetD("event", rule)
				events.RuleGetT("event", rule)
				events.RuleDelAll("event", rule)
			}
		})
	}

	waitGroup.Wait()
}

func TestSubscribeConcurrentAccess(t *testing.T) {
	t.Parallel()

	stateFile := filepath.Join(t.TempDir(), "state.json")
	sub, err := GetDB(stateFile)
	require.NoError(t, err)
	require.NoError(t, sub.Events.New("evt", nil))

	var waitGroup sync.WaitGroup

	for i := range 20 {
		waitGroup.Go(func() {
			contact := "contact_" + strconv.Itoa(i%5)
			runSubscribeOps(t, sub, contact)
		})
	}

	waitGroup.Wait()
}

func runSubscribeOps(t *testing.T, sub *Subscribe, contact string) {
	t.Helper()

	for j := range 100 {
		admin := j%2 == 0
		ignored := j%3 == 0
		subscriber := sub.CreateSub(contact, "api", admin, ignored)

		if subscriber == nil {
			t.Error("subscriber should not be nil")

			return
		}

		err := subscriber.Subscribe("evt")

		if err != nil && !errors.Is(err, ErrEventExists) {
			t.Errorf("subscribe failed: %v", err)

			return
		}

		err = subscriber.Events.Pause("evt", time.Millisecond)
		if err != nil {
			t.Errorf("pause failed: %v", err)

			return
		}

		err = subscriber.Events.UnPause("evt")
		if err != nil {
			t.Errorf("unpause failed: %v", err)

			return
		}

		_, _ = sub.GetSubscriber(contact, "api")
		sub.GetAdmins()
		sub.GetIgnored()
		sub.GetSubscribers("evt")
		_, _ = sub.StateGetJSON()

		err = sub.StateFileSave()
		if err != nil {
			t.Errorf("state save failed: %v", err)

			return
		}
	}
}

// TestSubscriberStateConcurrentAccess races the accessors against the library
// readers. The consuming application flips Admin/Ignored/Contact/Meta from its
// own goroutines while notifications fan out and the state file is written, and
// it holds no Subscribe lock while doing so — writing those fields directly is
// what the accessors exist to avoid.
func TestSubscriberStateConcurrentAccess(t *testing.T) {
	t.Parallel()

	stateFile := filepath.Join(t.TempDir(), "state.json")
	sub, err := GetDB(stateFile)
	require.NoError(t, err)
	require.NoError(t, sub.Events.New("evt", nil))

	subscriber := sub.CreateSubWithID(1, "contact", "api", false, false)
	require.NotNil(t, subscriber)
	require.NoError(t, subscriber.Subscribe("evt"))

	var waitGroup sync.WaitGroup

	for i := range 8 {
		waitGroup.Go(func() { writeSubscriberState(subscriber, i) })
		waitGroup.Go(func() { readSubscriberState(t, sub, subscriber) })
	}

	waitGroup.Wait()

	// The record must survive the hammering intact.
	require.Equal(t, int64(1), subscriber.ID)
	require.NotEmpty(t, subscriber.GetContact())
}

func writeSubscriberState(subscriber *Subscriber, id int) {
	key := "key_" + strconv.Itoa(id)

	for j := range 250 {
		subscriber.SetAdmin(j%2 == 0)
		subscriber.SetIgnored(j%3 == 0)
		subscriber.SetContact("contact_" + strconv.Itoa(j%4))
		subscriber.SetContactIfEmpty("filled")
		subscriber.SetMeta(key, j)
		subscriber.DeleteMeta(key)
	}
}

func readSubscriberState(t *testing.T, sub *Subscribe, subscriber *Subscriber) {
	t.Helper()

	for range 250 {
		subscriber.IsAdmin()
		subscriber.IsIgnored()
		subscriber.GetContact()
		subscriber.GetAllMeta()
		_, _ = subscriber.GetMeta("key_0")

		sub.GetAdmins()
		sub.GetIgnored()
		sub.GetSubscribers("evt")
		_, _ = sub.GetSubscriber("contact_0", "api")
		_, _ = sub.GetSubscriberByID(1, "api")
		_, _ = sub.StateGetJSON()
		sub.CreateSubWithID(1, "contact", "api", false, false)

		err := sub.StateFileSave()
		if err != nil {
			t.Errorf("state save failed: %v", err)

			return
		}
	}
}
