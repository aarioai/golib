package enum

type Act string

const (
	UnknownAct Act = "unknown"
	Read       Act = "read"
	Write      Act = "write"
	ReadWrite  Act = "(read)|(write)"
)

func (a Act) String() string {
	return string(a)
}
