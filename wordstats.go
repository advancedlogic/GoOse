/*
This is a golang port of "Goose" originaly licensed to Gravity.com
under one or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.

Golang port was written by Antonio Linari

Gravity.com licenses this file
to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

func (this *wordStats) getStopWords() *set.Set {
	return this.stopWords
}

func (this *wordStats) setStopWords(stopWords *set.Set) {
	this.stopWords = stopWords
}

func (this *wordStats) getStopWordCount() int {
	return this.stopWordCount
}

func (this *wordStats) setStopWordCount(stopWordCount int) {
	this.stopWordCount = stopWordCount
}

func (this *wordStats) getWordCount() int {
	return this.wordCount
}

func (this *wordStats) setWordCount(wordCount int) {
	this.wordCount = wordCount
}
