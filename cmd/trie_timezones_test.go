package cmd

import (
	"testing"
)

func TestInsertAndSearchWord(t *testing.T) {
	type searchResult struct {
		found    bool
		location []string
	}
	type testCase struct {
		given string
		want  searchResult
	}

	tests := []testCase{
		{given: "London", want: searchResult{true, []string{"London"}}},
		{given: "Berlin", want: searchResult{true, []string{"Berlin"}}},
		{given: "Lo", want: searchResult{true, []string{"London", "Los Angeles"}}},
		// {given: "Los", want: searchResult{true, []string{"London", "Los Angeles"}}},
		{given: "Paris", want: searchResult{false, []string{}}},
	}
	trie := newtrie()

	//insert some cities into the Trie
	trie.insertWord("London", "London")
	trie.insertWord("Los Angeles", "Los Angeles")
	trie.insertWord("Berlin", "Berlin")

	//test matches
	for _, test := range tests {
		found, got := trie.searchWordWithPrefix(test.given)
		if found != test.want.found {
			t.Fatalf("trie.searchWordWithPrefix(%v) = %v, want %v", test.given, found, test.want.found)
		}
		if !EqualSlices(got, test.want.location) {
			t.Fatalf("`trie.searchWordWithPrefix(%v)`= %v, want %v", test.given, got, test.want.location)
		}

	}
}

func EqualSlices[T comparable](p, q []T) bool {
	if len(p) != len(q) {
		return false
	}
	for i := range p {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

func TestLevenshteinDistance(t *testing.T) {
	type testCase struct {
		given [2]string
		want  int
	}

	tests := []testCase{
		{given: [2]string{"kitten", "kitten"}, want: 0},
		{given: [2]string{"kritib", "ksytip"}, want: 3},
		{given: [2]string{"henry", "ryan"}, want: 5},
	}
	for _, test := range tests {
		got := levenshteinDistance(test.given[0], test.given[1])
		if got != test.want {
			t.Fatalf("`levenshteinDistance(%v,%v)=%v`, want %v", test.given[0], test.given[1], got, test.want)
		}

	}
}

func TestCleanWord(t *testing.T) {
	type testCase struct {
		given string
		want  string
	}
	tests := []testCase{
		{given: "!abc2", want: "abc2"},
		{given: "123@abc$def", want: "123abcdef"},
		{given: "Hello, World!", want: "helloworld"},
		{given: "kritib", want: "kritib"},
	}
	for _, test := range tests {
		got := cleanWord(test.given)
		if got != test.want {
			t.Fatalf("`cleanWord(%v)=%v`, want %v", test.given, got, test.want)

		}
	}

}

func TestFindClosestMatches(t *testing.T) {
	words := []string{"test", "testing", "tested", "tent", "kritib"}
	got := findClosestMatches("test", words, 3)
	if len(got) != 3 {
		t.Fatalf("Expected 3 matches, but got %d", len(got))
	}
	if got[0] != "test" {
		t.Fatalf("Expected `test` as the closest match, but got %v", got[0])
	}
}

func TestGetMatchingLocationCity(t *testing.T) {

	type searchResult struct {
		location []string
		err      string
	}
	type testCase struct {
		given string
		want  searchResult
	}

	tests := []testCase{
		{given: "Berlin", want: searchResult{[]string{"Berlin"}, ""}},
		{given: "Lond", want: searchResult{[]string{"London"}, ""}},
		{given: "Xyz", want: searchResult{[]string{}, "City 'Xyz' not found!"}},
	}

	for _, test := range tests {
		gotCities, err := getMatchingLocation(test.given, "")
        	if err != nil {
			// Compare error messages
			if err.Error() != test.want.err {
				t.Fatalf("getMatchingLocation(%v,\"\") returned error '%v', want '%v'", test.given, err.Error(), test.want.err)
			}
		} else if test.want.err != "" {
			t.Fatalf("getMatchingLocation(%v,\"\") returned no error, want error '%v'", test.given, test.want.err)
		}
		if !EqualSlices(gotCities, test.want.location) {
			t.Fatalf("`getMatchingLocation(%v,\"\")=%v`, want %v", test.given, gotCities, test.want.location)
		}
	}
}


func TestGetMatchingLocationCountry(t *testing.T) {

	type searchResult struct {
		location []string
		err      string
	}
	type testCase struct {
		given string
		want  searchResult
	}

	tests := []testCase{
		{given: "US", want: searchResult{[]string{"United States of America"}, ""}},
		{given: "NPL", want: searchResult{[]string{"Nepal"}, ""}},
		{given: "Nep", want: searchResult{[]string{"Nepal"}, ""}},
		{given: "United", want: searchResult{[]string{"United Kingdom","United Arab Emirates", "United States of America","United States Minor Outlying Islands"}, ""}},
		{given: "Xyz", want: searchResult{[]string{}, "Country 'Xyz' not found!"}},
	}

	for _, test := range tests {
		gotCountries, err := getMatchingLocation("",test.given)
        	if err != nil {
			// Compare error messages
			if err.Error() != test.want.err {
				t.Fatalf("getMatchingLocation(\"\",%v) returned error '%v', want '%v'", test.given, err.Error(), test.want.err)
			}
		} else if test.want.err != "" {
			t.Fatalf("getMatchingLocation(\"\", %v) returned no error, want error '%v'", test.given, test.want.err)
		}
		if !EqualSlices(gotCountries, test.want.location) {
			t.Fatalf("`getMatchingLocation(\"\", %v)=%v`, want %v", test.given, gotCountries, test.want.location)

		}

	}

}
