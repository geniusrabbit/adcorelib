//
// @project Geniusrabbit corelib 2016 – 2017, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2022
//

package rand

import "github.com/google/uuid"

// UUID generated
func UUID() string {
	return uuid.New().String()
}
