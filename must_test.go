package container_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wbreza/container/v4"
)

func TestMustRegisterSingleton_It_Should_Panic_On_Error(t *testing.T) {
	assert.PanicsWithError(t, "container: the resolver must be a function", func() {
		c := container.New()
		container.MustRegisterSingleton(c, "not a resolver function")
	})
}

func TestMustNamedSingleton_It_Should_Panic_On_Error(t *testing.T) {
	assert.PanicsWithError(t, "container: the resolver must be a function", func() {
		c := container.New()
		container.MustRegisterNamedSingleton(c, "name", "not a resolver function")
	})
}

func TestMustRegisterTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustRegisterTransient(c, func() (Shape, error) {
		return nil, errors.New("error")
	})

	var resVal Shape
	container.MustResolve(c, &resVal)

	t.Errorf("panic expected.")
}

func TestMustRegisterNamedTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustRegisterNamedTransient(c, "name", func() (Shape, error) {
		return nil, errors.New("error")
	})

	var resVal Shape
	container.MustNamedResolve(c, &resVal, "name")

	t.Errorf("panic expcted.")
}

func TestMustCall_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustCall(c, func(s Shape) {
		s.GetArea()
	})
	t.Errorf("panic expcted.")
}

func TestMustResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	var s Shape

	defer func() { recover() }()
	container.MustResolve(c, &s)
	t.Errorf("panic expcted.")
}

func TestMustNamedResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	var s Shape

	defer func() { recover() }()
	container.MustNamedResolve(c, &s, "name")
	t.Errorf("panic expcted.")
}

func TestMustFill_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	myApp := struct {
		S Shape `container:"type"`
	}{}

	defer func() { recover() }()
	container.MustFill(c, &myApp)
	t.Errorf("panic expcted.")
}
