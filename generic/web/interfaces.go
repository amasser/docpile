package web

type locker interface {
	Lock()
	Unlock()
}
