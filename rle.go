package rle

import (
	"fmt"
	"io"
	"strings"
)

// RunLengthEncoding takes a string and returns a string representing the run
// length encoding of that string. Run length encoding is a text compression
// algorithm which replaces runs of n identical characters c with "cn".
//
// As an alternative to using a rune slice, the Go spec recommends using a
// []byte to store the result and appending to it with [utf8.AppendRune].
func RunLengthEncoding(s string) string {
	if s == "" {
		return s
	}

	stringReader := strings.NewReader(s)

	var result []rune

	var last rune
	var buffer rune
	var err error
	count := 0

	for {
		buffer, _, err = stringReader.ReadRune()
		if err == io.EOF {
			// The following alternatives do not work because of the historical
			// behavior of "string()"
			//
			// result = append(result, rune(string(count)))
			// result = append(result, []rune(string(count))...)

			countString := fmt.Sprintf("%d", count)
			result = append(result, []rune(countString)...)
			break
		}
		if err != nil {
			panic(err)
		}
		if len(result) == 0 {
			result = append(result, buffer)
			count = 1
			last = buffer
			continue
		}
		if buffer != last {
			countString := fmt.Sprintf("%d", count)
			result = append(result, []rune(countString)...)
			result = append(result, buffer)
			last = buffer
			count = 1
		} else {
			count += 1
		}
	}
	return string(result)
}
