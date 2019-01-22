# go-subscribe

A (mostly generic) Subscription Library for Go.

This library allows you to subscribe things (like users) to events. It has quite
a few options for defining different APIs (generic term for a thing you can notify).
The library holds subscribers in memory and writes them out to a json file when told
to do so. This file is the subscription database and should scale to thousands of
subscriptions. If you need more, this library may not work for you; I do not know.

TODO: Add an example usage. Add a few missing tests.
