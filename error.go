package tachikoma

import log "github.com/sirupsen/logrus"

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
