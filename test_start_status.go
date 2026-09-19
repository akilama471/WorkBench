
package main
import (
	"fmt"
	"github.com/akilama471/WorkBench/internal/app"
	"github.com/akilama471/WorkBench/internal/logger"
	"os"
	"time"
)
func main() {
	log := logger.New(logger.LevelDebug, os.Stdout)
	a, _ := app.New("I:\\Project\\Hobby\\WorkBench\\build")
	_ = log
	a.ServiceManager.Start("mariadb")
	svc, _ := a.ServiceManager.GetService("mariadb")
	fmt.Printf("Status immediately: %s\n", svc.Status())
	time.Sleep(1 * time.Second)
	fmt.Printf("Status after 1s: %s\n", svc.Status())
}

