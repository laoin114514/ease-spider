package debug

type Debug struct{}

func NewDebug() *Debug {
	return &Debug{}
}
func (d *Debug) Debug(v ...any) {

}
