package logs

var LogType = &logType{
	ConnectedRealmAdded:           "CONNECTED_REALM_ADDED",
	ConnectedRealmUpdated:         "CONNECTED_REALM_UPDATED",
	ConnectedRealmRemoved:         "CONNECTED_REALM_REMOVED",
	RealmAdded:                    "REALM_ADDED",
	RealmUpdated:                  "REALM_UPDATED",
	RealmRemoved:                  "REALM_REMOVED",
	MythicPlusAffixAdded:          "MPLUS_AFFIX_ADDED",
	MythicPlusAffixUpdated:        "MPLUS_AFFIX_UPDATED",
	MythicPlusAffixRemoved:        "MPLUS_AFFIX_REMOVED",
	MythicPlusDungeonAdded:        "MPLUS_DUNGEON_ADDED",
	MythicPlusDungeonUpdated:      "MPLUS_DUNGEON_UPDATED",
	MythicPlusDungeonRemoved:      "MPLUS_DUNGEON_REMOVED",
	MythicPlusSeasonAdded:         "MPLUS_SEASON_ADDED",
	MythicPlusSeasonUpdated:       "MPLUS_SEASON_UPDATED",
	MythicPlusSeasonRemoved:       "MPLUS_SEASON_REMOVED",
	MythicPlusPeriodAdded:         "MPLUS_PERIOD_ADDED",
	MythicPlusActiveSeasonChanged: "MPLUS_ACTIVE_SEASON_CHANGED",
}

type logType struct {
	ConnectedRealmAdded           string
	ConnectedRealmUpdated         string
	ConnectedRealmRemoved         string
	RealmAdded                    string
	RealmUpdated                  string
	RealmRemoved                  string
	MythicPlusAffixAdded          string
	MythicPlusAffixUpdated        string
	MythicPlusAffixRemoved        string
	MythicPlusDungeonAdded        string
	MythicPlusDungeonUpdated      string
	MythicPlusDungeonRemoved      string
	MythicPlusSeasonAdded         string
	MythicPlusSeasonUpdated       string
	MythicPlusSeasonRemoved       string
	MythicPlusPeriodAdded         string
	MythicPlusActiveSeasonChanged string
}
