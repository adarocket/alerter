package checker

type TypeChek int

const (
	IntervalT TypeChek = 1 + iota
	ChangeUpT
	ChangeDownT
	DateT
	MoreT
	LessT
	EqualT
)

var chekersTypes = [...]string{
	"Interval",
	"ChangeUp",
	"ChangeDown",
	"Date",
	"More",
	"Less",
	"Equal",
}

func (tc TypeChek) String() string {
	return chekersTypes[tc-1]
}
