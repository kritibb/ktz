package cmd

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/kritibb/ktz/tzdata"
)

// trieNode represents a node in the Trie data structure.
// Each node contains a map of child nodes, a flag indicating if it is the end of a word,
// and the original word stored at this node.
type trieNode struct {
	children     map[rune]*trieNode
	isWordEnd    bool
	originalWord string
}

// trie represents the Trie data structure.
// It contains a single root node from which all words are stored.
type trie struct {
	root *trieNode
}

// newtrie creates a new Trie and initializes its root node.
// It returns a pointer to the newly created Trie.
func newtrie() *trie {
	//same as trie:=new(trie); trie.root=new(trieNode{children:make(map[rune]*trieNode)});
	//return trie
	return &trie{root: &trieNode{children: make(map[rune]*trieNode)}}

}

// insertWord adds a word to the Trie.
// The 'word' parameter is the word to be added, and 'original' is the original
// word (without removing spaces and non-alphanumeric characters) for reference.
// It cleans the word and inserts each character into the Trie, creating new nodes as necessary.
func (t *trie) insertWord(word, original string) {
	word = cleanWord(word)
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &trieNode{children: make(map[rune]*trieNode)}
		}
		node = node.children[ch]

	}
	node.isWordEnd = true
	node.originalWord = original
}

// searchWordWithPrefix searches the trie for words starting with the given prefix.
// If the prefix exactly matches an original word, it returns that word.
// Otherwise, it returns the two closest matches based on Levenshtein distance.
//
// Parameters:
//   - prefix: The prefix to search for.
//
// Returns:
//   - bool: True if words with the given prefix were found, false otherwise.
//   - []string: A slice of strings containing the closest matching words.
func (t *trie) searchWordWithPrefix(prefix string) (bool, []string) {
	prefix = cleanWord(prefix)
	node := t.root
	for _, ch := range prefix {
		if node.children[ch] == nil {
			return false, []string{}
		}
		node = node.children[ch]
	}
	if node.isWordEnd {
		return true, []string{node.originalWord}
	}
	words := collectAllWords(node, prefix)
	if len(words) == 0 {
		return false, []string{}
	}
	return true, findClosestMatches(prefix, words, 10)
}

// collectAllWords collects all words in the trie starting from the given node and prefix.
// It a slice of strings containing all words.
//
// Parameters:
//   - node: The starting node in the trie.
//   - prefix: The prefix accumulated so far.
//
// Returns:
//   - []string: A slice of strings containing all words found.
func collectAllWords(node *trieNode, prefix string) []string {
	var words []string
	if node.isWordEnd {
		words = append(words, node.originalWord)

	}
	for char, childNode := range node.children {
		childWords := collectAllWords(childNode, prefix+string(char))
		words = append(words, childWords...)
	}

	return words
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
// The Levenshtein distance is a measure of the difference between two sequences.
// It is the minimum number of single-character edits (insertions, deletions or substitutions) required to change one word into the other.
//
// Parameters:
//   - s1: The first string.
//   - s2: The second string.
//
// Returns:
//   - int: The Levenshtein distance between the two strings.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) < len(s2) {
		return levenshteinDistance(s2, s1)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	previousRow := make([]int, len(s2)+1)
	for i := range previousRow {
		previousRow[i] = i
	}
	for i := range s1 {
		currentRow := make([]int, len(s2)+1)
		currentRow[0] = i + 1
		for j := range s2 {
			deletionCost := previousRow[j+1] + 1
			insertionCost := currentRow[j] + 1
			substitutionCost := previousRow[j]
			if s1[i] != s2[j] {
				substitutionCost++
			}
			currentRow[j+1] = min(insertionCost, deletionCost, substitutionCost)
		}
		previousRow = currentRow
	}
	return previousRow[len(s2)]

}

// cleanWord processes the input word to remove any non-letter and non-number characters,
// converting all letters to lowercase. It returns the cleaned word.
func cleanWord(word string) string {
	var sb strings.Builder
	for _, ch := range word {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			sb.WriteRune(unicode.ToLower(ch))
		}

	}
	return sb.String()
}

// findClosestMatches finds the closest matches to the target word from the list of words
// using the Levenshtein distance. It returns up to the specified number of closest matches.
//
// Parameters:
//   - target: The target word to compare against.
//   - words: The list of words to search through.
//   - maxResults: The maximum number of closest matches to return.
//
// Returns:
//   - []string: A slice of the closest matching words.
func findClosestMatches(target string, words []string, maxResults int) []string {
	type wordDistance struct {
		word     string
		distance int
	}

	var distances []wordDistance
	for _, word := range words {
		distances = append(distances, wordDistance{word: word, distance: levenshteinDistance(target, word)})
	}

	// Sort the words by distance
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Return the closest matches up to maxResults
	var closestMatches []string
	for i := 0; i < min(maxResults, int(len(distances))); i++ {
		closestMatches = append(closestMatches, distances[i].word)
	}

	return closestMatches
}

var (
	new_trie = newtrie() //singleton instance of the Trie
	once     sync.Once   //ensures the Trie is initialized only once (only works on online type setting
	// where a program is once started and keeps running until manually exit)
)

// initializetrie initializes the Trie with city names from CityToTimezone.
func initializeCityTrie() {
	for city := range tzdata.CityToIanaTimezone {
		new_trie.insertWord(city, city)
	}
}

// initializetrie initializes the Trie with city names from CityToTimezone.
func initializeCountryTrie() {
	for country := range tzdata.CountryToIanaTimezone {
		new_trie.insertWord(country, country)
	}
}

// getMatchingLocation retrieves matching cities/countries based on a given prefix/city string by performing fuzzy search.
//
// It ensures the Trie is initialized before performing the search.
// It returns matching cities/countires if any city/country matches the given string, otherwise an error.
func getMatchingLocation(city, country string) ([]string, error) {
	if city != "" {
		once.Do(initializeCityTrie)
		if found, matchingOriginalLocationNames := new_trie.searchWordWithPrefix(city); found {
			return matchingOriginalLocationNames, nil
		}
		return nil, fmt.Errorf("City '%s' not found!", city)
	} else {
		//check if the country is 2-letter alpha-2 code
		if countryName, ok := tzdata.Alpha2ToCountry[strings.ToUpper(country)]; ok {
			return []string{countryName}, nil
		} else if countryName, ok := tzdata.Alpha3ToCountry[strings.ToUpper(country)]; ok { //check if the country is 3-letter alpha-3 code
			return []string{countryName}, nil
		} else {
			once.Do(initializeCountryTrie)
			if found, matchingOriginalLocationNames := new_trie.searchWordWithPrefix(country); found {
				return matchingOriginalLocationNames, nil
			}
			return nil, fmt.Errorf("Country '%s' not found!", country)
		}
	}
}
