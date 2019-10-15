package k8seventwatcher

import "regexp"

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	regexString := ""
	err := unmarshal(&regexString)
	if err != nil {
		return err
	}

	(*r).Regexp, err = regexp.Compile(regexString)
	if err != nil {
		return err
	}

	return nil
}

func (r *Regexp) MarshalYAML() (interface{}, error) {
	return r.Regexp.String(), nil
}
