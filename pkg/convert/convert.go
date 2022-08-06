package convert

import "strconv"

// 对string进行二次封装
type StrTo string

func (s *StrTo) String() string {
	return string(*s)
}

// 转化成整数返回
func (s *StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s *StrTo) MustInt() int {
	v, _ := s.Int()

	return v
}

func (s *StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())

	return uint32(v), err
}

func (s *StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()

	return v
}
