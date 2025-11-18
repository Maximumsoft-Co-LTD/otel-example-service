package domain

import (
	"time"
)

type MongoOptions struct {
	URL    string // mongodb://user:pass@host:27017/?replicaSet=rs0
	DBName string // ชื่อ database หลัก

	// timeout (seconds)
	ConnTimeoutSec int // connect + ping
	OpTimeoutSec   int // default per operation context
	TxnTimeoutSec  int // run transaction context
}

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
	Prefix   string        // optional key prefix (e.g. "office:")
	DialTO   time.Duration // dial timeout
	ReadTO   time.Duration
	WriteTO  time.Duration
}
