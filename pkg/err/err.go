package err

import (
	log "github.com/sirupsen/logrus"
)

// RecoverError is a function that recover the panic
func RecoverError(msg string) {
	if err := recover(); err != nil {
		if msg != "" {
			log.Errorf("%s recover error: %v", msg, err)
		} else {
			log.Errorf("Recover error: %v", err)
		}
	}
}
