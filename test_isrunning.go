
package main
import (
	"fmt"
	"syscall"
)
func main() {
	pid := 2628
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	handle, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	fmt.Printf("OpenProcess err: %v, handle: %v\n", err, handle)
	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)
	fmt.Printf("GetExitCodeProcess err: %v, exitCode: %v\n", err, exitCode)
}

