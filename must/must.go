// A.k.a "I just wanna write python"
// Good for irrecoverable cases we absolutely want to fail fast with a
// stacktrace.
// Possibly useful for quick prototyping too, but ofc please treat every usage
// of this package as a FIXME item and clean it up before running in
// production.
package must

func Zero(err error) {
	if err != nil {
		panic(err)
	}
}

func One[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func Two[T any, U any](v1 T, v2 U, err error) (T, U) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

func Ok[T any](v T, ok bool) T {
	if !ok {
		panic("not ok")
	}
	return v
}
