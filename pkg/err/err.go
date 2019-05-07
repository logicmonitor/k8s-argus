package err

import (
	log "github.com/sirupsen/logrus"
)

func RecoverError(msg string) {
	if err := recover(); err != nil {
		if msg != "" {
			log.Errorf("%s recover error: %s", msg, err)
			msg += ", "
		} else {
			log.Errorf("Recover error: %s", err)
		}
	}
}
