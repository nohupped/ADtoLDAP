package syncer

import (
	"gopkg.in/ldap.v2"
	"regexp"
	"strings"
)

// AddRequest is a wrapperStruct to embed *ldap.AddRequest and add a custom method to it.
type AddRequest struct {
	*ldap.AddRequest
}

// SetDN is used to set the Dn in an AddRequest
func (a *AddRequest) SetDN(dn string) {
	a.DN = dn
}

// ConvertRealmToLower is a normalisation function for Windows Active Directory. This is required because
// the realm is returned capitalised in Active Directory, and needs to be normalised for the sake of LDAP.
func ConvertRealmToLower(upperrealm []*ldap.AddRequest) {

	r := regexp.MustCompile(`[A-Z]+=`)

	for _, i := range upperrealm {
		i := &AddRequest{i}
		input := i.DN
		i.SetDN(r.ReplaceAllStringFunc(input, func(m string) string {
			return strings.ToLower(m)
		}))
	}

}

// ConvertAttributesToLower is a normalisation function. This is for the sake of LDAP
func ConvertAttributesToLower(upperAttribute *[]string) *[]string {
	r := regexp.MustCompile(`[A-Z]+=`)
	var attributeAggregated []string
	for _, attribute := range *upperAttribute {
		tmpstring := r.ReplaceAllStringFunc(attribute, func(m string) string {
			return strings.ToLower(m)
		})
		attributeAggregated = append(attributeAggregated, tmpstring)

	}
	upperAttribute = &attributeAggregated
	return upperAttribute
}
