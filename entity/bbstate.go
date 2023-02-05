package entity

type BBState int

const (
	None BBState = iota
	UnderLower
	UnderBasic
	UnderUpper
	OverUpper
)

var bbStateMap = map[BBState]string{
	None:       "",
	UnderLower: "underLower",
	UnderBasic: "underBasic",
	UnderUpper: "underUpper",
	OverUpper:  "overUpper",
}

func (state BBState) ToString() string {
	return bbStateMap[state]
}
