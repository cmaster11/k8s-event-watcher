package pkg

import "testing"

func TestRegexpInverse(t *testing.T) {
	regexString := "hello"
	regex, err := NewRegexp(regexString)
	if err != nil {
		t.Fatal(err)
	}

	if !regex.MatchString("hello123") {
		t.Fatal("no match")
	}

	regexInverseString := "!hello"
	regexInverse, err := NewRegexp(regexInverseString)
	if err != nil {
		t.Fatal(err)
	}

	if regexInverse.MatchString("hello123") {
		t.Fatal("no match")
	}

	if !regexInverse.MatchString("mello123") {
		t.Fatal("no match")
	}

}
