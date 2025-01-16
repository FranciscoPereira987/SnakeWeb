package pkg_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"websocket/pkg"

	"github.com/coder/websocket"
)

/*
Connection Dummy that stores what is written into it
and returns objects on reads as long as they are of the same type
*/
type ConnectionDummy struct {
	ReadStack []interface{}
	WriteList []interface{}
	Open      bool
	Status    *websocket.StatusCode
	reason    *string
}

func (c *ConnectionDummy) Read(ctx context.Context, t interface{}) error {
	if !c.Open {
		return errors.New("Connection closed")
	}
	if len(c.ReadStack) == 0 {
		return errors.New("Nothing to read")
	}
	if !reflect.TypeOf(t).Elem().AssignableTo(reflect.TypeOf(c.ReadStack[0])) {
		return errors.New("Type mismatch")
	}
	reflect.ValueOf(t).Elem().Set(reflect.ValueOf(c.ReadStack[0]))
	c.ReadStack = c.ReadStack[1:]
	return nil
}

func (c *ConnectionDummy) Write(ctx context.Context, t interface{}) error {
	if !c.Open {
		return errors.New("Connection closed")
	}
	c.WriteList = append(c.WriteList, t)
	return nil
}

func (c *ConnectionDummy) Close(s websocket.StatusCode, reason string) error {
	if !c.Open {
		return nil
	}
	c.Open = false
	c.Status = &s
	c.reason = &reason
	return nil
}

func TestDummyR(t *testing.T) {
	value := pkg.GameMessage{
		pkg.GET_GAME,
		"some ID",
		nil,
		nil,
	}
	dumm := &ConnectionDummy{
		[]interface{}{
			value,
		},
		[]interface{}{},
		true,
		nil,
		nil,
	}
	var toRead pkg.GameMessage
	if err := dumm.Read(context.Background(), &toRead); err != nil {
		t.Fatalf("Failed with error: %s", err)
	}
	if toRead.Op != pkg.GET_GAME {
		t.Fatalf("Invalid Op, expected %d but got %d", toRead.Op, pkg.GET_GAME)
	}
	if toRead.Id != value.Id {
		t.Fatalf("Invalid ID, expected %s but got %s", value.Id, toRead.Id)
	}
}
