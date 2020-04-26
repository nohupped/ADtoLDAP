package syncer

type keyvalue map[string]interface{}

// LDAPElement holds the DN and attributes of an LDAP/AD entry.
type LDAPElement struct {
	DN         string
	attributes []keyvalue
}
