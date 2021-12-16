package demos

type Runnable interface {
	Run()
}

type Stoppable interface {
	Stop()
}

// Bus 用组合的方式来代替继承，因此支持"多重继承"
type Bus struct {
	Runnable
	Stoppable
}

func (b *Bus) Extends() {
	b.Run()
	b.Stop()
}
