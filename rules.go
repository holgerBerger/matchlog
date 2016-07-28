package main

import (
	"regexp"

	termbox "github.com/nsf/termbox-go"
)

// RuleT describes color rules
type RuleT struct {
	re    string            // re to match the line against
	regex *regexp.Regexp    // compiled regexp
	color termbox.Attribute // color to mark a match
}

// DefaultRules creates a slice with some usefull rules
func DefaultRules() *[]*RuleT {
	rules := []*RuleT{
		NewRuleNoerr(".*(?i)error|fail.*", termbox.ColorRed),
		NewRuleNoerr(".*(?i)warn.*", termbox.ColorYellow),
		NewRuleNoerr(".*(?i)ok|success.*", termbox.ColorGreen),
	}
	return &rules
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
