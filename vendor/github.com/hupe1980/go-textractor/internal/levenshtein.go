package internal

// ComputeLevenshteinDistance calculates the Levenshtein distance between two strings.
// It measures the minimum number of single-character edits (insertions, deletions,
// or substitutions) required to change one string into the other.
func ComputeLevenshteinDistance(s1, s2 string) int {
	// If one of the strings is empty, return the length of the other string.
	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	// If the strings are equal, the distance is zero.
	if s1 == s2 {
		return 0
	}

	// Swap to save memory (O(min(a,b)) instead of O(a)).
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	lenS1 := len(s1)
	lenS2 := len(s2)

	x := make([]uint16, lenS1+1)
	for i := 1; i <= lenS1; i++ {
		x[i] = uint16(i)
	}

	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)

		for j := 1; j <= lenS1; j++ {
			current := x[j-1]
			if s2[i-1] != s1[j-1] {
				current = minUint16(minUint16(x[j-1]+1, prev+1), x[j]+1)
			}

			x[j-1] = prev
			prev = current
		}

		x[lenS1] = prev
	}

	// Return the Levenshtein distance.
	return int(x[lenS1])
}

// minUint16 returns the minimum value among three uint16 integers.
func minUint16(a, b uint16) uint16 {
	if a < b {
		return a
	}

	return b
}
