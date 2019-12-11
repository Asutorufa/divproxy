package divproxyinit

import (
	"divproxy/net/matcher"
	"log"
)

func GetMatcher() (match *matcher.Match) {
	match, err := matcher.NewMatcherWithFile("114.114.114.114:53", "./rule/rule.config")
	if err != nil {
		log.Println(err)
		return
	}
	return
}
