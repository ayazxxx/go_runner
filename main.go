package main

import (
	"fmt"
	"os"
	"sync"
	"time"

    "github.com/ayazxxx/go_runner/engine"
    "github.com/ayazxxx/go_runner/utils/logger"
    "github.com/ayazxxx/go_runner/utils/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: ./main <domain>")
		return
	}

	target := os.Args[1]
	start := time.Now()

	fmt.Println("[*] Tarama başlatılıyor:", target)

	// HTTP isteklerini ZAP'ten al
	requests, err := parser.ExtractRequests("webapp/requests/" + parser.ExtractDomain(target))
	if err != nil {
		fmt.Println("[!] İstekler alınamadı:", err)
		return
	}

	fmt.Printf("[*] Toplam %d istek bulundu.\n", len(requests))

	var wg sync.WaitGroup

	// Her bir istek için ayrı tarama işlemi başlat
	for _, req := range requests {
		wg.Add(1)
		go func(r parser.Request) {
			defer wg.Done()
			scanner.RunSQLMapScan(r)
		}(req)
	}

	wg.Wait()

	fmt.Println("[+] Tarama tamamlandı:", target)
	fmt.Println("[✓] Süre:", time.Since(start))
}
