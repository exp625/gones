package plz

import "io"

func Just(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}

func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}
