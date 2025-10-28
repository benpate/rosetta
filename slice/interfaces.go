package slice

type stringOKGetter interface {
	GetStringOK(string) (string, bool)
}
