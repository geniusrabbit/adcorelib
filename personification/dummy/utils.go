package dummy

import "github.com/google/uuid"

func isEmptyUUIDPtr(uuid *uuid.UUID) bool {
	if uuid != nil {
		for i := range len(*uuid) {
			if (*uuid)[i] != 0 {
				return false
			}
		}
	}
	return true
}
