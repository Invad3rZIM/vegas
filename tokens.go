package main

import "fmt"

type WordNode struct {
	Word            string
	RawTokens       []string
	ConjoinedTokens []*ConjoinedToken
	HashSequence    []int
	SyllableHash    int64
	VowelCount      int
}

//must assign ids using reference map to track global progess
func (wn *WordNode) GenerateHashSequence(syllableMapReference *map[string]int) {
	syllableMap := *syllableMapReference

	wn.HashSequence = []int{}

	for i := 0; i < len(wn.ConjoinedTokens); i = i + 1 {
		ct := wn.ConjoinedTokens[i]
		syllableKey := ""
		syllableHash := 0

		if ct.IsVowel {
			syllableHash = syllableHash + 10000
			syllableKey = ct.Subtokens[0]
		} else {
			for j := 0; j < len(ct.Subtokens); j = j + 1 {
				syllableKey += ct.Subtokens[j]

				if j < len(ct.Subtokens)-1 {
					syllableKey += "."
				}
			}
		}
		if hash, found := syllableMap[syllableKey]; !found {
			syllableMap[syllableKey] = syllableHash + len(syllableMap)
		} else {
			syllableHash = hash
		}

		wn.HashSequence = append(wn.HashSequence, syllableHash)
	}

	syllableMapReference = &syllableMap
}

func (wn *WordNode) GenerateSyllableHash() {
	var sequence int64 = 0

	for i := 0; i < len(wn.ConjoinedTokens); i = i + 1 {
		token := wn.ConjoinedTokens[i]

		if token.IsVowel {
			sequence = sequence<<1 | 1
		} else {
			sequence = sequence << 1
		}
	}

	wn.SyllableHash = sequence
}

func (wn *WordNode) Output() {
	fmt.Print("   ")
	fmt.Print(resetfg + wn.Word + " ")
	for _, ct := range wn.ConjoinedTokens {
		fmt.Print(ct.ToString())
	}

	fmt.Print(ccfg(255, 100, 100) + " [V = " + fmt.Sprint(wn.VowelCount) + "] ")

	sh := wn.SyllableHash

	fmt.Print(ccbg(100, 255, 255) + ccfg(0, 0, 0))

	bitSequence := ""
	for range wn.ConjoinedTokens {
		i := sh & 1
		sh = sh >> 1
		bitSequence += fmt.Sprint(i)
	}

	reverseBitSequence := ""

	for i := len(bitSequence) - 1; i >= 0; i = i - 1 {
		reverseBitSequence += string(bitSequence[i])
	}

	fmt.Print(reverseBitSequence)
	fmt.Println(resetbg + resetfg + "\n")

}

type ConjoinedToken struct {
	Subtokens []string
	IsVowel   bool
}

func (ct *ConjoinedToken) GenerateConjoinedHash() {

}

func (ct *ConjoinedToken) ToString() string {
	returnStr := ""

	if ct.IsVowel {
		returnStr += ccfg(255, 100, 150) + fmt.Sprint(ct.Subtokens) + resetfg
	} else {

		returnStr += ccfg(255, 255, 100) + fmt.Sprint(ct.Subtokens) + resetfg
	}

	return returnStr
}

func (wn *WordNode) GenerateSpeechTokens() {
	//given

	vowelCount := 0
	currentSubtokens := []string{}
	currentIsVowel := false

	for i := 0; i < len(wn.RawTokens); i = i + 1 {
		token := wn.RawTokens[i]

		if len(token) > 0 {
			if len(token) < 3 {
				currentSubtokens = append(currentSubtokens, token)
				currentIsVowel = false
			} else {
				if len(currentSubtokens) > 0 {
					wn.ConjoinedTokens = append(wn.ConjoinedTokens, &ConjoinedToken{
						Subtokens: currentSubtokens,
						IsVowel:   currentIsVowel,
					})
				}
				if len(token) > 0 {
					wn.ConjoinedTokens = append(wn.ConjoinedTokens, &ConjoinedToken{
						Subtokens: []string{token},
						IsVowel:   true,
					})
					vowelCount = vowelCount + 1
				}
				currentSubtokens = []string{}
			}
		}

	}
	if len(currentSubtokens) > 0 {
		wn.ConjoinedTokens = append(wn.ConjoinedTokens, &ConjoinedToken{
			Subtokens: currentSubtokens,
			IsVowel:   currentIsVowel,
		})
	}

	wn.VowelCount = vowelCount
}

