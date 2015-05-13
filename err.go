package gs

import (
	"log"
)

func SetErrorHandler(f func(error)) {
	c = f
}

var c = func(e error) {
	if e != nil {
		panic(e)
		log.Fatal(e)
	}
}
