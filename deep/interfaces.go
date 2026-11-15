package deep

type Copier interface {
	Copy() (any, error)
}

type Remover interface {
	Remove(...string) error
}
