package k8seventwatcher

import (
	"regexp"
	"strings"
)

type Regexp struct {
	regex *regexp.Regexp

	inverseMatch bool
}

func NewRegexp(pattern string) (*Regexp, error) {
	r := &Regexp{}

	// Is this an inverse match?
	if strings.HasPrefix(pattern, "!") {
		r.inverseMatch = true
		pattern = pattern[1:]
	} else {
		r.inverseMatch = false
	}

	var err error

	r.regex, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	regexString := ""
	err := unmarshal(&regexString)
	if err != nil {
		return err
	}

	newRegex, err := NewRegexp(regexString)
	if err != nil {
		return err
	}

	r.inverseMatch = newRegex.inverseMatch
	r.regex = newRegex.regex

	return nil
}

func (r *Regexp) MarshalYAML() (interface{}, error) {
	return r.String(), nil
}

func (r *Regexp) String() string {
	regexString := r.regex.String()

	if r.inverseMatch {
		regexString = "!" + regexString
	}

	return regexString
}

func (r *Regexp) MatchString(value string) bool {
	matches := r.regex.MatchString(value)

	if r.inverseMatch {
		matches = !matches
	}

	return matches
}
