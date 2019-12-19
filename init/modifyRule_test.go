package divproxyinit

import (
	"testing"
)

func TestAddOneRule(t *testing.T) {
	if err := AddOneRule("ssss.sss url1", "../rule/rule.config"); err != nil {
		t.Log(err)
	}
}

func TestDeleteOneRule(t *testing.T) {
	if err := DeleteOneRule("google.com", "../rule/rule.config"); err != nil {
		t.Log(err)
	}
}
