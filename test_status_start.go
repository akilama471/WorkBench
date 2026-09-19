
package main
import (
	"fmt"
	"github.com/akilama471/WorkBench/internal/app"
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)
func main() {
	paths := filesystem.NewPaths("I:\\Project\\Hobby\\WorkBench\\build")
	log := logger.New(paths)
	a := app.New(paths, log)
	a.ServiceManager.Start("mariadb")
	svc, _ := a.ServiceManager.GetService("mariadb")
	fmt.Printf("Status immediately after start: %s\n", svc.Status())
}

