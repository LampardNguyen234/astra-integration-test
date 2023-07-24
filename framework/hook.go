package framework

type BeforeFunc func()
type BeforeEachFunc func()
type AfterFunc func()
type AfterEachFunc func()

func Process(f func()) ProcessFunc {
	return f
}

func Before(f func()) BeforeFunc {
	return f
}

func BeforeEach(f func()) BeforeEachFunc {
	return f
}

func After(f func()) AfterFunc {
	return f
}

func AfterEach(f func()) AfterEachFunc {
	return f
}
