package container_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wbreza/container/v4"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c Circle) GetArea() int {
	return c.a
}

type Square struct {
	a int
}

func (s *Square) SetArea(a int) {
	s.a = a
}

func (s Square) GetArea() int {
	return s.a
}

type DatabaseOptions struct {
	Host     string
	Port     int
	Username string
}

type Database interface {
	Connect() bool
	Options() *DatabaseOptions
}

type MySQL struct {
	options *DatabaseOptions
}

func (m MySQL) Connect() bool {
	return true
}

func (m MySQL) Options() *DatabaseOptions {
	return m.options
}

type SqlServer struct {
	options *DatabaseOptions
}

func (s SqlServer) Connect() bool {
	return true
}

func (s SqlServer) Options() *DatabaseOptions {
	return s.options
}

var instance = container.New()

func TestContainer_RegisterSingleton(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s1 Shape) {
		s1.SetArea(666)
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_RegisterSingleton_With_Missing_Dependency_Resolve(t *testing.T) {
	err := instance.RegisterSingleton(func(db Database) Shape {
		return &Circle{a: 13}
	})

	assert.NoError(t, err)

	var resolved Shape
	err = instance.Resolve(context.Background(), &resolved)
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_RegisterSingleton_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.RegisterSingleton(func() {})
	assert.Error(t, err, "container: resolver function signature is invalid")
}

func TestContainer_RegisterSingleton_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.RegisterSingleton("not a resolver")
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_RegisterSingleton_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.RegisterSingleton("STRING!")
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_RegisterSingleton_With_Resolvable_Arguments(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), 666)
		return &MySQL{}
	})
	assert.NoError(t, err)
}

func TestContainer_RegisterSingleton_With_Non_Resolvable_Arguments(t *testing.T) {
	instance.Reset()

	err := instance.RegisterSingleton(func(s Shape) Shape {
		return &Circle{a: s.GetArea()}
	})
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_RegisterNamedSingleton(t *testing.T) {
	err := instance.RegisterNamedSingleton("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.ResolvedNamed(context.Background(), &sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestContainer_RegisterTransient(t *testing.T) {
	err := instance.RegisterTransient(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s1 Shape) {
		s1.SetArea(13)
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_RegisterTransient_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.RegisterTransient(func() {})
	assert.Error(t, err, "container: resolver function signature is invalid")
}

func TestContainer_RegisterTransient_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.RegisterTransient(func() (Shape, error) {
		return nil, errors.New("app: error resolving Shape")
	})

	assert.NoError(t, err)

	err = instance.RegisterTransient(func() (Database, error) {
		return nil, errors.New("app: error resolving Database")
	})
	assert.NoError(t, err)

	var db Database
	err = instance.Resolve(context.Background(), &db)
	assert.Error(t, err, "app: error resolving Database")
}

func TestContainer_RegisterTransient_With_Resolve_With_Invalid_Signature_It_Should_Fail(t *testing.T) {
	err := instance.RegisterTransient(func() (Shape, Database, error) {
		return nil, nil, nil
	})
	assert.Error(t, err, "container: resolver function signature is invalid")
}

func TestContainer_RegisterNamedTransient(t *testing.T) {
	err := instance.RegisterNamedTransient("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.ResolvedNamed(context.Background(), &sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestContainer_Call_With_Multiple_Resolving(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.NoError(t, err)
}

func TestContainer_Call_With_Dependency_Missing_In_Chain(t *testing.T) {
	var instance = container.New()
	err := instance.RegisterSingleton(func() (Database, error) {
		var s Shape
		if err := instance.Resolve(context.Background(), &s); err != nil {
			return nil, err
		}
		return &MySQL{}, nil
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(m Database) {
		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Call_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Call(context.Background(), "STRING!")
	assert.ErrorIs(t, err, container.ErrInvalidReceiver)
}

func TestContainer_Call_With_Second_UnBounded_Argument(t *testing.T) {
	instance.Reset()

	err := instance.RegisterSingleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s Shape, d Database) {})
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Call_With_A_Returning_Error(t *testing.T) {
	instance.Reset()

	err := instance.RegisterSingleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s Shape) error {
		return errors.New("app: some context error")
	})
	assert.EqualError(t, err, "app: some context error")
}

func TestContainer_Call_With_A_Returning_Nil_Error(t *testing.T) {
	instance.Reset()

	err := instance.RegisterSingleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s Shape) error {
		return nil
	})
	assert.Nil(t, err)
}

func TestContainer_Call_With_Invalid_Signature(t *testing.T) {
	instance.Reset()

	err := instance.RegisterSingleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(context.Background(), func(s Shape) (int, error) {
		return 13, errors.New("app: some context error")
	})
	assert.ErrorIs(t, err, container.ErrInvalidReceiver)
}

func TestContainer_Resolve_With_Reference_As_Resolver(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	var (
		s Shape
		d Database
	)

	err = instance.Resolve(context.Background(), &s)
	assert.NoError(t, err)
	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	err = instance.Resolve(context.Background(), &d)
	assert.NoError(t, err)
	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestContainer_Resolve_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Resolve(context.Background(), "STRING!")
	assert.ErrorIs(t, err, container.ErrInvalidAbstraction)
}

func TestContainer_Resolve_With_NonReference_Receiver_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Resolve(context.Background(), s)
	assert.ErrorIs(t, err, container.ErrInvalidAbstraction)
}

