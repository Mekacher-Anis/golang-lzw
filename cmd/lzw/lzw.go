package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func LZWEncodeStr(text string) []uint32 {
	DICT_SIZE := uint32(math.Pow(2, 8))
	dictionary := make(map[string]uint32)
	for i := uint32(0); i < DICT_SIZE; i++ {
		dictionary[string(rune(i))] = i
	}

	res := make([]uint32, 0, len(text))
	var tokenSoFar strings.Builder
	for _, char := range text {
		currToken := tokenSoFar.String() + string(char)
		if _, ok := dictionary[currToken]; ok {
			tokenSoFar.WriteRune(char)
		} else {
			dictionary[currToken] = uint32(len(dictionary))
			if tokenSoFar.Len() != 0 {
				res = append(res, dictionary[tokenSoFar.String()])
			}
			tokenSoFar.Reset()
			tokenSoFar.WriteRune(char)
		}
	}

	if tokenSoFar.Len() != 0 {
		res = append(res, dictionary[tokenSoFar.String()])
	}

	return res
}

func LZWDecodeStr(data []uint32) string {
	if len(data) == 0 {
		return ""
	}
	DICT_SIZE := uint32(math.Pow(2, 8))
	dictionary := make(map[uint32]string)
	for i := uint32(0); i < DICT_SIZE; i++ {
		dictionary[i] = string(rune(i))
	}

	var tokenSoFar strings.Builder
	tokenSoFar.WriteString(dictionary[data[0]])
	data = data[1:]
	var res strings.Builder
	res.WriteString(tokenSoFar.String())
	for _, t := range data {
		currentToken, ok := dictionary[t]
		if !ok {
			currentToken = tokenSoFar.String() + string([]rune(tokenSoFar.String())[0])
		}
		res.WriteString(currentToken)
		dictionary[uint32(len(dictionary))] = tokenSoFar.String() + string([]rune(currentToken)[0])
		tokenSoFar.Reset()
		tokenSoFar.WriteString(currentToken)
	}

	return res.String()
}

func main() {

	data, err := os.ReadFile("tests/t8.shakespeare.txt")

	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}

	originalText := string(data)
	compressed := LZWEncodeStr(originalText)
	decompressed := LZWDecodeStr(compressed)

	if originalText != decompressed {
		fmt.Println("Decompressed text doesn't match original")
		if len(originalText) != len(decompressed) {
			fmt.Printf("Orginal len %v and decompressed %v\n", len(originalText), len(decompressed))
		}
		originalR := []rune(originalText)
		decomporessedR := []rune(decompressed)
		smallerSize := math.Min(float64(len(originalR)-1), float64(len(originalR)-1))
		for i := 0; i < int(smallerSize); i++ {
			if originalR[i] != decomporessedR[i] {
				lowerBound := int(math.Max(0, float64(i-10)))
				upperBound := int(math.Min(smallerSize, float64(i+10)))
				fmt.Printf("Character \"%v\" at position %v doesn't match\n", originalR[i], i)
				fmt.Printf("Original: \"%v\"\n", string(originalR[lowerBound:upperBound]))
				fmt.Printf("Decompressed: \"%v\"\n", string(decomporessedR[lowerBound:upperBound]))
				break
			}
		}
		return
	}

	fmt.Printf("Size of original text : %v\n", len(originalText))
	fmt.Printf("Size of compressed text: %v\n", len(compressed))
	fmt.Printf("Compression Ratio: %v\n", float32(len(originalText))/float32(len(compressed)))
}
