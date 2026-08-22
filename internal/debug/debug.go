package debug

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
)

func PrintMem(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s:\n", label)
	fmt.Printf("  Alloc: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("  HeapInuse: %.2f MB\n", float64(m.HeapInuse)/1024/1024)
}

func Debug() {
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
