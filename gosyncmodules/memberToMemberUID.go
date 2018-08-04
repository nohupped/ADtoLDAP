package gosyncmodules

import (
	"reflect"
	"strings"
)

/*
func memberTomemberUid(member *interface{}, uid *keyvalue)  []string{
	Info.Println("Found member attribute ", *member, "converting it to memberUid")
	matchCN := regexp.MustCompile("CN=")
	memberlist := make([]string, 0)
	for _, members := range (*member).([]string) {
		tmpmember := strings.Split(members, ",")[0]
		memberclean := matchCN.ReplaceAllString(tmpmember, "")
		memberlist = append(memberlist, memberclean)
	}
	Info.Println("retrieved members as ", memberlist)
	//fmt.Println(memberlist, "\n\n\n\n")
	return memberlist

}
*/

// memberTomemberUid function will populate memberUid attribute with the corresponding uid field
// from the entire ldap request slice. Parameters are the member attribute slice which is of
// type interface{}
func memberTomemberUid(member *interface{}, fullmap *[]LDAPElement) []string {
	uids := make([]string, 0)
	for _, members := range (*member).([]string) {
		uid := checkMemberUIDInLDAPElements(members, fullmap)
		if uid != nil {
			uids = append(uids, *uid)
		}

	}
	logger.Debugln("retrieved members as ", uids)
	return uids

}

//
func checkMemberUIDInLDAPElements(members string, fullmap *[]LDAPElement) *string {
	for _, i := range *fullmap {
		if reflect.DeepEqual(i.DN, members) {
			for _, maps := range i.attributes {
				for k, v := range maps {
					if k == "uid" {
						uid := strings.Join(v.([]string), "")
						return &uid
					}

				}
			}
		}
	}
	//fmt.Println(uids)
	return nil
}
