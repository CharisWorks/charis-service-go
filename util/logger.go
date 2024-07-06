package util

import "log"

func Logger(text string) {
	if PRODUCTION {
		log.Print(text)
	}
}
