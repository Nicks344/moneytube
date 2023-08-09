package enigma

var impl iEnigma

type iEnigma interface {
	Init()

	GetHWID() (string, error)

	CheckAndSaveKey(name string, key string) bool

	LoadAndCheckKey() bool

	GetDaysLeft() int
}

func GetHWID() (string, error) {
	return impl.GetHWID()
}

func CheckAndSaveKey(name string, key string) bool {
	return impl.CheckAndSaveKey(name, key)
}

func LoadAndCheckKey() bool {
	return impl.LoadAndCheckKey()
}

func GetDaysLeft() int {
	return impl.GetDaysLeft()
}
