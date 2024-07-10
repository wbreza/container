package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wbreza/container/v4"
)

func TestRegisterInstance(t *testing.T) {
	container.Reset()

	err := container.RegisterInstance(&Circle{a: 13})
	assert.NoError(t, err)
}

func TestRegisterNamedInstance(t *testing.T) {
	container.Reset()

	err := container.RegisterNamedInstance("rounded", &Circle{a: 13})
	assert.NoError(t, err)
}

func TestRegisterSingleton(t *testing.T) {
	container.Reset()

	err := container.RegisterSingleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestRegisterNamedSingleton(t *testing.T) {
	container.Reset()

	err := container.RegisterNamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestRegisterTransient(t *testing.T) {
	container.Reset()

	err := container.RegisterTransient(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestRegisterNamedTransient(t *testing.T) {
	container.Reset()

	err := container.RegisterNamedTransient("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestRegisterScoped(t *testing.T) {
	container.Reset()

	err := container.RegisterScoped(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestRegisterNamedScoped(t *testing.T) {
	container.Reset()

	err := container.RegisterNamedScoped("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestCall(t *testing.T) {
	container.Reset()

	err := container.Call(func() {})
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {
	container.Reset()

	var s Shape

	err := container.RegisterSingleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = container.Resolve(&s)
	assert.NoError(t, err)
}

func TestResolveNamed(t *testing.T) {
	container.Reset()

	var s Shape

	err := container.RegisterNamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = container.ResolveNamed(&s, "rounded")
	assert.NoError(t, err)
}

func TestFill(t *testing.T) {
	container.Reset()

	err := container.RegisterSingleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape `Global:"type"`
	}{}

	err = container.Fill(&myApp)
	assert.NoError(t, err)
}
