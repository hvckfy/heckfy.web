package saveshare

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func tarantool_connect_to_saveshare() (*tarantool.Connection, error) {
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

func add_new_data(data, hash string, connsaveshare *tarantool.Connection) error {
	if connsaveshare == nil {
		var err error
		connsaveshare, err = tarantool_connect_to_saveshare()
		if err != nil {
			return err
		}
	}
	statuscode, err := connsaveshare.Call("add_new_data", []interface{}{data, hash})
	if err != nil {
		fmt.Println(statuscode)
		return fmt.Errorf("failed to add new data to Tarantool: %w", err)
	}
	return nil
}

func get_data_from_tarantool(hash string, connsaveshare *tarantool.Connection) (string, []interface{}, error) {
	// Call the Tarantool function to get the data
	if connsaveshare == nil {
		var err error
		connsaveshare, err = tarantool_connect_to_saveshare()
		if err != nil {
			return "500", nil, err
		}
	}
	result, err := connsaveshare.Call("get_data", []interface{}{hash})
	if err != nil {
		return "500", nil, err
	}

	return result[0].(string), result[1].([]interface{}), nil
}
