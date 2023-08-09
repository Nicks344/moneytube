//go:build !release

package serverAPI

func init() {
	host = "http://127.0.0.1:4444"
}
