//
// @project GeniusRabbit corelib 2016 – 2017, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2022
//

package rand

import "github.com/google/uuid"

const (
	shortIDAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// shortIDLen produces ~143 bits of entropy (24 * log2(62) ≈ 142.9)
	// which exceeds UUID v4 (122 bits) while being 12 chars shorter.
	shortIDLen = 24
)

// UUID generated
func UUID() string {
	return uuid.New().String()
}

// ShortID returns a 22-character alphanumeric unique ID.
// Uses the base62 alphabet [0-9A-Za-z] — no hyphens, URL-safe.
func ShortID() string {
	return StrFromDict(shortIDLen, shortIDAlphabet)
}
