package db

import "time"

type Option func(p *PostgreSQL)

func LoadConnAttempts(attempts int) Option {
	return func(p *PostgreSQL) {
		p.connAttempts = attempts
	}
}

func LoadMaxPoolSize(size int) Option {
	return func(p *PostgreSQL) {
		p.maxPoolSize = size
	}
}

func LoadConnTimeout(timeout time.Duration) Option {
	return func(p *PostgreSQL) {
		p.connTimeout = timeout
	}
}
