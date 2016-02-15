package SyncModules


type keyvalue map[string]interface{}

type ADElement struct {
	DN string
	attributes []keyvalue
	}

