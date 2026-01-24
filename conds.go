package conds

func IF(expr bool, vthen, velse string) string {
	if expr {
		return vthen
	} else {
		return velse
	}
}
