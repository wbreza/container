package container_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wbreza/container/v4"
)

func TestMustRegisterSingleton_It_Should_Panic_On_Error(t *testing.T) {
	expectedErr := fmt.Sprintf("%s, the resolver must be a function", container.ErrInvalidResolver)

	assert.PanicsWithError(t, expectedErr, func() {
		c := container.New()
		container.MustRegisterSingleton(c, "not a resolver function")
	})
}

func TestMustNamedSingleton_It_Should_Panic_On_Error(t *testing.T) {
	expectedErr := fmt.Sprintf("%s, the resolver must be a function", container.ErrInvalidResolver)

	assert.PanicsWithError(t, expectedErr, func() {
		c := container.New()
		container.MustRegisterNamedSingleton(c, "name", "not a resolver function")
	})
}

func TestMustRegisterTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	container.MustRegisterTransient(c, func() (Shape, error) {
		return nil, errors.New("custom error")
	})

	assert.PanicsWithError(t, "failed making instance for type 'container_test.Shape'. Error: custom error", func() {
		var resVal Shape
		container.MustResolve(context.Background(), c, &resVal)
	})
}

func TestMustRegisterNamedTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for type 'container_test.Shape' with name 'name'. Error: no binding found for abstraction 'container_test.Shape'"

	container.MustRegisterNamedTransient(c, "custom-name", func() (Shape, error) {
		return nil, errors.New("custom error")
	})

	assert.PanicsWithError(t, expectedErr, func() {
		var resVal Shape
		container.MustNamedResolve(context.Background(), c, &resVal, "name")
	})
}

func TestMustRegisterScoped_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	container.MustRegisterScoped(c, func() (Shape, error) {
		return nil, errors.New("custom error")
	})

	assert.PanicsWithError(t, "failed making instance for type 'container_test.Shape'. Error: custom error", func() {
		scope, err := c.NewScope()
		assert.NoError(t, err)

		var resVal Shape
		container.MustResolve(context.Background(), scope, &resVal)
	})
}

func TestMustRegisterNamedScoped_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for type 'container_test.Shape' with name 'name'. Error: no binding found for abstraction 'container_test.Shape'"

	container.MustRegisterNamedScoped(c, "custom-name", func() (Shape, error) {
		return nil, errors.New("custom error")
	})

	assert.PanicsWithError(t, expectedErr, func() {
		scope, err := c.NewScope()
		assert.NoError(t, err)

		var resVal Shape
		container.MustNamedResolve(context.Background(), scope, &resVal, "name")
	})
}

func TestMustRegisterInstance_It_Should_Panic_On_Error(t *testing.T) {
	expectedErr := "invalid resolver, cannot register a function as an instance"

	assert.PanicsWithError(t, expectedErr, func() {
		c := container.New()
		container.MustRegisterInstance(c, func() {})
	})
}

func TestMustRegisterNamedInstance_It_Should_Panic_On_Error(t *testing.T) {
	expectedErr := "invalid resolver, cannot register a function as an instance"

	assert.PanicsWithError(t, expectedErr, func() {
		c := container.New()
		container.MustRegisterNamedInstance(c, "name", func() {})
	})
}

func TestMustCall_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for type 'container_test.Shape', Error: no binding found for abstraction 'container_test.Shape'"

	assert.PanicsWithError(t, expectedErr, func() {
		container.MustCall(context.Background(), c, func(s Shape) {
			s.GetArea()
		})
	})
}

func TestMustResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for type 'container_test.Shape'. Error: no binding found for abstraction 'container_test.Shape'"

	assert.PanicsWithError(t, expectedErr, func() {
		var s Shape
		container.MustResolve(context.Background(), c, &s)
	})
}

func TestMustNamedResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for type 'container_test.Shape' with name 'name'. Error: no binding found for abstraction 'container_test.Shape'"

	assert.PanicsWithError(t, expectedErr, func() {
		var s Shape
		container.MustNamedResolve(context.Background(), c, &s, "name")
	})
}

func TestMustFill_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()
	expectedErr := "failed making instance for field 'S', Error: no binding found for abstraction 'container_test.Shape'"

	myApp := struct {
		S Shape `container:"type"`
	}{}

	assert.PanicsWithError(t, expectedErr, func() {
		container.MustFill(context.Background(), c, &myApp)
	})
}
