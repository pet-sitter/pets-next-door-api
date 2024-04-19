package tests

func CreateForEach(setUp, tearDown func()) func(func()) {
	return func(test func()) {
		setUp()
		test()
		tearDown()
	}
}
