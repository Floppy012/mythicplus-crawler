package logs

import "github.com/r3labs/diff/v3"

type ConnectedRealmAddedPayload struct {
	ID uint
}

type ConnectedRealmUpdatedPayload struct {
	ID        uint
	Changelog diff.Changelog
}

type ConnectedRealmRemovedPayload struct {
	ID uint
}

type RealmAddedPayload struct {
	ID uint
}

type RealmUpdatedPayload struct {
	ID        uint
	Changelog diff.Changelog
}

type RealmRemovedPayload struct {
	ID uint
}
