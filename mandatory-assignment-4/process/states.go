package process

type State int

const (
	RELEASED State = iota
	WANTED
	HELD
)

func (s State) String() string {
	switch s {
	case RELEASED:
		return "RELEASED"
	case WANTED:
		return "WANTED"
	case HELD:
		return "HELD"
	default:
		return "UNKNOWN"
	}
}
