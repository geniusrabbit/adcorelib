//
// @project GeniusRabbit corelib 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package adtype

// AccessPoint is the DSP source
type AccessPoint interface {
	// ID of the access point
	ID() uint64
}
