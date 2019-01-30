# go-subscribe

A (mostly generic) Subscription Library for Go.

This library allows you to subscribe things (like users) to events.
The library holds subscribers in memory. It can be also be configured
to persist subscriptions to disk. The subscription database should scale
to thousands of subscriptions with hundreds of request per second.
If you need more, this library may not work for you; I do not know.

TODO: Add example use.
