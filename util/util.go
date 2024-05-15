package util

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/waterdragen/slf-cache/v2/assert"
	"github.com/waterdragen/slf-cache/v2/slf"
)

type Ngram struct {
	chars []rune
	freq  float64
}

type Metric int

const (
	Sfb Metric = iota
	Sft
	Sfr
	Alt
	AltSfs
	Red
	BadRed
	RedSfs
	BadRedSfs
	InOne
	OutOne
	InRoll
	OutRoll
	Unknown
)

const MetricNum = 14

var Table [4096]Metric = LoadTable()

var mapStrMetric = map[string]Metric{
	"sfb":         Sfb,
	"sft":         Sft,
	"alt":         Alt,
	"alt-sfs":     AltSfs,
	"red":         Red,
	"red-sfs":     RedSfs,
	"bad-red":     BadRed,
	"bad-red-sfs": BadRedSfs,
	"inoneh":      InOne,
	"outoneh":     OutOne,
	"inroll":      InRoll,
	"outroll":     OutRoll,
}

var mapStrFinger = map[string]uint16{
	"LP": 0,
	"LR": 1,
	"LM": 2,
	"LI": 3,
	"LT": 4,
	"RT": 5,
	"RI": 6,
	"RM": 7,
	"RR": 8,
	"RP": 9,
}

func strToMetric(str string) Metric {
	return mapStrMetric[str]
}

func strToFinger(str string) uint16 {
	return mapStrFinger[str]
}

func LoadTable() [4096]Metric {
	jsonData, err := os.ReadFile("table.json")
	assert.Ok(err)

	var rawTable map[string]string
	err = json.Unmarshal(jsonData, &rawTable)
	assert.Ok(err)

	table := [4096]Metric{Unknown}
	for fingerStr, metricStr := range rawTable {
		finger0 := strToFinger(fingerStr[0:2])
		finger1 := strToFinger(fingerStr[2:4])
		finger2 := strToFinger(fingerStr[4:6])
		index := finger0<<8 | finger1<<4 | finger2
		table[index] = strToMetric(metricStr)
	}
	return table
}

func LoadCorpus(path string) []Ngram {
	jsonData, err := os.ReadFile(path)
	assert.Ok(err)

	var rawCorpus map[string]float64
	err = json.Unmarshal(jsonData, &rawCorpus)
	assert.Ok(err)

	var corpus []Ngram
	for chars, freq := range rawCorpus {
		corpus = append(corpus, Ngram{[]rune(chars), freq})
	}
	return corpus
}

func getFingerHash(keyMap map[rune]uint16, gram0, gram1, gram2 rune) (uint16, bool) {
	finger0, found := keyMap[gram0]
	if !found {
		return 0, false
	}
	finger1, found := keyMap[gram1]
	if !found {
		return 0, false
	}
	finger2, found := keyMap[gram2]
	if !found {
		return 0, false
	}
	return finger0<<8 | finger1<<4 | finger2, true
}

// AnalyzeTrigrams
// - using cmini implementation
// - ignore spaces
// - case-insensitive
// - sfb is double counted
func AnalyzeTrigrams(layout *slf.Layout, corpus *[]Ngram) []float64 {
	keyMap := make(map[rune]uint16)
	counter := make([]float64, MetricNum)

	for _, key := range layout.Keys {
		s := strings.ToLower(key.Char)
		if s == "" {
			continue
		}
		char := []rune(s)[0]
		finger := uint16(key.Finger)
		keyMap[char] = finger
	}
	for _, trigram := range *corpus {
		gram0 := trigram.chars[0]
		gram1 := trigram.chars[1]
		gram2 := trigram.chars[2]

		if gram0 == ' ' || gram1 == ' ' || gram2 == ' ' {
			continue
		}
		if gram0 == gram1 || gram1 == gram2 {
			counter[Sfr] += trigram.freq
			continue
		}

		fingerHash, ok := getFingerHash(keyMap, gram0, gram1, gram2)
		if !ok {
			counter[Unknown] += trigram.freq
			continue
		}
		gramType := Table[fingerHash]
		counter[gramType] += trigram.freq
	}

	var total float64
	for _, freq := range counter {
		total += freq
	}
	for index, freq := range counter {
		counter[index] = freq / total
	}
	return counter
}
