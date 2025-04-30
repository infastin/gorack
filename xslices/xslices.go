package xslices

// Removes the element s[i] from s, returning the modified slice.
// Panics if i > len(s).
// Performs fast removal without preserving order of the elements.
func FastDelete[S ~[]E, E any](s S, i int) S {
	var zero E
	s[i] = s[len(s)-1]
	s[len(s)-1] = zero
	return s[:len(s)-1]
}
