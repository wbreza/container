// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	// Errors encountered while registering bindings
	ErrInvalidResolver    = errors.New("invalid resolver")
	ErrInvalidAbstraction = errors.New("invalid abstraction")
	ErrInvalidReceiver    = errors.New("invalid receiver")
	ErrInvalidStructure   = errors.New("invalid structure")

	// Errors encountered while resolving, calling or filling
	ErrContextRequired  = errors.New("context is required. If you don't have a context pass 'context.Background()' or 'context.TODO()'")
	ErrResolutionFailed = errors.New("failed making instance")
	ErrBindingNotFound  = errors.New("no binding found")
)

// Container holds the bindings and provides methods to interact with them.
// It is the entry point in the package.
type Container struct {
	parent   *Container
	bindings map[reflect.Type]map[string]*binding
}

// New creates a new instance of the Container.
func New() *Container {
	return &Container{
		bindings: make(map[reflect.Type]map[string]*binding),
	}
}

// NewScope creates a new child container scope.
// Scoped bindings are copied from the parent container to the child container and act as singletons within the new scope
func (c *Container) NewScope() (*Container, error) {
	childContainer := New()
	childContainer.parent = c

	for _, outerBinding := range c.bindings {
		for name, binding := range outerBinding {
			if binding.lifetime == Scoped {
				if err := childContainer.bind(binding.resolver, name, binding.lifetime); err != nil {
					return nil, err
				}
			}
		}
	}

	return childContainer, nil
}

// Reset deletes all the existing bindings and empties the container.
func (c *Container) Reset() {
	for k := range c.bindings {
		delete(c.bindings, k)
	}
}

type RegisterOptions struct {
	Resolver interface{}
	Name     string
	Lifetime Lifetime
}

// Registers the resolver with the specified options.
func (c *Container) Register(options RegisterOptions) error {
	if options.Lifetime == "" {
		options.Lifetime = Singleton
	}

	return c.bind(options.Resolver, options.Name, options.Lifetime)
}

// Invokes the resolver and registers the instance with the specified options.
func (c *Container) InvokeAndRegister(ctx context.Context, options RegisterOptions) error {
	if ctx == nil {
		return ErrContextRequired
	}

	if options.Lifetime == "" {
		options.Lifetime = Singleton
	}

	instance, err := c.invoke(ctx, options.Resolver)
	if err != nil {
		return err
	}

	return c.bind(instance, options.Name, options.Lifetime)
}

// RegisterInstance binds an instance to the container in singleton mode.
func (c *Container) RegisterInstance(instance interface{}) error {
	return c.RegisterNamedInstance("", instance)
}

// RegisterNamedInstance binds an instance to the container in singleton mode with a name.
func (c *Container) RegisterNamedInstance(name string, instance interface{}) error {
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Func {
		return fmt.Errorf("%w, cannot register a function as an instance", ErrInvalidResolver)
	}

	return c.bind(instance, name, Singleton)
}

// Singleton binds an abstraction to concrete in singleton mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c *Container) RegisterSingleton(resolver interface{}) error {
	return c.RegisterNamedSingleton("", resolver)
}

// NamedSingleton binds a named abstraction to concrete in singleton mode.
func (c *Container) RegisterNamedSingleton(name string, resolver interface{}) error {
	t := reflect.TypeOf(resolver)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("%w, the resolver must be a function", ErrInvalidResolver)
	}

	return c.bind(resolver, name, Singleton)
}

// Transient binds an abstraction to concrete in transient mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c *Container) RegisterTransient(resolver interface{}) error {
	return c.RegisterNamedTransient("", resolver)
}

// NamedTransient binds a named abstraction to concrete lazily in transient mode.
func (c *Container) RegisterNamedTransient(name string, resolver interface{}) error {
	t := reflect.TypeOf(resolver)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("%w, the resolver must be a function", ErrInvalidResolver)
	}

	return c.bind(resolver, name, Transient)
}

// Scoped binds an abstraction to concrete in scoped mode.
func (c *Container) RegisterScoped(resolver interface{}) error {
	return c.RegisterNamedScoped("", resolver)
}

// NamedScoped binds a named abstraction to concrete in scoped mode.
func (c *Container) RegisterNamedScoped(name string, resolver interface{}) error {
	t := reflect.TypeOf(resolver)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("%w, the resolver must be a function", ErrInvalidResolver)
	}

	return c.bind(resolver, name, Scoped)
}

// Call takes a receiver function with one or more arguments of the abstractions (interfaces).
// It invokes the receiver function and passes the related concretes.
func (c *Container) Call(ctx context.Context, function interface{}) error {
	if ctx == nil {
		return ErrContextRequired
	}

	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return ErrInvalidReceiver
	}

	arguments, err := c.arguments(ctx, function)
	if err != nil {
		return err
	}

	result := reflect.ValueOf(function).Call(arguments)

	if len(result) == 0 {
		return nil
	} else if len(result) == 1 && result[0].CanInterface() {
		if result[0].IsNil() {
			return nil
		}
		if err, ok := result[0].Interface().(error); ok {
			return err
		}
	}

	return ErrInvalidReceiver
}

// ResolveWithContext takes an abstraction and a context and fills it with the related concrete.
func (c *Container) Resolve(ctx context.Context, abstraction interface{}) error {
	return c.ResolveNamed(ctx, "", abstraction)
}

