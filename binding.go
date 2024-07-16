package container

import "context"

type Lifetime string

const (
	// Singleton lifetime means that the container will resolve the binding only once and return the same instance every time.
	Singleton Lifetime = "singleton"
	// Transient lifetime means that the container will resolve the binding every time it is requested.
	Transient Lifetime = "transient"
	// Scoped lifetime means that the container will resolve the binding once per created scope.
	Scoped Lifetime = "scoped"
)

// binding holds a resolver and a concrete (if already resolved).
// It is the break for the Container wall!
type binding struct {
	resolver interface{} // resolver is the function that is responsible for making the concrete.
	concrete interface{} // concrete is the stored instance for singleton / scoped bindings.
	lifetime Lifetime
}

// make resolves the binding if needed and returns the resolved concrete.
func (b *binding) make(ctx context.Context, c *Container) (interface{}, error) {
	if b.concrete != nil {
		return b.concrete, nil
	}

	retVal, err := c.invoke(ctx, b.resolver)
	if b.lifetime != Transient && err == nil {
		b.concrete = retVal
	}

	return retVal, err
}
