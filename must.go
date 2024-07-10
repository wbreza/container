package container

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
func MustCall(c *Container, receiver interface{}) {
	if err := c.Call(receiver); err != nil {
		panic(err)
	}
}

// MustResolve wraps the `Resolve` method and panics on errors instead of returning the errors.
func MustResolve(c *Container, abstraction interface{}) {
	if err := c.Resolve(abstraction); err != nil {
		panic(err)
	}
}

// MustNamedResolve wraps the `NamedResolve` method and panics on errors instead of returning the errors.
func MustNamedResolve(c *Container, abstraction interface{}, name string) {
	if err := c.ResolvedNamed(abstraction, name); err != nil {
		panic(err)
	}
}

// MustFill wraps the `Fill` method and panics on errors instead of returning the errors.
func MustFill(c *Container, receiver interface{}) {
	if err := c.Fill(receiver); err != nil {
		panic(err)
	}
}
