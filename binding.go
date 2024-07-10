package container

type Lifetime string

const (
	Singleton Lifetime = "singleton"
	Transient Lifetime = "transient"
	Scoped    Lifetime = "scoped"
)

// binding holds a resolver and a concrete (if already resolved).
// It is the break for the Container wall!
type binding struct {
	resolver interface{} // resolver is the function that is responsible for making the concrete.
	concrete interface{} // concrete is the stored instance for singleton bindings.
	lifetime Lifetime
}

// make resolves the binding if needed and returns the resolved concrete.
func (b *binding) make(c *Container) (interface{}, error) {
	if b.concrete != nil {
		return b.concrete, nil
	}

	retVal, err := c.invoke(b.resolver)
	if b.lifetime == Singleton && err == nil {
		b.concrete = retVal
	}

	return retVal, err
}
