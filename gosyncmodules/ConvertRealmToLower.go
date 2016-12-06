package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"strings"
	"regexp"
)
//WrapperStruct to embed *ldap.AddRequest and add a custom method to it.
type AddRequest struct {
	*ldap.AddRequest
}

func (a *AddRequest) SetDN(dn string) {
	a.DN = dn
}

func ConvertRealmToLower(upperrealm []*ldap.AddRequest)  {

	r := regexp.MustCompile(`[A-Z]+=`)

	for _, i := range upperrealm {
		i := &AddRequest{i}
		input := i.DN
		i.SetDN(r.ReplaceAllStringFunc(input, func(m string) string {
			return strings.ToLower(m)
		}))
	}


}

func ConvertAttributesToLower(upperAttribute *[]string)  *[]string {
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
