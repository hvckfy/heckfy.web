package rql

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func RQLConnectToTarantool() (*tarantool.Connection, error) {
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
	if err != nil {
		return nil, fmt.Errorf("error connecting to Tarantool: %w", err)
	}
	return conn, nil
}
