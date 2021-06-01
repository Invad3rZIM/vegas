package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Vegas")
	testMode1 := true
	testMode2 := false

	//testDataFile := "./staticdata/smalltestset.txt"
	realDataFile := "./staticdata/cmudict.0.7a.txt"

	///read file, scan it line by line, parse it into a node structure
	file, err := os.Open(realDataFile)

	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	// and then a loop iterates through
	// and prints each of the slice values.

	wordNodes := map[string]*WordNode{}
	allPhenomeTokens := map[string]int{}

	allConjoinedSyllables := &map[string]int{}
	allConjoinedSyllablesByHash := &map[int]string{}

	for _, each_line := range text {
		isCommentLine := false
		isEmptyLine := false

		if len(strings.Trim(each_line, " ")) == 0 {
			isEmptyLine = true
		}

		if strings.Index(each_line, ";;;") != -1 {
			isCommentLine = true
		}
		if isCommentLine == false && isEmptyLine == false {
			tokens := strings.Split(each_line, " ")

			node := &WordNode{
				Word:            tokens[0],
				RawTokens:       tokens[1:],
				ConjoinedTokens: []*ConjoinedToken{},
			}
			node.GenerateSpeechTokens()
			node.GenerateHashSequence(allConjoinedSyllables)
			node.GenerateSyllableHash()

			wordNodes[tokens[0]] = node

			for i := 1; i < len(tokens); i = i + 1 {
				if _, found := allPhenomeTokens[tokens[i]]; found {
					allPhenomeTokens[tokens[i]] += 1
				} else {
					allPhenomeTokens[tokens[i]] = 1
				}
			}
		}
	}

	allConjoinedSyllablesByHashDeref := *allConjoinedSyllablesByHash
	for key, val := range *allConjoinedSyllables {
		allConjoinedSyllablesByHashDeref[val] = key
	}
	allConjoinedSyllablesByHash = &allConjoinedSyllablesByHashDeref

	fmt.Println("\ntotal syllables : ", len(*allConjoinedSyllables))

	nodeArr := []*WordNode{}

	for _, val := range wordNodes {
		nodeArr = append(nodeArr, val)
	}

	bitTracker := map[int64][]*WordNode{}

	for _, word := range wordNodes {
		if arr, found := bitTracker[word.SyllableHash]; found {
			bitTracker[word.SyllableHash] = append(arr, word)
		} else {
			bitTracker[word.SyllableHash] = []*WordNode{word}
		}
	}

	for bitSequence, wordNodeArr := range bitTracker {
		fmt.Println("bs : " + fmt.Sprint(bitSequence))
		for i := 0; i < len(wordNodeArr); i = i + 1 {
			if i%2 == 0 {
				fmt.Print(ccfg(230, 130, 255))
			} else {
				fmt.Print(ccfg(100, 100, 255))
			}
			fmt.Print((wordNodeArr[i]).Word + " ")
		}
		fmt.Println()
	}

	if testMode2 == true {
		testNode := wordNodes["AFTER"]

		if testNode != nil {
			bitHash := testNode.SyllableHash
			bitHash = 3
			reverseRoot_STN := buildSTN_Reverse(bitTracker[bitHash], allConjoinedSyllablesByHash)
			reverseRoot_STN.Iterate(0)
		}
		for i := 0; i < 5; i = i + 1 {
			fmt.Println()
		}
	}

	if testMode1 == true {
		fmt.Println("RUNNING TESTMODE 1 \n\n\n")
		testWords2 := []string{"CLEVER", "CLEAVER", "LEVEL", "LEVER", "FOREVER"}
		testWordNodes := []*WordNode{}

		for _, word := range testWords2 {
			fmt.Println("ADDING WORD : " + ccfg(200, 200, 100) + word + resetfg)
			if _, found := wordNodes[word]; found == true {
				testWordNodes = append(testWordNodes, wordNodes[word])
			}
		}

		reverseRoot_STN := buildSTN_Reverse(testWordNodes, allConjoinedSyllablesByHash)

		reverseRoot_STN.Iterate(0)
	}
	for i := 0; i < 5; i = i + 1 {
		fmt.Println()
	}

	dictBinding := &DictionaryBinding{
		Words:                        &wordNodes,
		ReverseRoot:                  buildSTN_Reverse(nodeArr, allConjoinedSyllablesByHash),
		ForwardRoot:                  buildSTN_Forward(nodeArr, allConjoinedSyllablesByHash),
		SyllableBitTracker:           &bitTracker,
		AllConjoinedSyllables:        allConjoinedSyllables,
		AllCoinjoinedSyllablesByHash: allConjoinedSyllablesByHash,
	}

	fmt.Println(dictBinding != nil)
	//	fmt.Println(dictBinding)

	getRhymes("BARREL", dictBinding)
}

