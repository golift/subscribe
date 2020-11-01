# go-subscribe

A (mostly generic) Subscription Library for Go.

This library allows you to subscribe things (like users) to events.
The library holds subscribers in memory. It can be also be configured
to persist subscriptions to disk. The subscription database should scale
to thousands of subscriptions with hundreds of request per second.
If you need more, this library may not work for you; I do not know.

Thread safe.

The following is a very simple example. This library provides many other methods
not shown here to deal with event and notification rules, pause/un-pausing, etc.

```golang
package main

import (
	"fmt"
	"time"

	"golift.io/subscribe"
)

func main() {
	// Instantiate an in-memory database. Passing a file-path here loads
	// a DB from disk, or saves a new database to the file-path provided.
	db, _ := subscribe.GetDB("")

	// Create two new subscribers. Initial state has no subscriptions.
	newSub := db.CreateSub("you@email.com.tw", "smtp", false, false)
	newSub2 := db.CreateSub("+18089117234", "sms", false, false)

	// Subscribe the users to an event. Errors only if subscription already exists.
	_ = newSub.Subscribe("party invites")
	_ = newSub2.Subscribe("party invites")

	// Limit subscriber search to only email recipients. Initial state returns any API.
	db.EnableAPIs = []string{"smtp"}

	// Now that your party invites event has subscribers, you can find them when a party invite arrives.
	subs := db.GetSubscribers("party invites")
	fmt.Printf("Sending email to %d subscriber(s):\n", len(subs))

	for _, sub := range subs {
		fmt.Println(sub.Contact)
		// send email.
	}

	// Limit subscriber search to only sms.
	db.EnableAPIs[0] = "sms"
	// Now that your party invites event has subscribers, you can find them when a party invite arrives.
	subs = db.GetSubscribers("party invites")
	fmt.Printf("Sending Text Msg to %d subscriber(s):\n", len(subs))

	for _, sub := range subs {
		fmt.Println(sub.Contact)
		// send sms.
	}

	// if you want to save the DB:
	err := db.StateFileRelocate("/var/lib/somewhere/for/a/file.json")
	if err != nil {
		fmt.Println("Unable to relocate DB:", err)
		return
	}

	// save the DB once in a while, or after making changes.
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		err = db.StateFileSave()
		// This always returns nil if state file path is empty: ""
		if err != nil {
			fmt.Println("Unable to save DB:", err)
		}
	}
}
```

Output:

```
Sending email to 1 subscriber(s):
you@email.com.tw
Sending Text Msg to 1 subscriber(s):
+18089117234
Unable to relocate DB: open /var/lib/somewhere/for/a/file.json: no such file or directory
```

Feedback, ideas and contributions welcomed!
