package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"strings"
	"regexp"
)

func ConvertRealmToLower(upperrealm []*ldap.AddRequest)  {
	r := regexp.MustCompile(`[A-Z]+=`)

	for _, i := range upperrealm {
		input := i.GetDN()
		i.SetDN(r.ReplaceAllStringFunc(input, func(m string) string {
			return strings.ToLower(m)
		}))
	}


}
