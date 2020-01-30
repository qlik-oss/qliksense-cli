//go:generate packr2
package qliksense

import (
	"github.com/gobuffalo/packr/v2"
)

// Qliksense is the logic behind the qliksense client
type Qliksense struct {
	PorterExe     string
	QliksenseHome string
	CrdBox        *packr.Box
}

// New qliksense client, initialized with useful defaults.
func New(porterExe, qliksenseHome string) *Qliksense {
	return &Qliksense{
		PorterExe:     porterExe,
		QliksenseHome: qliksenseHome,
		CrdBox:        packr.New("crds", "./crds"),
	}
}
