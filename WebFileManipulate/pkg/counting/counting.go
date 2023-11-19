
package counting

import "wordcount/internal/calculation"

// Count performs the character count for a given segment of data.
func Count(data []byte, resultCh chan<- calculation.Calculation, doneCh <-chan struct{}) {
	var calculation calculation.Calculation

	for _, char := range data {
		select {
		case <-doneCh:
			return
		default:
			if char == '\n' {
				calculation.LineCount++
			} else if char == '!' || char == '"' || char == '$' || char == '*' || char == '.' || char == '<' || char == '?' || char == '~' || char == '{' || char == '`' {
				calculation.PunctuationCount++
			} else if char == 'a' || char == 'A' || char == 'e' || char == 'E' || char == 'i' || char == 'I' || char == 'O' || char == 'o' || char == 'u' || char == 'U' {
				calculation.VowelCount++
			} else if char == '\t' || char == ' ' || char == '\n' || char == '\r' {
				calculation.WordCount++
			}
		}
	}

	resultCh <- calculation
}
