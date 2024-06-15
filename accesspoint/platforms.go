package accesspoint

// DSPPlatformAccessor provides methods access to the handlers
type DSPPlatformAccessor interface {
	ListFactories() []Factory
	PlatformByProtocol(protocol, codename string) (Platformer, error)
}
