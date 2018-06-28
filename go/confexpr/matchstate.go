package confexpr

type MatchState int

const (
	Missing MatchState = iota
	Incomplete
	Complete
	Collision
)

func (this MatchState) String() string {
	switch this {
	case Missing:
		return "Missing"
	case Incomplete:
		return "Incomplete"
	case Complete:
		return "Complete"
	case Collision:
		return "Collision"
	}
	return "Error"

}
