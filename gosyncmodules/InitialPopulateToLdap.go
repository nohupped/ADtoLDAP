package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"gopkg.in/ini.v1"
//	"strings"
/*	"bytes"
	"encoding/gob"
	"strings"
	"strconv"*/
)

func InitialPopulateToLdap(ADElements *[]ADElement, connectLDAP *ldap.Conn, ReplaceAttributes, MapAttributes, RequiredAttributes *ini.Section)  {

	userObjectClass, err := ReplaceAttributes.GetKey("userObjectClass")
	CheckForError(err)
	groupObjectClass, err := ReplaceAttributes.GetKey("groupObjectClass")
	CheckForError(err)
	mapping := make(map[string] string)
	for _, i := range MapAttributes.KeyStrings() {
		tmpvar, err := MapAttributes.GetKey(i)
		CheckForError(err)
		mapping[i] = tmpvar.String()
	}
	Attributes, err := RequiredAttributes.GetKey("ADattr")
	Info.Println(Attributes)
	CheckForError(err)

	Info.Println("Creating mappings for the following attributes,", mapping)
	Info.Println("Userobjectclass of AD will be replaced with", userObjectClass)
	Info.Println("Groupobjectclass of AD will be replaced with", groupObjectClass)
	for _, i := range *ADElements {
		//fmt.Println(i.DN)
		Add := ldap.NewAddRequest(i.DN)
		//fmt.Println(i.DN)
		//	if len(i.attributes) == 0 {
		//		Warning.Println("Dropping", i.DN, "because of null attributes.")
		//		continue
		//	}
		//var primaryGroupID int
		//var objectSid []byte
		for _, maps := range i.attributes {
			for key, value := range maps {
				//fmt.Println(value)
		/*		if key == "primaryGroupID" {
					tmpprimaryGroupID := strings.Join(value.([]string), "")
					primaryGroupID, _ = strconv.Atoi(tmpprimaryGroupID)
					continue
				}
				if key == "objectSid" {
					var buf bytes.Buffer
					enc := gob.NewEncoder(&buf)
					err := enc.Encode(value)
					fmt.Println(err)
					objectSid = buf.Bytes()
					continue
				}*/
				if key == "objectClass" {
					//Add.Attribute(key, []string{"posixAccount", "top", "inetOrgPerson"})
					Add.Attribute(key, userObjectClass.Strings(","))
					continue
				}
				mappingvalue, ok := mapping[key]
				if ok == true {
					Add.Attribute(mappingvalue, value.([]string))
					continue
				}

				/*if key == "unixHomeDirectory" {
					Add.Attribute("homeDirectory", value.([]string))
					continue
				}
				if key == "primaryGroupID" {
					Add.Attribute("gidNumber", value.([]string))
					continue
				}*/

				Add.Attribute(key, value.([]string))





			}
		}
		Info.Println(Add)
		//GetPrimaryGroup(primaryGroupID, objectSid)
		err := connectLDAP.Add(Add)
		fmt.Println(Add)
		Info.Println(err)
		fmt.Println(err)


	}
}