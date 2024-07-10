// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
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

func (c *Container) NewScope() (*Container, error) {
	childContainer := New()
	childContainer.parent = c

	for _, outerBinding := range c.bindings {
		for name, binding := range outerBinding {
			if binding.lifetime == Scoped {
				if err := childContainer.bind(binding.resolver, name, Singleton); err != nil {
					return nil, err
				}
			}
		}
	}

	return childContainer, nil
}

// bind maps an abstraction to concrete and instantiates if it is a singleton binding.
func (c *Container) bind(resolver interface{}, name string, lifetime Lifetime) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("container: the resolver must be a function")
	}

	if reflectedResolver.NumOut() > 0 {
		if _, exist := c.bindings[reflectedResolver.Out(0)]; !exist {
			c.bindings[reflectedResolver.Out(0)] = make(map[string]*binding)
		}
	}

	if err := c.validateResolverFunction(reflectedResolver); err != nil {
		return err
	}

	c.bindings[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, lifetime: lifetime}

	return nil
}

func (c *Container) validateResolverFunction(funcType reflect.Type) error {
	retCount := funcType.NumOut()

	if retCount == 0 || retCount > 2 {
		return errors.New("container: resolver function signature is invalid - it must return abstract, or abstract and error")
	}

	resolveType := funcType.Out(0)
	for i := 0; i < funcType.NumIn(); i++ {
		if funcType.In(i) == resolveType {
			return fmt.Errorf("container: resolver function signature is invalid - depends on abstract it returns")
		}
	}

	return nil
}

// invoke calls a function and its returned values.
// It only accepts one value and an optional error.
func (c *Container) invoke(function interface{}) (interface{}, error) {
	arguments, err := c.arguments(function)
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
func (c *Container) make(t reflect.Type, name string) (interface{}, error) {
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
		return nil, errors.New("container: no binding found for: " + t.String())
	}

	return binding.make(current)
}

// arguments returns the list of resolved arguments for a function.
func (c *Container) arguments(function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if instance, err := c.make(abstraction, ""); err == nil {
			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, fmt.Errorf("container: encountered error while making instance for: %s. Error encountered: %w", abstraction.String(), err)
		}
	}

	return arguments, nil
}

// Reset deletes all the existing bindings and empties the container.
func (c *Container) Reset() {
	for k := range c.bindings {
		delete(c.bindings, k)
	}
}

// RegisterInstance binds an instance to the container in singleton mode.
func (c *Container) RegisterInstance(instance interface{}) error {
	return c.RegisterNamedInstance("", instance)
}

// RegisterNamedInstance binds an instance to the container in singleton mode with a name.
func (c *Container) RegisterNamedInstance(name string, instance interface{}) error {
	t := reflect.TypeOf(instance)

	if t.Kind() == reflect.Func {
		return errors.New("container: cannot register a function as an instance")
	}

	c.bindings[t] = map[string]*binding{
		name: {
			concrete: instance,
			lifetime: Singleton,
		},
	}

	return nil
}

// Singleton binds an abstraction to concrete in singleton mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c *Container) RegisterSingleton(resolver interface{}) error {
	return c.bind(resolver, "", Singleton)
}

// NamedSingleton binds a named abstraction to concrete in singleton mode.
func (c *Container) RegisterNamedSingleton(name string, resolver interface{}) error {
	return c.bind(resolver, name, Singleton)
}

// Transient binds an abstraction to concrete in transient mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c *Container) RegisterTransient(resolver interface{}) error {
	return c.bind(resolver, "", Transient)
}

// NamedTransient binds a named abstraction to concrete lazily in transient mode.
func (c *Container) RegisterNamedTransient(name string, resolver interface{}) error {
	return c.bind(resolver, name, Transient)
}

// Scoped binds an abstraction to concrete in scoped mode.
func (c *Container) RegisterScoped(resolver interface{}) error {
	return c.bind(resolver, "", Scoped)
}

// NamedScoped binds a named abstraction to concrete in scoped mode.
func (c *Container) RegisterNamedScoped(name string, resolver interface{}) error {
	return c.bind(resolver, name, Scoped)
}

// Call takes a receiver function with one or more arguments of the abstractions (interfaces).
// It invokes the receiver function and passes the related concretes.
func (c *Container) Call(function interface{}) error {
	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return errors.New("container: invalid function")
	}

	arguments, err := c.arguments(function)
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

	return errors.New("container: receiver function signature is invalid")
}

// Resolve takes an abstraction (reference of an interface type) and fills it with the related concrete.
func (c *Container) Resolve(abstraction interface{}) error {
	return c.ResolvedNamed(abstraction, "")
}

// ResolvedNamed takes abstraction and its name and fills it with the related concrete.
func (c *Container) ResolvedNamed(abstraction interface{}, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("container: invalid abstraction")
	}

	if receiverType.Kind() != reflect.Ptr {
		return errors.New("container: invalid abstraction")

	}

	elem := receiverType.Elem()

	if instance, err := c.make(elem, name); err == nil {
		reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
		return nil
	} else {
		return fmt.Errorf("container: encountered error while making instance for: %s. Error encountered: %w", elem.String(), err)
	}
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func (c *Container) Fill(structure interface{}) error {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return errors.New("container: invalid structure")
	}

	if receiverType.Kind() != reflect.Ptr {
		return errors.New("container: invalid structure")
	}

	elem := receiverType.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("container: invalid structure")
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
				return fmt.Errorf("container: %v has an invalid struct tag", s.Type().Field(i).Name)
			}

			if instance, err := c.make(f.Type(), name); err == nil {
				ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(instance))

				continue
			} else {
				return fmt.Errorf("container: encountered error while making %v field. Error encountered: %w", s.Type().Field(i).Name, err)
			}
		}
	}

	return nil
}
