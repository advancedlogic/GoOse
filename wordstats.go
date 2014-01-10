package goose

import (
	"github.com/fatih/set"
)

//some word statistics
type wordStats struct {
	//total number of stopwords or good words that we can calculate
	stopWordCount int
	//total number of words on a node
	wordCount int
	//holds an actual list of the stop words we found
	stopWords *set.Set
}

func (this *WordStats) getStopWords() *set.Set {
	return this.stopWords
}

func (this *WordStats) setStopWords(stopWords *set.Set) {
	this.stopWords = stopWords
}

func (this *WordStats) getStopWordCount() int {
	return this.stopWordCount
}

func (this *WordStats) setStopWordCount(stopWordCount int) {
	this.stopWordCount = stopWordCount
}

func (this *WordStats) getWordCount() int {
	return this.wordCount
}

func (this *WordStats) setWordCount(wordCount int) {
	this.wordCount = wordCount
}
