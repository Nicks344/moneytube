//go:build !release

package enigma

func init() {
	impl = &fakeImpl{}
}

type fakeImpl struct {
}

func (this *fakeImpl) Init() {

}

func (this *fakeImpl) GetHWID() (string, error) {
	return "6B32F84-AB5B1F1-E53AF0F-1A5C40B", nil
}

func (this *fakeImpl) CheckAndSaveKey(name string, key string) bool {
	return true
}

func (this *fakeImpl) LoadAndCheckKey() bool {
	return true
}

func (this *fakeImpl) GetDaysLeft() int {
	return 1
}
