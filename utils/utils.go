package utils

func ExitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
