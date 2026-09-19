
package main
import (
	"fmt"
	"github.com/akilama471/WorkBench/internal/app"
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)
func main() {
	paths := filesystem.NewPaths("I:\\Project\\Hobby\\WorkBench\\build")
	log := logger.NewLogger(paths)
	a := app.NewApplication(paths, log)
	svc, _ := a.ServiceManager.GetService("mariadb")
	fmt.Printf("Status: %s\n", svc.Status())
}

