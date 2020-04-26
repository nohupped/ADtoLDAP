package syncer

import (
	"gopkg.in/ldap.v2"
	//	"fmt"
	"gopkg.in/ini.v1"
)

// InitialPopulateToLdap is a wrapper function on top of other functions in this same package.
func InitialPopulateToLdap(ADElements *[]LDAPElement, connectLDAP *ldap.Conn,
	ReplaceAttributes, MapAttributes *ini.Section, ReturnData bool) []*ldap.AddRequest {
	var ReturnConvertedData []*ldap.AddRequest
	userObjectClass, err := ReplaceAttributes.GetKey("userObjectClass")
	CheckForError(err)
	groupObjectClass, err := ReplaceAttributes.GetKey("groupObjectClass")
	CheckForError(err)
	mapping := make(map[string]string) //mapping of AD values to ldap values
	for _, i := range MapAttributes.KeyStrings() {
		tmpvar, err := MapAttributes.GetKey(i)
		CheckForError(err)
		mapping[i] = tmpvar.String()
	} //keys = AD attributes, values = ldap values to which it would be mapped

	logger.Debugln("Creating mappings for the following attributes,", mapping)
	logger.Debugln("Userobjectclass of AD will be replaced with", userObjectClass)
	logger.Debugln("Groupobjectclass of AD will be replaced with", groupObjectClass)
	for _, i := range *ADElements {
		Add := ldap.NewAddRequest(i.DN)
		for _, maps := range i.attributes {
			for key, value := range maps {

				if key == "objectClass" {
					if StringInSlice("user", value.([]string)) {
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
					if mappingvalue == "memberUid" {
						Add.Attribute(mappingvalue, memberToMemberUID(&value, ADElements))
						continue
					}
					Add.Attribute(mappingvalue, value.([]string))
					continue
				}

				Add.Attribute(key, value.([]string))

			}
		}
		if ReturnData == false {
			err := connectLDAP.Add(Add)
			logger.Errorln(err)
		} else {
			ReturnConvertedData = append(ReturnConvertedData, Add)
		}


	}
	return ReturnConvertedData
}
