package utils

import (
	"fmt"
)

type Debug struct {
	isUse bool
}

func NewDebug(isUse bool) *Debug {
	return &Debug{
		isUse: isUse,
	}
}
func (d *Debug) Debug(v ...any) {
	if d.isUse {
		fmt.Print("[Debug]:")
		fmt.Println(v...)
	}
}
func (d *Debug) Debugf(format string, v ...any) {
	if d.isUse {
		fmt.Print("[Debug]:")
		fmt.Printf(format, v...)
	}
}
