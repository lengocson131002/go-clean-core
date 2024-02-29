package xslt

import (
	"github.com/wamuir/go-xslt"
)

type dXslt struct {
}

func NewDefaultXslt() Xslt {
	return &dXslt{}
}

// Transform implements ports.Xslt.
func (*dXslt) Transform(style []byte, input []byte) ([]byte, error) {
	xs, err := xslt.NewStylesheet(style)
	if err != nil {
		panic(err)
	}
	defer xs.Close()

	// doc is an XML document to be transformed and res is the result of
	// the XSL transformation, both as []byte.
	res, err := xs.Transform(input)
	if err != nil {
		return nil, err
	}

	return res, nil
}
