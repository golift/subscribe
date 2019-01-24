# go-subscribe

A (mostly generic) Subscription Library for Go.

This library allows you to subscribe things (like users) to events.
The library holds subscribers in memory. It can be also be configured to persist
subscriptions to disk. This file is the subscription database and should scale
to thousands of subscriptions. If you need more, this library may not work for
you; I do not know.

TODO: Add example use. Add a missing test.
