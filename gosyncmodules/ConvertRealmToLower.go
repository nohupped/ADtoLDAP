package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"regexp"
	"strings"
)

//WrapperStruct to embed *ldap.AddRequest and add a custom method to it.
type AddRequest struct {
	*ldap.AddRequest
}

func (a *AddRequest) SetDN(dn string) {
	a.DN = dn
}

//func ConvertRealmToLower(upperrealm []*ldap.AddRequest) {
//
//	r := regexp.MustCompile(`[A-Z]+=`)
//
//	for _, i := range upperrealm {
//		i := &AddRequest{i}
//		input := i.DN
//		i.SetDN(r.ReplaceAllStringFunc(input, func(m string) string {
//			return strings.ToLower(m)
//		}))
//	}
//
//}


func normalizeResult(elements []LDAPElement, ADBasedn, ldapbasedn string)  {
	logger.Debugln("Un-normalized elements from AD:", elements)

	r := regexp.MustCompile(`[A-Z]+=`)


	for i := range elements {
		//elements[i].DN = strings.Replace(strings.ToLower(ADElements[i].DN), ADBaseDN, LDAPBaseDN, -1)
		elements[i].DN = r.ReplaceAllStringFunc(elements[i].DN, func(m string) string {
			return strings.ToLower(m)
		})
		elements[i].DN = strings.Replace(elements[i].DN, ADBasedn, ldapbasedn, -1)
		for i1 := range elements[i].attributes {
			for k, v := range elements[i].attributes[i1] {
				normalizedattribute := ConvertAttributesToLower(&v)
				// replace each element in normalizedattribute with ldapbasedn
				var attributes []string
				for _, attr := range *normalizedattribute {
					a := strings.Replace(attr, ADBasedn, ldapbasedn, -1)
					attributes = append(attributes, a)
				}
				elements[i].attributes[i1][k] = attributes
			}
		}
	}
	logger.Debugln("Normalized AD elements that match the LDAP BaseDN:", elements)

}

func ConvertAttributesToLower(upperAttribute *[]string) *[]string {
	r := regexp.MustCompile(`[A-Z]+=`)
	var attributeAggregated []string
	for _, attribute := range *upperAttribute {
		tmpstring := r.ReplaceAllStringFunc(attribute, func(m string) string {
			return strings.ToLower(m)
		})
		attributeAggregated = append(attributeAggregated, tmpstring)

	}
	return upperAttribute
}
