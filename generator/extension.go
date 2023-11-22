package generator

type extension string
type extSet    []extension

// Extension set
var crypto = map[extension]struct{}{
	"zvbb":  {},
	"zvknc": {},
	"zvkng": {},
	"zvksc": {},
	"zvksg": {},
}

func (i *Insn) isExtension(exts map[extension]struct{}) bool {
	for _, insnExt := range i.Ext {
		if _, ok := exts[insnExt]; ok {
			return true
		}
	}
	return false
}