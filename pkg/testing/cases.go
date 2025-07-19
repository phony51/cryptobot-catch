package testing

type TestCase[D, E any] struct {
	Data     D
	Expected E
}
