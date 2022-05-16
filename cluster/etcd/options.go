package etcd

import "time"

type Option func(*Provider)

func WithTTL(ttl time.Duration) Option {
	return func(p *Provider) {
		p.keepAliveTTL = ttl
	}
}

func WithRetryInterval(retryInterval time.Duration) Option {
	return func(p *Provider) {
		p.retryInterval = retryInterval
	}
}

func WithBaseKey(baseKey string) Option {
	return func(p *Provider) {
		p.baseKey = baseKey
	}
}
