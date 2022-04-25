package checker

type TypeChek int

const (
	IntervalT       = TypeChek(1)
	ChangeUpT       = TypeChek(2)
	ChangeDownT     = TypeChek(3)
	DateT           = TypeChek(4)
	MoreT           = TypeChek(5)
	LessT           = TypeChek(6)
	EqualT          = TypeChek(7)
	CheckCardanoVer = TypeChek(8)
	CheckBool       = TypeChek(9)
)

var chekersTypes = map[TypeChek]string{
	IntervalT:       "Interval",
	ChangeUpT:       "ChangeUp",
	ChangeDownT:     "ChangeDown",
	DateT:           "Date",
	MoreT:           "More",
	LessT:           "Less",
	EqualT:          "Equal",
	CheckCardanoVer: "CheckCardanoVer",
	CheckBool:       "CheckBool",
}

func (tc TypeChek) String() string {
	return chekersTypes[tc]
}
