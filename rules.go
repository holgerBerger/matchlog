package main

/*
	rule handling, here is the regular expressions determining color
	of a line to be displayed

	(c) Holger Berger 2016, under GPL
*/

import (
	"regexp"

	termbox "github.com/nsf/termbox-go"
)

// RulesT is list of rules
type RulesT []*RuleT

// DefaultRules creates a slice with some usefull rules
func DefaultRules() RulesT {
	rules := []*RuleT{
		NewRuleNoerr(".*(?i)error|fail|offline|unable|cannot|no.*found.*", termbox.ColorRed),
		NewRuleNoerr(".*(?i)warn.*", termbox.ColorYellow),
		NewRuleNoerr(".*(?i) ok|ok |success.*", termbox.ColorGreen),
	}
	return rules
}

// Match checks all rules and returns first match
func (r *RulesT) Match(buffer []byte) termbox.Attribute {
	for _, rule := range *r {
		t := rule.Match(buffer)
		if t != termbox.ColorDefault {
			return t
		}
	}

	return termbox.ColorDefault

}

// RuleT describes color rules
type RuleT struct {
	re    string            // re to match the line against
	regex *regexp.Regexp    // compiled regexp
	color termbox.Attribute // color to mark a match
}

// NewRule returns a Rule object with the compiled regex
func NewRule(re string, color termbox.Attribute) (*RuleT, error) {
	var newrule RuleT
	var err error
	newrule.re = re
	newrule.color = color
	newrule.regex, err = regexp.Compile(re)
	return &newrule, err
}

// NewRuleNoerr returns a Rule object with the compiled regex
func NewRuleNoerr(re string, color termbox.Attribute) *RuleT {
	var newrule RuleT
	var err error
	newrule.re = re
	newrule.color = color
	newrule.regex, err = regexp.Compile(re)
	if err != nil {
		panic("error in regexp " + re)
	}
	return &newrule
}

// Match matches the rule against a byte buffer
func (r *RuleT) Match(buffer []byte) termbox.Attribute {
	if r.regex.Match(buffer) {
		return r.color
	}
	return termbox.ColorDefault
}
