package opts

type RedisOpts struct {
	Host          string
	Port          int
	Address       []string
	Username      string
	Password      string
	Database      int
	MasterName    string
	MinIdleConns  int
	MaxActive     int
	EnableCluster bool
	Timeout       int
}
