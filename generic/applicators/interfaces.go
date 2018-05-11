package applicators

type Applicator interface {
	Apply([]interface{})
}

type Locker interface {
	Lock()
	Unlock()
}
