package container

import "context"

// Global is the global concrete of the Container.
var Global = New()

// RegisterInstance calls the same method of the global concrete.
func RegisterInstance(instance interface{}) error {
	return Global.RegisterInstance(instance)
}

// RegisterNamedInstance calls the same method of the global concrete.
func RegisterNamedInstance(name string, instance interface{}) error {
	return Global.RegisterNamedInstance(name, instance)
}

// RegisterSingleton calls the same method of the global concrete.
func RegisterSingleton(resolver interface{}) error {
	return Global.RegisterSingleton(resolver)
}

// RegisterNamedSingleton calls the same method of the global concrete.
func RegisterNamedSingleton(name string, resolver interface{}) error {
	return Global.RegisterNamedSingleton(name, resolver)
}

// RegisterTransient calls the same method of the global concrete.
func RegisterTransient(resolver interface{}) error {
	return Global.RegisterTransient(resolver)
}

// RegisterNamedTransient calls the same method of the global concrete.
func RegisterNamedTransient(name string, resolver interface{}) error {
	return Global.RegisterNamedTransient(name, resolver)
}

// RegisterScoped calls the same method of the global concrete.
func RegisterScoped(resolver interface{}) error {
	return Global.RegisterScoped(resolver)
}

// RegisterNamedScoped calls the same method of the global concrete.
func RegisterNamedScoped(name string, resolver interface{}) error {
	return Global.RegisterNamedScoped(name, resolver)
}

// Reset calls the same method of the global concrete.
func Reset() {
	Global.Reset()
}

// Call calls the same method of the global concrete.
func Call(ctx context.Context, receiver interface{}) error {
	return Global.Call(ctx, receiver)
}

// Resolve calls the same method of the global concrete.
func Resolve(ctx context.Context, abstraction interface{}) error {
	return Global.Resolve(ctx, abstraction)
}

// ResolveNamed calls the same method of the global concrete.
func ResolveNamed(ctx context.Context, abstraction interface{}, name string) error {
	return Global.ResolveNamed(ctx, name, abstraction)
}

// Fill calls the same method of the global concrete.
func Fill(ctx context.Context, receiver interface{}) error {
	return Global.Fill(ctx, receiver)
}
