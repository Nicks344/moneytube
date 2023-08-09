//go:build release

package enigma

import (
	"errors"
	"log"
	"github.com/meandrewdev/logger"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/charmap"
)

func init() {
	impl = &dllImpl{}
	impl.Init()
}

type dllImpl struct {
	lib syscall.Handle
}

func (this *dllImpl) Init() {
	var err error

	this.lib, err = syscall.LoadLibrary("enigma_ide64.dll")
	if err != nil {
		panic(err)
	}
}

func (this *dllImpl) GetHWID() (string, error) {

	addr, err := syscall.GetProcAddress(this.lib, "EP_RegHardwareID")
	if err != nil {
		return "", err
	}
	ret, _, e := syscall.Syscall(addr, uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if e != 0 {
		return "", e
	}
	if ret == 0 {
		return "", errors.New("Не могу получить идентификатор устройства")
	}
	return uintptrToStr(ret), nil
}

func (this *dllImpl) CheckAndSaveKey(name string, key string) bool {

	addr, err := syscall.GetProcAddress(this.lib, "EP_RegCheckAndSaveKey")
	if err != nil {
		return false
	}
	ret, _, e := syscall.Syscall(addr, uintptr(2), strToUintptr(name), strToUintptr(key), uintptr(0))
	if e != 0 {
		return false
	}
	return ret == 1
}

func (this *dllImpl) LoadAndCheckKey() bool {

	addr, err := syscall.GetProcAddress(this.lib, "EP_RegLoadAndCheckKey")
	if err != nil {
		return false
	}
	ret, _, e := syscall.Syscall(addr, uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if e != 0 {
		logger.Error(e)
		return false
	}
	return ret == 1
}

func (this *dllImpl) GetDaysLeft() int {

	addr, err := syscall.GetProcAddress(this.lib, "EP_RegKeyDaysLeft")
	if err != nil {
		return 0
	}
	ret, _, e := syscall.Syscall(addr, uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if e != 0 {
		return 0
	}
	return int(ret)
}

func strToUintptr(s string) uintptr {
	enc := charmap.Windows1252.NewEncoder()
	s1252, err := enc.Bytes([]byte(s))
	if err != nil {
		log.Fatal(err)
	}
	b := append([]byte(s1252), 0)
	return uintptr(unsafe.Pointer(&b[0]))
}

func uintptrToStr(u uintptr) string {
	p := (*byte)(unsafe.Pointer(u))
	data := make([]byte, 0)
	for *p != 0 {
		data = append(data, *p)
		u += unsafe.Sizeof(byte(0))
		p = (*byte)(unsafe.Pointer(u))
	}
	dataStr := toUTF8(string(data))
	return dataStr
}

func toUTF8(enc string) string {
	dec := charmap.Windows1252.NewDecoder()
	out, _ := dec.Bytes([]byte(enc))
	return string(out)
}
