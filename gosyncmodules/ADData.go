package gosyncmodules

type keyvalue map[string][]string

type LDAPElement struct {
	DN         string
	attributes []keyvalue
}
