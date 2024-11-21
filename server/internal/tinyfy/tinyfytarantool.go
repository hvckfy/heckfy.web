package tinyfy

import (
	"context"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func connectToTarantool() (*tarantool.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  Address,
		User:     User,
		Password: Password,
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}
	conn, err := tarantool.Connect(ctx, dialer, opts)
	return conn, err
}
