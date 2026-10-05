package types

type FocalPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (f FocalPoint) IsZero() bool {
	return f.X == 0 && f.Y == 0
}

func (f FocalPoint) IsValid() bool {
	return f.X >= 0 && f.X <= 1 && f.Y >= 0 && f.Y <= 1
}
