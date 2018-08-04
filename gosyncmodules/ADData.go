package gosyncmodules

type keyvalue map[string]interface{}

type LDAPElement struct {
	DN         string
	attributes []keyvalue
}
