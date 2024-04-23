package router

type routerConfig struct {
	serversAddresses map[string]string
	serverLimits     map[string]uint64
}
