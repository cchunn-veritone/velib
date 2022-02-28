package velib

import "log"

func l(text ...interface{}) {
	if v.verbose {
		log.Println("[VELIB] ", text)
	}
}
