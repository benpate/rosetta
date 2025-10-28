package sliceof

type stringOKGetter interface {
	GetStringOK(string) (string, bool)
}
