package huflit

type Register struct {
	Name       string
	Code       string
	FirstCode  string
	SecondCode string
	RequestId  string
}

func (r *Register) IsSingleSubject() bool {
	return r.SecondCode == ""
}
