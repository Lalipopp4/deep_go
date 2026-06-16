package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

var (
	errUndefinedType = errors.New("undefined type")
)

type UserService struct {
	NotEmptyStruct bool
}
type MessageService struct {
	NotEmptyStruct bool
}

type Constructor interface {
	interface{} | int
}

type Container[T Constructor] struct {
	constructors map[string]func() T
}

func NewContainer[T Constructor]() *Container[T] {
	return &Container[T]{
		constructors: make(map[string]func() T),
	}
}

func (c *Container[T]) RegisterType(name string, constructor func() T) {
	c.constructors[name] = constructor
}

func (c *Container[T]) Resolve(name string) (T, error) {
	cons, ok := c.constructors[name]
	if !ok {
		var zero T
		return zero, errUndefinedType
	}

	return cons(), nil

}

func (c *Container[T]) RegisterSingletonType(name string, constructor func() T) {
	obj := constructor()
	c.constructors[name] = func() T {
		return obj
	}
}

func TestDIContainer(t *testing.T) {
	container := NewContainer[interface{}]()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}

func TestRegisterSingletonType(t *testing.T) {
	container := NewContainer[interface{}]()
	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	assert.True(t, userService1 == userService2)
}
