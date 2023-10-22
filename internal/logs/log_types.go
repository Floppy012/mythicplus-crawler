package logs

var LogType = &logType{
	ConnectedRealmAdded:   "CONNECTED_REALM_ADDED",
	ConnectedRealmUpdated: "CONNECTED_REALM_UPDATED",
	ConnectedRealmRemoved: "CONNECTED_REALM_REMOVED",
	RealmAdded:            "REALM_ADDED",
	RealmUpdated:          "REALM_UPDATED",
	RealmRemoved:          "REALM_REMOVED",
}

type logType struct {
	ConnectedRealmAdded   string
	ConnectedRealmUpdated string
	ConnectedRealmRemoved string
	RealmAdded            string
	RealmUpdated          string
	RealmRemoved          string
}
