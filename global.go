package container

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

// Singleton calls the same method of the global concrete.
func RegisterSingleton(resolver interface{}) error {
	return Global.RegisterSingleton(resolver)
}

// NamedSingleton calls the same method of the global concrete.
func RegisterNamedSingleton(name string, resolver interface{}) error {
	return Global.RegisterNamedSingleton(name, resolver)
}

// Transient calls the same method of the global concrete.
func RegisterTransient(resolver interface{}) error {
	return Global.RegisterTransient(resolver)
}

// NamedTransient calls the same method of the global concrete.
func RegisterNamedTransient(name string, resolver interface{}) error {
	return Global.RegisterNamedTransient(name, resolver)
}

func RegisterScoped(resolver interface{}) error {
	return Global.RegisterScoped(resolver)
}

func RegisterNamedScoped(name string, resolver interface{}) error {
	return Global.RegisterNamedScoped(name, resolver)
}

// Reset calls the same method of the global concrete.
func Reset() {
	Global.Reset()
}

// Call calls the same method of the global concrete.
func Call(receiver interface{}) error {
	return Global.Call(receiver)
}

// Resolve calls the same method of the global concrete.
func Resolve(abstraction interface{}) error {
	return Global.Resolve(abstraction)
}

// NamedResolve calls the same method of the global concrete.
func ResolveNamed(abstraction interface{}, name string) error {
	return Global.ResolvedNamed(abstraction, name)
}

// Fill calls the same method of the global concrete.
func Fill(receiver interface{}) error {
	return Global.Fill(receiver)
}