type DictionaryBinding struct {
	Words                        *map[string]*WordNode
	ReverseRoot                  *STN
	ForwardRoot                  *STN
	SyllableBitTracker           *map[int64][]*WordNode
	AllConjoinedSyllables        *map[string]int
	AllCoinjoinedSyllablesByHash *map[int]string
}

func getRhymes(needsRhyme string, dictionaryBinding *DictionaryBinding) {
	//1. lookup word in dictionary to see if it exists
	wordMapping := *dictionaryBinding.Words

	wordExistsInDictionary := false

	var wordToRhyme *WordNode = nil

	if node, found := wordMapping[needsRhyme]; found {
		wordExistsInDictionary = true
		wordToRhyme = node
	}

	wordToRhyme.Output()
	fmt.Println("word found : ", wordExistsInDictionary)

	if wordToRhyme != nil {
		startingIndex := 0

		for i := 0; i < len(wordToRhyme.HashSequence) && startingIndex == 0; i = i + 1 {
			if wordToRhyme.HashSequence[i] >= 10000 {
				startingIndex = i
			}
		}
		for key, toCheckAgainst := range wordMapping {

			// for i := 0; i < len(wordNode.HashSequence); i = i + 1 {
			// 	hashVal := wordNode.HashSequence[i]

			// 	fmt.Println("hashVal : ", hashVal)

			// 	for j := 0; j < 10
			// }
			if key != needsRhyme {

				foundStartingHook := false
				score := 0
				increments := 0

				for j := 0; j < len(toCheckAgainst.HashSequence) && increments+startingIndex < len(wordToRhyme.HashSequence)-2; j = j + 1 {
					if foundStartingHook == false {
						if toCheckAgainst.HashSequence[j] == wordToRhyme.HashSequence[startingIndex] {
							foundStartingHook = true
							score = score + 1
							increments += 1

							if len(toCheckAgainst.HashSequence) > j+1 && toCheckAgainst.HashSequence[j+1] == wordToRhyme.HashSequence[startingIndex+1] {
								score += 1
							}

							if len(toCheckAgainst.HashSequence) > j+2 && toCheckAgainst.HashSequence[j+2] == wordToRhyme.HashSequence[startingIndex+2] {
								score += 1
							}
						}
					}
				}

				if score > 1 {
					toCheckAgainst.Output()
				}
				// 				i := 0
				// 				j := 0

				// 				//find first vowel sequence in initial word.
				// 				//find first vowel sequence in checkAgainst word
				// 				//if they are the same, add to score
				// 				rhymeScore := 0

				// 				//matching vowels
				// 				//different consonants

				// // 				for i < len(wordToRhyme.HashSequence) && j < len(toCheckAgainst.HashSequence) {
				// // 					seqI := wordToRhyme.HashSequence[i]
				// // 					seqJ := wordToRhyme.HashSequence[j]
				// // ioi
				//	}

			}
		}
	}
}

/*
	TODO - add rhymeschemes
		exactly different consonants,

	add word associations


*/
