package meander

// Facade exposes the public method
type Facade interface {
	Public() interface{}
}

// Public implements Facade interface to check if it has a Public() interface{} method,
// if it does then it calls the method and returns the results otherwise it just returns
// the original object untouched
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
