package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"strings"
)

func InitialPopulateToLdap(ADElements *[]ADElement, connectLDAP *ldap.Conn)  {

	for _, i := range *ADElements {
		//fmt.Println(i.DN)
		Add := ldap.NewAddRequest(i.DN)
		//	if len(i.attributes) == 0 {
		//		Warning.Println("Dropping", i.DN, "because of null attributes.")
		//		continue
		//	}
		for _, maps := range i.attributes {
			for key, value := range maps {
				//fmt.Println(value)
				if key == "objectClass" {
					for _, OClass := range value.([]string) {
						Add.Attribute(key, strings.Fields(OClass))
					}
					continue
				}
				Add.Attribute(key, value.([]string))




			}
		}
		Info.Println(Add)
		err := connectLDAP.Add(Add)
		fmt.Println(err)


	}
}