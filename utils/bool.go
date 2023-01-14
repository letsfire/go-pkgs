package utils

const (
	True  Bool = 1
	False Bool = 2
)

type Bool int

func (b Bool) Not() Bool {
	if b.IsTrue() {
		return False
	} else {
		return True
	}
}

func (b Bool) IsTrue() bool {
	return b == True
}

func (b Bool) IsFalse() bool {
	return b == False
}

func (b Bool) Validate() bool {
	return b.IsTrue() || b.IsFalse()
}

func NewBool(value bool) Bool {
	if value {
		return True
	}
	return False
}
