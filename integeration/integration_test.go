package integeration_test

import (
	"net/http"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/myshkin5/hasher/integeration"
	"github.com/myshkin5/hasher/logs"
)

func BenchmarkIntegration(b *testing.B) {
	cmd := exec.Command("go", "run", "../api/main.go")
	go func() {
		err := cmd.Run()
		if err != nil {
			b.Errorf("API exited with non-zero status: %s", err.Error())
		}

		logs.Logger.Info("API shutdown complete.")
	}()
	waitForStartup()

	logs.Logger.Info("Successfully started API")
	b.ResetTimer()

	requestMgr := requestManager{requesterCount: 1000, requestCount: b.N, b: b}
	requestMgr.start()
	requestMgr.wait()

	logs.Logger.Info("Shutting down API...")

	response, err := http.Post("http://localhost:8080/shutdown", "", nil)
	if err != nil {
		b.Errorf("Could not shutdown API: %s", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		b.Errorf("Did not receive OK on shutdown: %d", response.StatusCode)
	}
	response.Body.Close()
}

func waitForStartup() {
	for {
		response, _ := http.Get("http://localhost:8080/stats")
		if response != nil && response.StatusCode == http.StatusOK {
			break
		}
		response.Body.Close()
		time.Sleep(100 * time.Millisecond)
	}
}

type requestManager struct {
	requesterCount int
	requestCount   int
	b              *testing.B

	wg *sync.WaitGroup
}

func (m *requestManager) start() {
	m.wg = &sync.WaitGroup{}
	m.wg.Add(m.requesterCount)
	integeration.StartRequesters(m.wg, m.requesterCount, m.requestCount)
}

func (m *requestManager) wait() {
	m.wg.Wait()
}
