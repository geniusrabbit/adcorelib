package personification

import "github.com/google/uuid"

func isEmptyUUIDPtr(uuid *uuid.UUID) bool {
	if uuid != nil {
		for i := 0; i < len(*uuid); i++ {
			if (*uuid)[i] != 0 {
				return false
			}
		}
	}
	return true
}
