package moneylib

type MoneyService interface {
	FloatToInt(value float32) int
	IntToFloat32(value int) float32
}

func NewMoneyService(delimiter int) MoneyService {
	return service{Delimiter: delimiter}
}

type service struct {
	Delimiter int
}

func (s service) FloatToInt(value float32) int {
	return int(value * float32(s.Delimiter))
}

func (s service) IntToFloat32(value int) float32 {
	return float32(value) / float32(s.Delimiter)
}