func TestContainer_Resolve_With_UnBounded_Reference_It_Should_Fail(t *testing.T) {
	instance.Reset()

	var s Shape
	err := instance.Resolve(context.Background(), &s)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Resolve_With_Error_Should_Not_Cache_Concrete(t *testing.T) {
	c := container.New()

	resolveCount := 0
	err := c.RegisterSingleton(func() (Shape, error) {
		resolveCount++
		if resolveCount == 1 {
			return nil, errors.New("first resolve error")
		}

		return &Circle{a: 5}, nil
	})

	assert.NoError(t, err)

	var s Shape

	err = c.Resolve(context.Background(), &s)
	assert.Error(t, err, "first resolve error")
	assert.Nil(t, s)

	err = c.Resolve(context.Background(), &s)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestContainer_Fill_With_Struct_Pointer(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterNamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `container:"type"`
		D Database `container:"type"`
		C Shape    `container:"name"`
		X string
	}{}

	err = instance.Fill(context.Background(), &myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.S)
	assert.IsType(t, &MySQL{}, myApp.D)
}

func TestContainer_Fill_Unexported_With_Struct_Pointer(t *testing.T) {
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape    `container:"type"`
		d Database `container:"type"`
		y int
	}{}

	err = instance.Fill(context.Background(), &myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.s)
	assert.IsType(t, &MySQL{}, myApp.d)
}

func TestContainer_Fill_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	err := instance.RegisterNamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	err = instance.Fill(context.Background(), &myApp)
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Fill_With_Invalid_Tag_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `container:"invalid"`
	}

	myApp := App{}

	err := instance.Fill(context.Background(), &myApp)
	assert.ErrorIs(t, err, container.ErrInvalidStructure)
}

func TestContainer_Fill_With_Invalid_Field_Name_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	err := instance.Fill(context.Background(), &myApp)
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Fill_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	invalidStruct := 0
	err := instance.Fill(context.Background(), &invalidStruct)
	assert.ErrorIs(t, err, container.ErrInvalidStructure)
}

func TestContainer_Fill_With_Invalid_Pointer_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Fill(context.Background(), s)
	assert.ErrorIs(t, err, container.ErrInvalidStructure)
}

func TestContainer_Fill_With_Dependency_Missing_In_Chain(t *testing.T) {
	var instance = container.New()
	err := instance.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.RegisterNamedSingleton("C", func() (Shape, error) {
		var s Shape
		if err := instance.ResolvedNamed(context.Background(), &s, "foo"); err != nil {
			return nil, err
		}
		return &Circle{a: 5}, nil
	})
	assert.NoError(t, err)

	err = instance.RegisterSingleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `container:"type"`
		D Database `container:"type"`
		C Shape    `container:"name"`
		X string
	}{}

	err = instance.Fill(context.Background(), &myApp)
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_ResolveScoped_Is_Same_Instance_Within_Scope(t *testing.T) {
	root := container.New()

	root.RegisterScoped(func() Database {
		return &MySQL{}
	})

	scope, err := root.NewScope()
	assert.NoError(t, err)

	var database1 Database
	err = scope.Resolve(context.Background(), &database1)
	assert.NoError(t, err)

	var database2 Database
	err = scope.Resolve(context.Background(), &database2)
	assert.NoError(t, err)

	// Both instances are resolved from the same scope so the same cached instance should be returned for both.
	assert.Same(t, database1, database2)
}

func TestContainer_ResolveScoped_At_Root_Acts_Like_Singleton(t *testing.T) {
	root := container.New()

	root.RegisterScoped(func() Database {
		return &MySQL{}
	})

	var db1 Database
	err := root.Resolve(context.Background(), &db1)
	assert.NoError(t, err)
	assert.NotNil(t, db1)

	var db2 Database
	err = root.Resolve(context.Background(), &db2)
	assert.NoError(t, err)
	assert.NotNil(t, db2)

	// When scoped elements are resolved at the root container, they act like singleton elements.
	assert.Same(t, db1, db2)
}

