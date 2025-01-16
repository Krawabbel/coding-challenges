package diff

import (
	"reflect"
	"testing"
)

func TestLcsStr(t *testing.T) {

	type Args struct {
		s1, s2 string
		want   string
	}

	args := []Args{
		{"ABCDEF", "ABCDEF", "ABCDEF"},
		{"ABCDEF", "AXYZ", "A"},
		{"ABC", "XYZ", ""},
		{"AABCXY", "XYZ", "XY"},
		{"", "", ""},
		{"ABCD", "AD", "AD"},
	}

	for i, arg := range args {
		have := lcsstr(arg.s1, arg.s2)
		if have != arg.want {
			t.Errorf("%d: have '%s', want '%s'", i, have, arg.want)
		}
	}
}

func TestLcsLines(t *testing.T) {

	type Args struct {
		s1, s2 []string
		want   []string
	}

	args := []Args{
		{[]string{"This is a test which contains:", "this is the lcs"}, []string{"this is the lcs", "we're testing"}, []string{"this is the lcs"}},
		{
			[]string{"Coding Challenges helps you become a better software engineer through that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"I’ve used or am using these coding challenges as exercise to learn a new programming language or technology.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities."},
			[]string{"Helping you become a better software engineer through coding challenges that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"These are challenges that I’ve used or am using as exercises to learn a new programming language or technology.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities."},
			[]string{"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities."},
		},
	}

	for i, arg := range args {
		have := lcsImpl(arg.s1, arg.s2)
		if !reflect.DeepEqual(arg.want, have) {
			t.Errorf("%d: have '%s', want '%s'", i, have, arg.want)
		}
	}
}

func TestDiff(t *testing.T) {

	type Args struct {
		s1, s2 []string
		want   []string
	}

	args := []Args{
		{[]string{"This is a test which contains:", "this is the lcs"}, []string{"this is the lcs", "we're testing"}, []string{"< This is a test which contains:", "> we're testing"}},
		{
			[]string{
				"Coding Challenges helps you become a better software engineer through that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"I’ve used or am using these coding challenges as exercise to learn a new programming language or technology.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities.",
			},
			[]string{
				"Helping you become a better software engineer through coding challenges that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities.",
				"These are challenges that I’ve used or am using as exercises to learn a new programming language or technology.",
			},
			[]string{
				"< Coding Challenges helps you become a better software engineer through that build real applications.",
				"> Helping you become a better software engineer through coding challenges that build real applications.",
				"< I’ve used or am using these coding challenges as exercise to learn a new programming language or technology.",
				"> These are challenges that I’ve used or am using as exercises to learn a new programming language or technology.",
			},
		},
	}

	for i, arg := range args {
		have := Diff(arg.s1, arg.s2)
		if !reflect.DeepEqual(arg.want, have) {
			t.Errorf("%d: have '%s', want '%s'", i, have, arg.want)
		}
	}
}
