package container

// RegisterInstanceAs registers an instance as a specific type within the container
func RegisterInstanceAs[T any](c *Container, instance T) error {
	return RegisterNamedInstanceAs(c, "", instance)
}

// RegisterNamedInstanceAs registers an instance as a specific type within the container with a name
func RegisterNamedInstanceAs[T any](c *Container, name string, instance T) error {
	options := RegisterOptions{
		Name: name,
		Resolver: func() T {
			return instance
		},
		Lifetime: Singleton,
	}

	return c.Register(options)
}
