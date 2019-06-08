package github

import (
	"fmt"
	"time"

	"github.com/jeffail/tunny"
)

type executable func() error

type connPool struct {
	p  *tunny.Pool
	to time.Duration
}

// TODO: retry policy?
// TODO: some way to monitor queue sizes
func newPool(opts DatabaseOpts) *connPool {
	pool := tunny.NewFunc(opts.PoolSize, func(payload interface{}) interface{} {
		fn, ok := payload.(executable)
		if !ok {
			return fmt.Errorf("payload is not an executable: %T", payload)
		}
		return fn()
	})
	return &connPool{
		p:  pool,
		to: opts.TransactionTimeout,
	}
}

func (p *connPool) Exec(exec executable) error {
	resp, err := p.p.ProcessTimed(exec, p.to)
	if err != nil {
		return fmt.Errorf("failed to execute job: %v", err)
	}

	if resp == nil {
		return nil
	}

	execErr, ok := resp.(error)
	if !ok {
		return fmt.Errorf("unexpected non-error return type: %T", resp)
	}
	return execErr
}

func (p *connPool) Close() error {
	return p.Close()
}
