package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode"
)

type trieNode struct {
	children     map[rune]*trieNode
	isWordEnd    bool
	originalWord string
}

type trie struct {
	root *trieNode
}

func newtrie() *trie {
	//same as trie:=new(trie); trie.root=new(trieNode{children:make(map[rune]*trieNode)});
	//return trie
	return &trie{root: &trieNode{children: make(map[rune]*trieNode)}}

}

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

func (t *trie) searchWord(word string) (bool, string) {
    word=cleanWord(word)
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			return false, ""
		}
		node = node.children[ch]
	}
	if node.isWordEnd {
		return true, node.originalWord
	}
	return false, ""

}

func cleanWord(word string) string {
	var sb strings.Builder
	for _, ch := range word {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			sb.WriteRune(unicode.ToLower(ch))
		}

	}
	return sb.String()
}

var (
	new_trie = newtrie()
	once     sync.Once
)

func initializetrie() {
	for city, _ := range CityToTimezone {
		new_trie.insertWord(city,city)
	}
}

func GetCityTZ(city string) (string,error) {
	once.Do(initializetrie)
    found, original:= new_trie.searchWord(city)
    if found {
		return CityToTimezone[original]["tz"], nil
	}
    e:= fmt.Sprintf("City '%s' not found!", city)
	return "", errors.New(e)
}
