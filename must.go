package container

import "context"

// MustRegisterInstance wraps the `RegisterInstance` method and panics on errors instead of returning the errors.
func MustRegisterInstance(c *Container, instance interface{}) {
	if err := c.RegisterInstance(instance); err != nil {
		panic(err)
	}
}

// MustRegisterNamedInstance wraps the `RegisterNamedInstance` method and panics on errors instead of returning the errors.
func MustRegisterNamedInstance(c *Container, name string, instance interface{}) {
	if err := c.RegisterNamedInstance(name, instance); err != nil {
		panic(err)
	}
}

// MustRegisterInstanceAs wraps the `RegisterInstanceAs` method and panics on errors instead of returning the errors.
func MustRegisterInstanceAs[T any](c *Container, instance T) {
	if err := RegisterInstanceAs(c, instance); err != nil {
		panic(err)
	}
}

// MustRegisterNamedInstanceAs wraps the `RegisterNamedInstanceAs` method and panics on errors instead of returning the errors.
func MustRegisterNamedInstanceAs[T any](c *Container, name string, instance T) {
	if err := RegisterNamedInstanceAs(c, name, instance); err != nil {
		panic(err)
	}
}

// MustRegisterSingleton wraps the `RegisterSingleton` method and panics on errors instead of returning the errors.
func MustRegisterSingleton(c *Container, resolver interface{}) {
	if err := c.RegisterSingleton(resolver); err != nil {
		panic(err)
	}
}

// MustRegisterNamedSingleton wraps the `RegisterNamedSingleton` method and panics on errors instead of returning the errors.
func MustRegisterNamedSingleton(c *Container, name string, resolver interface{}) {
	if err := c.RegisterNamedSingleton(name, resolver); err != nil {
		panic(err)
	}
}

// MustRegisterTransient wraps the `RegisterTransient` method and panics on errors instead of returning the errors.
func MustRegisterTransient(c *Container, resolver interface{}) {
	if err := c.RegisterTransient(resolver); err != nil {
		panic(err)
	}
}

// MustRegisterNamedTransient wraps the `RegisterNamedTransient` method and panics on errors instead of returning the errors.
func MustRegisterNamedTransient(c *Container, name string, resolver interface{}) {
	if err := c.RegisterNamedTransient(name, resolver); err != nil {
		panic(err)
	}
}

// MustRegisterScoped wraps the `RegisterScoped` method and panics on errors instead of returning the errors.
func MustRegisterScoped(c *Container, resolver interface{}) {
	if err := c.RegisterScoped(resolver); err != nil {
		panic(err)
	}
}

// MustRegisterNamedScoped wraps the `RegisterNamedScoped` method and panics on errors instead of returning the errors.
func MustRegisterNamedScoped(c *Container, name string, resolver interface{}) {
	if err := c.RegisterNamedScoped(name, resolver); err != nil {
		panic(err)
	}
}

// MustCall wraps the `Call` method and panics on errors instead of returning the errors.
func MustCall(ctx context.Context, c *Container, receiver interface{}) {
	if err := c.Call(ctx, receiver); err != nil {
		panic(err)
	}
}

// MustResolve wraps the `Resolve` method and panics on errors instead of returning the errors.
func MustResolve(ctx context.Context, c *Container, abstraction interface{}) {
	if err := c.Resolve(ctx, abstraction); err != nil {
		panic(err)
	}
}

// MustNamedResolve wraps the `NamedResolve` method and panics on errors instead of returning the errors.
func MustResolveNamed(ctx context.Context, c *Container, abstraction interface{}, name string) {
	if err := c.ResolveNamed(ctx, name, abstraction); err != nil {
		panic(err)
	}
}

// MustFill wraps the `Fill` method and panics on errors instead of returning the errors.
func MustFill(ctx context.Context, c *Container, receiver interface{}) {
	if err := c.Fill(ctx, receiver); err != nil {
		panic(err)
	}
}
