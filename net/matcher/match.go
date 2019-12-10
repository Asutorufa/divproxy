package matcher

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"divproxy/net/cidrmatch"
	"divproxy/net/dns"
	"divproxy/net/domainmatch"
)

type Match struct {
	DNSServer   string
	cidrMatch   *cidrmatch.CidrMatch
	domainMatch *domainmatch.DomainMatcher
}

func (newMatch *Match) InsertOne(str,mark string) error{
	if _, _, err := net.ParseCIDR(str); err == nil {
		if err = newMatch.cidrMatch.InsetOneCIDR(str,mark); err != nil {
			return err
		}
	} else {
		newMatch.domainMatch.Insert(str,mark)
	}
	return nil
}

func NewMatcher(DNSServer string) *Match{
	cidrMatch:= cidrmatch.NewCidrMatch()
	domainMatch := domainmatch.NewDomainMatcher()
	return &Match{
		DNSServer:   DNSServer,
		cidrMatch:   cidrMatch,
		domainMatch: domainMatch,
	}
}

func NewMatcherWithFile(DNSServer string, MatcherFile string) (matcher *Match,err error) {
	cidrMatch:= cidrmatch.NewCidrMatch()
	domainMatch := domainmatch.NewDomainMatcher()
	matcher = &Match{
		DNSServer:   DNSServer,
		cidrMatch:   cidrMatch,
		domainMatch: domainMatch,
	}
	configTemp, err := ioutil.ReadFile(MatcherFile)
	if err != nil{
		return
	}
	for _, s := range strings.Split(string(configTemp), "\n") {
		div := strings.Split(s," ")
		if len(div) < 2{
			return matcher,errors.New("format error: "+s)
		}
		if err := matcher.InsertOne(div[0],div[1]); err != nil {
			log.Println(err)
			continue
		}
	}
	return matcher, nil
}

func (newMatch *Match) MatchStr(str string) (target string, proxy string) {
	isMatch := false
	target = str
	if net.ParseIP(str) != nil {
		isMatch,proxy = newMatch.cidrMatch.MatchOneIP(str)
		log.Println(isMatch,proxy)
	} else {
		isMatch,proxy = newMatch.domainMatch.Search(str)
		if !isMatch {
			if dnsS, isSuccess := dns.DNS(newMatch.DNSServer, str); isSuccess {
				isMatch,proxy = newMatch.cidrMatch.MatchOneIP(dnsS[0])
				target = dnsS[0]
			}
		}
	}
	if isMatch {
		return
	}
	return target, "direct"
}
