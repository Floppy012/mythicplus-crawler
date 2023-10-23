package logs

var LogType = &logType{
	ConnectedRealmAdded:      "CONNECTED_REALM_ADDED",
	ConnectedRealmUpdated:    "CONNECTED_REALM_UPDATED",
	ConnectedRealmRemoved:    "CONNECTED_REALM_REMOVED",
	RealmAdded:               "REALM_ADDED",
	RealmUpdated:             "REALM_UPDATED",
	RealmRemoved:             "REALM_REMOVED",
	MythicPlusAffixAdded:     "MPLUS_AFFIX_ADDED",
	MythicPlusAffixUpdated:   "MPLUS_AFFIX_UPDATED",
	MythicPlusAffixRemoved:   "MPLUS_AFFIX_REMOVED",
	MythicPlusDungeonAdded:   "MPLUS_DUNGEON_ADDED",
	MythicPlusDungeonUpdated: "MPLUS_DUNGEON_UPDATED",
	MythicPlusDungeonRemoved: "MPLUS_DUNGEON_REMOVED",
}

type logType struct {
	ConnectedRealmAdded      string
	ConnectedRealmUpdated    string
	ConnectedRealmRemoved    string
	RealmAdded               string
	RealmUpdated             string
	RealmRemoved             string
	MythicPlusAffixAdded     string
	MythicPlusAffixUpdated   string
	MythicPlusAffixRemoved   string
	MythicPlusDungeonAdded   string
	MythicPlusDungeonUpdated string
	MythicPlusDungeonRemoved string
}