func TestContainer_ResolveScoped_With_Singleton_Dependency(t *testing.T) {
	root := container.New()

	root.RegisterSingleton(func() *DatabaseOptions {
		return &DatabaseOptions{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
		}
	})

	root.RegisterScoped(func(options *DatabaseOptions) Database {
		return &MySQL{
			options: options,
		}
	})

	scope1, err := root.NewScope()
	assert.NoError(t, err)

	var database1 Database
	err = scope1.Resolve(context.Background(), &database1)
	assert.NoError(t, err)

	scope2, err := root.NewScope()
	assert.NoError(t, err)

	var database2 Database
	err = scope2.Resolve(context.Background(), &database2)
	assert.NoError(t, err)

	assert.NotSame(t, database1, database2)
	assert.Same(t, database1.Options(), database2.Options())
}

func TestContainer_ResolveScoped_With_Transient_Dependency(t *testing.T) {
	root := container.New()

	root.RegisterTransient(func() *DatabaseOptions {
		return &DatabaseOptions{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
		}
	})

	root.RegisterScoped(func(options *DatabaseOptions) Database {
		return &MySQL{
			options: options,
		}
	})

	scope1, err := root.NewScope()
	assert.NoError(t, err)

	var database1 Database
	err = scope1.Resolve(context.Background(), &database1)
	assert.NoError(t, err)

	scope2, err := root.NewScope()
	assert.NoError(t, err)

	var database2 Database
	err = scope2.Resolve(context.Background(), &database2)
	assert.NoError(t, err)

	assert.NotSame(t, database1, database2)
	assert.NotSame(t, database1.Options(), database2.Options())
}

func TestContainer_Fill_With_Scoped_Elements(t *testing.T) {
	root := container.New()
	root.RegisterNamedScoped("square", func() Shape {
		return &Square{a: 10}
	})

	root.RegisterNamedScoped("circle", func() Shape {
		return &Circle{a: 5}
	})

	type request struct {
		square Shape `container:"name"`
		circle Shape `container:"name"`
	}

	scope, err := root.NewScope()
	assert.NoError(t, err)

	var req1 request

	err = scope.Fill(context.Background(), &req1)
	assert.NoError(t, err)
	assert.NotNil(t, req1)

	assert.IsType(t, &Square{}, req1.square)
	assert.IsType(t, &Circle{}, req1.circle)

	var req2 request
	err = scope.Fill(context.Background(), &req2)
	assert.NoError(t, err)
	assert.NotNil(t, req2)

	assert.Same(t, req1.square, req2.square)
	assert.Same(t, req1.circle, req2.circle)
}

func TestContainer_Call_With_Scoped_Elements(t *testing.T) {
	root := container.New()

	root.RegisterScoped(func() Shape {
		return &Circle{a: 5}
	})

	scope, err := root.NewScope()
	assert.NoError(t, err)

	// First call should already have area set to 5 from the resolver
	err = scope.Call(context.Background(), func(s1 Shape) {
		assert.Equal(t, 5, s1.GetArea())
		s1.SetArea(20)
	})

	assert.NoError(t, err)

	// Second call should have the area set to 20 from the previous call
	err = scope.Call(context.Background(), func(s2 Shape) {
		assert.Equal(t, 20, s2.GetArea())
	})

	assert.NoError(t, err)
}

func TestContainer_RegisterScoped_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.RegisterScoped(func() {})
	assert.Error(t, err, "container: resolver function signature is invalid")
}

func TestContainer_RegisterScoped_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.RegisterScoped("STRING!")
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_ResolveInstance(t *testing.T) {
	c := container.New()
	circle := &Circle{a: 5}
	err := c.RegisterInstance(circle)
	assert.NoError(t, err)

	var resolvedCircle *Circle
	err = c.Resolve(context.Background(), &resolvedCircle)
	assert.NoError(t, err)
	assert.Same(t, circle, resolvedCircle)
}

