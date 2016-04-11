package gosyncmodules

import (
	"gopkg.in/ldap.v2"
//	"fmt"
	"gopkg.in/ini.v1"
)

func InitialPopulateToLdap(ADElements *[]ADElement, connectLDAP *ldap.Conn, ReplaceAttributes, MapAttributes *ini.Section)  {

	userObjectClass, err := ReplaceAttributes.GetKey("userObjectClass")
	CheckForError(err)
	groupObjectClass, err := ReplaceAttributes.GetKey("groupObjectClass")
	CheckForError(err)
	mapping := make(map[string] string) //mapping of AD values to ldap values
	for _, i := range MapAttributes.KeyStrings() {
		tmpvar, err := MapAttributes.GetKey(i)
		CheckForError(err)
		mapping[i] = tmpvar.String()
	} //keys = AD attributes, values = ldap values to which it would be mapped

	Info.Println("Creating mappings for the following attributes,", mapping)
	Info.Println("Userobjectclass of AD will be replaced with", userObjectClass)
	Info.Println("Groupobjectclass of AD will be replaced with", groupObjectClass)
	for _, i := range *ADElements {
		//fmt.Println(i.DN)
		Add := ldap.NewAddRequest(i.DN)
		for _, maps := range i.attributes {
			for key, value := range maps {
				//fmt.Println(value)

				if key == "objectClass" {
					if StringInSlice("user", value.([]string)) {
						//Add.Attribute(key, []string{"posixAccount", "top", "inetOrgPerson"})
						Add.Attribute(key, userObjectClass.Strings(","))
						continue
					}
					if StringInSlice("group", value.([]string)) {
						Add.Attribute(key, groupObjectClass.Strings(","))
						continue
					}
				}
				mappingvalue, ok := mapping[key]
				if ok == true {
					if mappingvalue == "memberUid"{
						//members := memberTomemberUid(&value)
						Add.Attribute(mappingvalue, memberTomemberUid(&value))
						continue
					}
					Add.Attribute(mappingvalue, value.([]string))
					continue
				}

				/*if key == "unixHomeDirectory" {
					Add.Attribute("homeDirectory", value.([]string))
					continue
				}*/

				Add.Attribute(key, value.([]string))





			}
		}
		Info.Println(Add)
		err := connectLDAP.Add(Add)
		//fmt.Println(Add)
		Error.Println(err)
		//fmt.Println(err)


	}
}