////STN Stands for SYLLABLE TRACKING NODE. It is a recursive tree that is used to organize phenome
///tags. there is forward + reverse functionallity in building
type STN struct {
	Children    map[int]*STN
	SymbolicVal string
	IsLeaf      bool
	IsRoot      bool
	Leaf        *WordNode
}

func buildSTN_Reverse(wordNodes []*WordNode, allConjoinedSyllablesByHash *map[int]string) *STN {
	///forward iteration
	//capture both forward order and reverse order of syllables - might have to do a pure
	//reversal of all consonant nodes later
	reverseRoot_STN := &STN{
		IsRoot:      true,
		Children:    map[int]*STN{},
		SymbolicVal: "-(rR)-",
	}
	///reverse iteration
	for _, node := range wordNodes {
		hashSequence := node.HashSequence

		currentSTN_Pointer := reverseRoot_STN //copy a pointer to the head of the tree

		for i := len(hashSequence) - 1; i >= 0; i = i - 1 {
			hashKey := hashSequence[i]

			if pointer, found := currentSTN_Pointer.Children[hashKey]; found {
				currentSTN_Pointer = pointer
			} else {
				newNode := &STN{
					SymbolicVal: (*allConjoinedSyllablesByHash)[hashKey],
					IsLeaf:      false,
					IsRoot:      false,
					Children:    map[int]*STN{},
				}

				currentSTN_Pointer.Children[hashKey] = newNode
				currentSTN_Pointer = newNode
			}
		}

		currentSTN_Pointer.Leaf = node
		currentSTN_Pointer.IsLeaf = true

	}

	return reverseRoot_STN
}

func buildSTN_Forward(wordNodes []*WordNode, allConjoinedSyllablesByHash *map[int]string) *STN {
	///forward iteration
	//capture both forward order and reverse order of syllables - might have to do a pure
	//reversal of all consonant nodes later
	forwardRoot_STN := &STN{
		IsRoot:      true,
		Children:    map[int]*STN{},
		SymbolicVal: "-(rF)-",
	}

	for _, node := range wordNodes {
		hashSequence := node.HashSequence

		currentSTN_Pointer := forwardRoot_STN //copy a pointer to the head of the tree

		for i := 0; i < len(hashSequence); i = i + 1 {
			hashKey := hashSequence[i]

			if pointer, found := currentSTN_Pointer.Children[hashKey]; found {
				currentSTN_Pointer = pointer
			} else {
				newNode := &STN{
					SymbolicVal: (*allConjoinedSyllablesByHash)[hashKey],
					IsLeaf:      false,
					IsRoot:      false,
					Children:    map[int]*STN{},
				}

				currentSTN_Pointer.Children[hashKey] = newNode
				currentSTN_Pointer = newNode
			}
		}
		currentSTN_Pointer.Leaf = node
		currentSTN_Pointer.IsLeaf = true
	}

	return forwardRoot_STN
}

func (stn *STN) Iterate(iteration int) {
	writeBuffer := ""

	for i := 0; i < iteration; i = i + 1 {
		writeBuffer += "    "
	}

	fmt.Println(writeBuffer + stn.SymbolicVal)

	if stn.Children != nil {
		for _, child := range stn.Children {
			child.Iterate(iteration + 1)
		}
	}

	if stn.Leaf != nil {
		fmt.Print(writeBuffer)
		stn.Leaf.Output()
	}

}
