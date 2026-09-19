
package main
import (
	"fmt"
	"github.com/akilama471/WorkBench/internal/app"
	"github.com/akilama471/WorkBench/internal/logger"
	"os"
)
func main() {
	log := logger.New(logger.LevelDebug, os.Stdout)
	a, _ := app.New("I:\\Project\\Hobby\\WorkBench\\build")
	_ = log
	svc, _ := a.ServiceManager.GetService("mariadb")
	fmt.Printf("Status: %s\n", svc.Status())
}