// ResolveNamed takes abstraction and its name and fills it with the related concrete.
func (c *Container) ResolveNamed(ctx context.Context, name string, abstraction interface{}) error {
	if ctx == nil {
		return ErrContextRequired
	}

	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return ErrInvalidAbstraction
	}

	if receiverType.Kind() != reflect.Ptr {
		return ErrInvalidAbstraction
	}

	elem := receiverType.Elem()

	if instance, err := c.make(ctx, elem, name); err == nil {
		reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
		return nil
	} else {
		if name == "" {
			return fmt.Errorf("%w for type '%s'. Error: %w", ErrResolutionFailed, elem.String(), err)
		} else {
			return fmt.Errorf("%w for type '%s' with name '%s'. Error: %w", ErrResolutionFailed, elem.String(), name, err)
		}
	}
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func (c *Container) Fill(ctx context.Context, structure interface{}) error {
	if ctx == nil {
		return ErrContextRequired
	}

	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return ErrInvalidStructure
	}

	if receiverType.Kind() != reflect.Ptr {
		return ErrInvalidStructure
	}

	elem := receiverType.Elem()
	if elem.Kind() != reflect.Struct {
		return ErrInvalidStructure
	}

	s := reflect.ValueOf(structure).Elem()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if t, exist := s.Type().Field(i).Tag.Lookup("container"); exist {
			var name string

			if t == "type" {
				name = ""
			} else if t == "name" {
				name = s.Type().Field(i).Name
			} else {
				return fmt.Errorf("%w, %v has an invalid struct tag", ErrInvalidStructure, s.Type().Field(i).Name)
			}

			if instance, err := c.make(ctx, f.Type(), name); err == nil {
				ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(instance))

				continue
			} else {
				return fmt.Errorf("%w for field '%v', Error: %w", ErrResolutionFailed, s.Type().Field(i).Name, err)
			}
		}
	}

	return nil
}

// Validate checks the container for any errors and ensures all registered types can be resolved.
func (c *Container) Validate(ctx context.Context) error {
	if ctx == nil {
		return ErrContextRequired
	}

	for t, binding := range c.bindings {
		for name := range binding {
			if _, err := c.make(ctx, t, name); err != nil {
				return err
			}
		}
	}

	return nil
}

// bind maps an abstraction to concrete and instantiates if it is a singleton binding.
func (c *Container) bind(resolver interface{}, name string, lifetime Lifetime) error {
	reflectedResolver := reflect.TypeOf(resolver)

	// For function based bindings
	if reflectedResolver.Kind() == reflect.Func {
		if reflectedResolver.NumOut() > 0 {
			if _, exist := c.bindings[reflectedResolver.Out(0)]; !exist {
				c.bindings[reflectedResolver.Out(0)] = make(map[string]*binding)
			}
		}

		if err := c.validateResolverFunction(reflectedResolver); err != nil {
			return err
		}

		c.bindings[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, lifetime: lifetime}

	} else { // For instance based bindings
		if _, exist := c.bindings[reflectedResolver]; !exist {
			c.bindings[reflectedResolver] = make(map[string]*binding)
		}

		c.bindings[reflectedResolver][name] = &binding{concrete: resolver, lifetime: lifetime}
	}

	return nil
}

func (c *Container) validateResolverFunction(funcType reflect.Type) error {
	retCount := funcType.NumOut()

	if retCount == 0 || retCount > 2 {
		return fmt.Errorf("%w, signature is invalid - it must return abstract, or abstract and error", ErrInvalidResolver)
	}

	resolveType := funcType.Out(0)
	for i := 0; i < funcType.NumIn(); i++ {
		if funcType.In(i) == resolveType {
			return fmt.Errorf("%w, signature is invalid - depends on abstract it returns", ErrInvalidResolver)
		}
	}

	return nil
}

// invoke calls a function and its returned values.
// It only accepts one value and an optional error.
func (c *Container) invoke(ctx context.Context, function interface{}) (interface{}, error) {
	arguments, err := c.arguments(ctx, function)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(function).Call(arguments)
	if len(values) == 2 && values[1].CanInterface() {
		if err, ok := values[1].Interface().(error); ok {
			return values[0].Interface(), err
		}
	}
	return values[0].Interface(), nil
}

// make resolves the binding and returns the concrete.
// Search up any parent container scopes if the binding is not found in current scope.
func (c *Container) make(ctx context.Context, t reflect.Type, name string) (interface{}, error) {
	current := c
	var binding *binding

	for {
		if found, exist := current.bindings[t][name]; exist {
			binding = found
			break
		}

		if current.parent == nil {
			break
		} else {
			current = current.parent
		}
	}

	if binding == nil {
		return nil, fmt.Errorf("%w for abstraction '%s'", ErrBindingNotFound, t.String())
	}

	return binding.make(ctx, c)
}

// arguments returns the list of resolved arguments for a function.
func (c *Container) arguments(ctx context.Context, function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if abstraction.Implements(contextType) {
			arguments[i] = reflect.ValueOf(ctx)
		} else {
			if instance, err := c.make(ctx, abstraction, ""); err == nil {
				arguments[i] = reflect.ValueOf(instance)
			} else {
				return nil, fmt.Errorf("%w for type '%s', Error: %w", ErrResolutionFailed, abstraction.String(), err)
			}
		}
	}

	return arguments, nil
}
