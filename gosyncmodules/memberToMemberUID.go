package gosyncmodules

import (
	"strings"
	"regexp"
	"fmt"
)

func memberTomemberUid(member *interface{})  []string{
	Info.Println("Found member attribute ", member, "converting it to memberUid")
	matchCN := regexp.MustCompile("CN=")
	memberlist := make([]string, 0)
	for _, members := range (*member).([]string) {
		tmpmember := strings.Split(members, ",")[0]
		memberclean := matchCN.ReplaceAllString(tmpmember, "")
		memberlist = append(memberlist, memberclean)
	}
	Info.Println("retrieved members as ", memberlist)
	fmt.Println(memberlist, "\n\n\n\n")
	return memberlist

}
