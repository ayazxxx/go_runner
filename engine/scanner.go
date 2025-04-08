package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang_engine/utils/logger"  // Doğru yola göre import edin
	"golang_engine/utils/parser"  // Doğru yola göre import edin
)

// RunSQLMapScan, tek bir HTTP isteğini SQLMap ile analiz eder
func RunSQLMapScan(req parser.Request) {
	targetFile := fmt.Sprintf("webapp/requests/%s", req.Filename)
	params := append(req.QueryParams, req.BodyParams...)
	params = append(params, req.MultipartParams...)
	params = append(params, req.HeaderParams...)
	params = append(params, req.CookieParams...)
	params = append(params, req.PathParams...)

	logger.LogInfo(fmt.Sprintf("📄 Başlatıldı: %s (%d parametre)", targetFile, len(params)))

	var wg sync.WaitGroup
	for _, param := range params {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			runSQLMap(targetFile, p)
		}(param)
	}
	wg.Wait()
	logger.LogInfo("🎯 SQLMap taraması tamamlandı.")
}

type ScanResult struct {
	Filename   string        `json:"filename"`
	Parameter  string        `json:"parameter"`
	Vulnerable bool          `json:"vulnerable"`
	Duration   time.Duration `json:"duration"`
}

// runSQLMap, SQLMap komutunu subprocess ile çalıştırır
func runSQLMap(targetFile string, param string) {
	start := time.Now()

	cmd := exec.Command(
		"python3", "sqlmap/sqlmap.py",
		"-r", targetFile,
		"-p", param,
		"--batch",
		"--level=5", "--risk=3",
		"--output-dir=output",
	)
	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	result := ScanResult{
		Filename:   filepath.Base(targetFile),
		Parameter:  param,
		Vulnerable: strings.Contains(string(out), "is vulnerable"),
		Duration:   duration,
	}

	saveScanResult(result)

	if err != nil {
		logger.LogError(fmt.Sprintf("HATA (%s): %v", param, err))
	} else if result.Vulnerable {
		logger.LogInfo(fmt.Sprintf("🚨 ZAFİYET TESPİT EDİLDİ (%s)", param))
	} else {
		logger.LogInfo(fmt.Sprintf("✅ Temiz (%s) [%s]", param, duration))
	}
}

// saveScanResult, tarama sonucunu JSON olarak log dosyasına kaydeder
func saveScanResult(result ScanResult) {
	file, err := os.OpenFile("output/scan_results.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogError("Log dosyasına yazılamadı: " + err.Error())
		return
	}
	defer file.Close()

	jsonData, _ := json.Marshal(result)
	file.WriteString(string(jsonData) + "\n")
}
