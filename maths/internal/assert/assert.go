package assert

type Equality struct {
	Message  string
	Expected string
	Got      string
	IsEqual  bool
}

func (e Equality) Error() string {
	return e.Message + ", Expected " + e.Expected + " Got " + e.Got
}
func (e Equality) String() string { return e.Error() }

func ErrorEquality(expErr, gotErr error) Equality {
	if expErr != gotErr {
		if expErr == nil && gotErr != nil {
			return Equality{
				"unexpected error",
				"nil",
				gotErr.Error(),
				false,
			}
		}
		if expErr != nil && gotErr == nil {
			return Equality{
				"expected error",
				expErr.Error(),
				"nil",
				false,
			}
		}
		if expErr.Error() != gotErr.Error() {
			return Equality{
				"incorrect error value",
				expErr.Error(),
				gotErr.Error(),
				false,
			}
		}
		return Equality{IsEqual: true}
	}
	if expErr != nil {
		return Equality{}
	}
	return Equality{IsEqual: true}
}