func TestContainer_ResolveInstance_With_Invalid_Receiver(t *testing.T) {
	c := container.New()
	err := c.RegisterInstance(func() Database {
		return &MySQL{}
	})
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_ResolveInstance_With_Value(t *testing.T) {
	c := container.New()
	var i int = 5
	err := c.RegisterInstance(i)
	assert.NoError(t, err)

	var resolvedInt int
	err = c.Resolve(context.Background(), &resolvedInt)
	assert.NoError(t, err)

	assert.Equal(t, i, resolvedInt)
}

func TestContainer_ResolveNamedInstance(t *testing.T) {
	c := container.New()
	circle := &Circle{a: 5}
	err := c.RegisterNamedInstance("circle", circle)
	assert.NoError(t, err)

	var resolvedCircle *Circle
	err = c.ResolvedNamed(context.Background(), &resolvedCircle, "circle")
	assert.NoError(t, err)
	assert.Same(t, circle, resolvedCircle)
}

func TestContainer_ResolveInstance_As_Dependency(t *testing.T) {
	c := container.New()
	value := "value"
	err := c.RegisterInstance(value)
	assert.NoError(t, err)

	err = c.RegisterSingleton(func(s string) Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	var s Shape
	err = c.Resolve(context.Background(), &s)
	assert.NoError(t, err)
}

func TestContainer_ResolvedNamedInstance_With_Invalid_Receiver(t *testing.T) {
	c := container.New()
	err := c.RegisterNamedInstance("circle", func() Database {
		return &MySQL{}
	})
	assert.ErrorIs(t, err, container.ErrInvalidResolver)
}

func TestContainer_Validate_With_Empty(t *testing.T) {
	c := container.New()
	err := c.Validate(context.Background())
	assert.NoError(t, err)
}

func TestContainer_Validate_All_Valid(t *testing.T) {
	c := container.New()
	err := c.RegisterSingleton(func(s Shape) Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	err = c.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = c.Validate(context.Background())
	assert.NoError(t, err)
}

func TestContainer_Validate_With_Missing_Dependency(t *testing.T) {
	c := container.New()
	err := c.RegisterSingleton(func(db Database) Shape {
		return &Circle{a: 5}
	})

	assert.NoError(t, err)
	err = c.Validate(context.Background())
	assert.ErrorIs(t, err, container.ErrResolutionFailed)
	assert.ErrorIs(t, err, container.ErrBindingNotFound)
}

func TestContainer_Resolve_With_Custom_Context(t *testing.T) {
	c := container.New()

	ctx := context.Background()
	var refCtx context.Context

	c.RegisterSingleton(func(innerCtx context.Context) Shape {
		refCtx = innerCtx
		return &Circle{a: 5}
	})

	var s Shape
	err := c.Resolve(ctx, &s)
	assert.NoError(t, err)
	assert.Equal(t, refCtx, ctx)
}

func TestContainer_Resolve_Nil_Context(t *testing.T) {
	c := container.New()

	var s Shape
	err := c.Resolve(nil, &s)
	assert.ErrorIs(t, err, container.ErrContextRequired)
	assert.Nil(t, s)
}

func TestContainer_Call_With_Custom_Context(t *testing.T) {
	c := container.New()

	ctx := context.Background()

	err := c.RegisterSingleton(func() Shape {
		return &Circle{a: 5}
	})

	assert.NoError(t, err)

	err = c.Call(ctx, func(refCtx context.Context, s Shape) {
		assert.Equal(t, refCtx, ctx)
		assert.NotNil(t, s)
	})

	assert.NoError(t, err)
}

func TestContainer_Call_Nil_Context(t *testing.T) {
	c := container.New()

	err := c.Call(nil, func(s Shape) {
		assert.NotNil(t, s)
	})

	assert.ErrorIs(t, err, container.ErrContextRequired)
}

func TestContainer_NewScope_Nested_Scopes(t *testing.T) {
	c := container.New()
	called := 0

	err := c.RegisterScoped(func() Shape {
		called++
		return &Circle{}
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, called)

	// Resolve the same type twice, since it's scoped, the resolver should only one once, the second time we just return the cached instance.
	var resolved1 Shape
	err = c.Resolve(context.Background(), &resolved1)
	assert.NoError(t, err)
	assert.Equal(t, 1, called)

	var resolved2 Shape
	err = c.Resolve(context.Background(), &resolved2)
	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Same(t, resolved1, resolved2)

	// Create a new scope, and then resolve in that scope, the resolver should be called again, since we are in a new scope.
	sub, err := c.NewScope()
	assert.NoError(t, err)

	var resolved3 Shape
	err = sub.Resolve(context.Background(), &resolved3)
	assert.NoError(t, err)
	assert.Equal(t, 2, called)

	var resolved4 Shape
	err = sub.Resolve(context.Background(), &resolved4)
	assert.NoError(t, err)
	assert.Equal(t, 2, called)
	assert.Same(t, resolved3, resolved4)

	// Now, create a scope from this container we got from NewScope on the previous container and run the resolvers again.
	sub2, err := sub.NewScope()
	assert.NoError(t, err)

	var resolved5 Shape
	err = sub2.Resolve(context.Background(), &resolved5)
	assert.NoError(t, err)
	assert.Equal(t, 3, called)

	var resolved6 Shape
	err = sub2.Resolve(context.Background(), &resolved6)
	assert.NoError(t, err)
	assert.Equal(t, 3, called)
	assert.Same(t, resolved5, resolved6)
}
