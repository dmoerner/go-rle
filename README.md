# What This Repo Contains

This repo contains a simple implementation of run-length encoding in Go by
Daniel Moerner and Damilola Israel Oluwole. Over the time we spent writing this
we ran into some unexpected behavior in Go's string conversion, which was the
subject of a blog post here:
https://moerner.com/posts/for-historical-reasons-go-string-int-conversion/.
This blog post is reproduced below.

# Blog Post: "For Historical Reasons": Go String and Int Conversion

Damilola Israel Oluwole and I have started learning Go, and decided to play
around with the language to write a toy implementation of [Run-length
encoding](https://en.wikipedia.org/wiki/Run-length_encoding). While working on
this little project we ran into a corner case with Go type conversion, which
turns out to be documented in the Go spec but neither of us were aware of.

Run-length encoding is a simple lossless text compression algorithm which
encodes each sequence of *n* identical characters *c* into the substring "cn".
For example, "Hello" is encoded as "H1e1l2o1", and "AAAAAAAAAAH" is
encoded as "A10H1". Obviously in the worst case, a string with no adjacent
identical characters, an RLE encoding of a string of length *n* requires length
*2n*. However, in the best case, a string which consists of a single sequence
of a single character, the encoded length scales as the square root of the
original length. Our implementation can be found on Github:
https://github.com/dmoerner/go-rle.

But in this blog post I want to briefly cover something that tripped us up. We
store the result in a rune slice, and then keep track of the last character
written and the sequential count of such counters. We then load a
new character into a buffer, and if it's distinct from the current sequence, we
append the count and then the new character to the result slice, and restart
the count. Here is a naive way to do this:

```go
if buffer != last {
    result = append(result, rune(count)) // critical line
    result = append(result, buffer)
    last = buffer
    count = 1
}
```

However, this does not work, and nor should we expect it to. In Go, a rune is
an integer encoding a Unicode code point. Our `count` variable is not a
representation of a Unicode code point, it's an integer. Converting a count
representing 65 letters to a rune with `rune()` will result in the Unicode code
point `0x41`, the letter 'A', which is not our intention.

Fortunately we realized this quite quickly, and I thought that we could solve
it by first converting the count into a string literal, and then converting
that into a proper rune slice. This can then be appended using a spread:

```go
    result = append(result, rune[](string(count))...)
```

However, this produces the same result! What's appended is the Unicode code

point represented by `count`. What's worse, we were mostly testing with small
test cases. It turns out that the lowest Unicode Characters in the single
digits are all [control
characters](https://en.wikipedia.org/wiki/Unicode_control_characters) like "End
of Text" (U+0003). If you try to print them out following a standard debugging
procedure of littering your code with print statements, nothing is printed and
your debugging is not going very well.

This really stumped us, including with Googling, until we finally came across a
hint: To use `fmt.Sprintf` to convert the integer to a string instead:

```go
    countString := fmt.Sprintf("%d", count)
    result = append(result, []rune(countString)...)
```

This worked, but we didn't understand why. It wasn't even easy to search for an
answer; for example, the Go [builtin](https://pkg.go.dev/builtin) docs do not list a
`string()` function but only the type. Thanks to some users on IRC for noting
that the answer lies in the [Go
spec](https://go.dev/ref/spec#Conversions_to_and_from_a_string_type):

> Finally, for historical reasons, an integer value may be converted to a
> string type. This form of conversion yields a string containing the (possibly
> multi-byte) UTF-8 representation of the Unicode code point with the given
> integer value. Values outside the range of valid Unicode code points are
> converted to "\uFFFD". [...Examples Omitted...] Note: This form of conversion
> may eventually be removed from the language. The go vet tool flags certain
> integer-to-string conversions as potential errors. Library functions such as
> utf8.AppendRune or utf8.EncodeRune should be used instead.

"For historical reasons", `string(int)` behaves the same as `rune(int)`. In
fact, from what I understand (although I'd like to learn more about this),
tools like `string()` are not actually Go functions at all but primitives built
into the language.
 
The first moral of this story is to always read the spec. Although Go has
excellent documentation, it's organized around Go modules. Something low-level
like this is documented in the specification itself.

The second moral of this story is that the language server protocol
[gopls](https://pkg.go.dev/golang.org/x/tools/gopls) is not a complete
replacement for [go vet](https://pkg.go.dev/cmd/vet). The two are
complementary. As the spec notes, `go vet` would have caught our error
immediately and suggested the solution, but neither of thought to run it.    

The third moral of this story is to rethink what a "simple" test case looks
like. We thought that we could debug the problem by focusing on the most simple
test cases like "a". But this meant we were only looking at the Unicode [Start
of Heading](https://www.compart.com/en/unicode/U+0001) character, which looks
like nothing at all! If we had built up a much larger set of test strings that
got out of the Unicode control characters, our debugging would have likely gone
faster.
