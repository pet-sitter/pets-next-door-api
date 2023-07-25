package tests

func CreateForEach(setUp func(), tearDown func()) func(func()) {
	return func(test func()) {
		setUp()
		test()
		tearDown()
	}
}
