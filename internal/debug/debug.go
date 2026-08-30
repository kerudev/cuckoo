package debug

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Frame = 0
var Iter = 0

var MinFPS = int32(-1)
var MaxFPS = int32(0)

func DrawFPS() {
	rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 0, 0, 20, rl.Black)
}

func WatchMaxFPS() {
	fps := rl.GetFPS()
	if MaxFPS < fps {
		MaxFPS = fps
	}
}

func WatchMinFPS() {
	fps := rl.GetFPS()
	if MinFPS < 0 || MinFPS > fps {
		MinFPS = fps
	}
}

func PrintDebug(label string, v *int, increment bool) {
	fmt.Println(label, *v)
	if increment {
		*v++
	}
}

func PrintMem(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s:\n", label)
	fmt.Printf("  Alloc: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("  HeapInuse: %.2f MB\n", float64(m.HeapInuse)/1024/1024)
}

func DebugServer() {
	// Create a separate mux for pprof
	debugMux := http.NewServeMux()

	// Register pprof handlers manually
	debugMux.HandleFunc("/debug/pprof/", pprof.Index)
	debugMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	debugMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	debugMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	debugMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Start debug server on a different port (not exposed publicly)
	go func() {
		log.Println("Debug server starting on :6060")
		log.Fatal(http.ListenAndServe("localhost:6060", debugMux))
	}()
}
