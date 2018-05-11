package applicators

type Applicator interface {
	Apply([]interface{})
}

type locker interface {
	Lock()
	Unlock()
}
