package main

import (
	"github.com/Pis0sion/rblog/internal/rblog"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// app entry
func main() {

	// rand seed
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	rblog.NewApp("rblog").Run()
}
