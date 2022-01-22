package opts

import "time"

type MysqlOpts struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	Charset               string
	ParseTime             bool
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}
