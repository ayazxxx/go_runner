package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Filename        string
	TargetURL       string
	Method          string
	PostData        string
	Headers         []string
	Cookies         string
	InjectionPoints []string
}

// ParseRequest parses a ZAP-style request .txt file
func ParseRequest(filename string) (*Request, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("dosya açılamadı: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	req := &Request{Filename: filename}

	headerSection := true
	var rawBodyLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// İlk satır: GET /path HTTP/1.1 veya POST /abc.php HTTP/1.1
		if strings.HasPrefix(line, "GET") || strings.HasPrefix(line, "POST") {
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				req.Method = parts[0]
				req.TargetURL = parts[1]
			}
			continue
		}

		if headerSection {
			if line == "" {
				headerSection = false
				continue
			}

			// Header ve Cookie ayrımı
			if strings.HasPrefix(strings.ToLower(line), "cookie:") {
				req.Cookies = strings.TrimSpace(strings.TrimPrefix(line, "Cookie:"))
			} else {
				req.Headers = append(req.Headers, line)
			}
		} else {
			rawBodyLines = append(rawBodyLines, line)
		}
	}

	req.PostData = strings.Join(rawBodyLines, "&")

	// Injection noktalarını bul
	req.InjectionPoints = detectInjectionParams(req)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner hatası: %v", err)
	}

	return req, nil
}

// detectInjectionParams finds potential injection points
func detectInjectionParams(req *Request) []string {
	var points []string

	// Query string parametreleri
	if strings.Contains(req.TargetURL, "?") {
		parts := strings.Split(req.TargetURL, "?")
		params := strings.Split(parts[1], "&")
		for _, p := range params {
			if kv := strings.SplitN(p, "=", 2); len(kv) == 2 {
				points = append(points, kv[0])
			}
		}
	}

	// POST body parametreleri
	if req.PostData != "" {
		params := strings.Split(req.PostData, "&")
		for _, p := range params {
			if kv := strings.SplitN(p, "=", 2); len(kv) == 2 {
				points = append(points, kv[0])
			}
		}
	}

	// Header parametreleri (User-Agent, X-Forwarded-For vs.)
	for _, h := range req.Headers {
		if strings.Contains(h, ":") {
			key := strings.SplitN(h, ":", 2)[0]
			points = append(points, key)
		}
	}

	// Cookie parametreleri
	if req.Cookies != "" {
		pairs := strings.Split(req.Cookies, ";")
		for _, c := range pairs {
			if kv := strings.SplitN(strings.TrimSpace(c), "=", 2); len(kv) == 2 {
				points = append(points, kv[0])
			}
		}
	}

	return points
}
