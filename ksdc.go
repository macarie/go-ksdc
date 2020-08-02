package ksdc

type tBigram struct {
	firstRune  rune
	secondRune rune
}

func createBigrams(word string) (map[tBigram]int, int, string) {
	bigrams := make(map[tBigram]int)

	runes := []rune(word)
	runesCount := len(runes) - 1

	for index := 0; index < runesCount; index++ {
		bigrams[tBigram{firstRune: runes[index], secondRune: runes[index+1]}]++
	}

	return bigrams, runesCount, word
}

func calculateScore(
	refBigrams map[tBigram]int,
	refCount int,
	ref string,
) func(
	inputBigrams map[tBigram]int,
	inputCount int,
	input string,
) float64 {
	return func(
		inputBigrams map[tBigram]int,
		inputCount int,
		input string,
	) float64 {
		if len(ref) <= 2 || len(input) <= 2 {
			if ref == input {
				return 1
			}

			return 0
		}

		expendableRefBigrams := make(map[tBigram]int)

		for bigram, count := range refBigrams {
			expendableRefBigrams[bigram] = count
		}

		overlaps := float64(0)

		for bigram, count := range inputBigrams {

			for index := 0; index < count; index++ {
				refBigramCount := expendableRefBigrams[bigram]

				if refBigramCount <= 0 {
					break
				}

				overlaps++

				expendableRefBigrams[bigram]--
			}
		}

		return 2 * overlaps / float64(refCount+inputCount)
	}
}

type Match struct {
	reference string
	score     float64
}

type BestMatch struct {
	reference string
	score     float64
	index     int
}

func FindMatchInReferences(references []string) func(input string) (BestMatch, []Match) {
	refsCount := len(references)
	refsFuncs := make([]func(map[tBigram]int, int, string) float64, refsCount)

	for index := 0; index < refsCount; index++ {
		refsFuncs[index] = calculateScore(createBigrams(references[index]))
	}

	return func(input string) (BestMatch, []Match) {
		inputBigrams, inputRunes, _ := createBigrams(input)

		matches := make([]Match, refsCount)

		bestScore := float64(-1)
		bestMatchIndex := 0

		for index := 0; index < refsCount; index++ {
			score := refsFuncs[index](inputBigrams, inputRunes, input)

			matches[index] = Match{
				score:     score,
				reference: references[index],
			}

			if score > bestScore {
				bestScore = score
				bestMatchIndex = index
			}
		}

		return BestMatch{
			reference: references[bestMatchIndex],
			score:     matches[bestMatchIndex].score,
			index:     bestMatchIndex,
		}, matches
	}
}

func FindMatch(references []string, input string) (BestMatch, []Match) {
	return FindMatchInReferences(references)(input)
}

func CompareStringToReference(reference string) func(input string) float64 {
	finder := FindMatchInReferences([]string{reference})

	return func(input string) float64 {
		bestMatch, _ := finder(input)

		return bestMatch.score
	}
}

func CompareStrings(reference string, input string) float64 {
	return CompareStringToReference(reference)(input)
}
