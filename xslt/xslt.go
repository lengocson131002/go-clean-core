package xslt

type Xslt interface {
	Transform(style []byte, input []byte) ([]byte, error)
}
