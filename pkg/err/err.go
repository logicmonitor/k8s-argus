package err

import (
	log "github.com/sirupsen/logrus"
)

func RecoverError(msg string) {
	if err := recover(); err != nil {
		if msg != "" {
			msg += ", "
		}
		log.Errorf("%srecover error: %s", msg, err)
	}
}
