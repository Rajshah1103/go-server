package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseRequestLine(line string) (method, path, version string) {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

func ReadHeaders(reader *bufio.Reader) map[string]string {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		line = strings.TrimSpace(line)
		colon := strings.Index(line, ":")
		if colon != -1 {
			key := strings.TrimSpace(line[:colon])
			val := strings.TrimSpace(line[colon+1:])
			headers[key] = val
		}
	}
	return headers
}

func BuildHTTPResponse(body string, statusCode ...int) string {
	code := 200
	statusText := "OK"

	if len(statusCode) > 0 {
		switch statusCode[0] {
		case 404:
			code = 404
			statusText = "Not Found"
		case 400:
			code = 400
			statusText = "Bad Request"
		case 500:
			code = 500
			statusText = "Internal Server Error"
		}
	}

	return fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s",
		code, statusText, len(body), body,
	)
}

func ServeStatic(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return BuildHTTPResponse("500 Internal Server Error", 500)
	}
	contentType := guessContentType(filePath)
	header := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n",
		len(data), contentType,
	)
	return header + string(data)
}

func guessContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
