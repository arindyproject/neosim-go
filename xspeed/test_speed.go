package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Endpoint Target yang akan diuji
type Endpoint struct {
	Name string
	URL  string
}

func main() {
	// Token Authorization
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InN1cGVyYWRtaW4iLCJpc19zdXBlcmFkbWluIjp0cnVlLCJpc19zdGFmZiI6dHJ1ZSwidG9rZW5fdHlwZSI6ImFjY2VzcyIsImlzcyI6Im5lb3NpbSIsImV4cCI6MTc4OTAwNzcyNywiaWF0IjoxNzg4NDAyOTI3LCJqdGkiOiJkNjFjNTg3Zi0xY2NhLTRmY2UtYjcxNi1hYTU3MjRmN2U0YzgifQ.ZowHs54Thm4x7C0AVhxHDARCjub4Qnrk1Xu7iBjy3Bk"
	base_url := "http://localhost:1323/api/v1/"
	page_size := 100000
	// Jumlah perulangan pengujian per endpoint
	repeatCount := 10

	// Daftar endpoint yang ingin diuji
	endpoints := []Endpoint{
		{
			Name: "Get Artikel ",
			URL:  base_url + "artikel?page_size=" + fmt.Sprint(page_size),
		},
		{
			Name: "Get Kategori Artikel",
			URL:  base_url + "artikel/kategori?page_size=" + fmt.Sprint(page_size),
		},
		{
			Name: "Get Tag Kategori Artikel",
			URL:  base_url + "artikel/kategori/tags?page_size=" + fmt.Sprint(page_size),
		},
		{
			Name: "Get Identifier Kepegawaian",
			URL:  base_url + "kepegawaian/identifier?page_size=" + fmt.Sprint(page_size),
		},
		{
			Name: "Get Identifier Kepegawaian By PegawaiID",
			URL:  base_url + "kepegawaian/identifier/3/pegawai?page_size=" + fmt.Sprint(page_size),
		},
		{
			Name: "Get Kualifikasi Kepegawaian",
			URL:  base_url + "kepegawaian/kualifikasi?page_size=" + fmt.Sprint(page_size),
		},
	}

	fmt.Printf("=== MEMULAI PENGUJIAN RESPOBN TIME (%d KALI PER ENDPOINT) ===\n\n", repeatCount)

	for _, ep := range endpoints {
		fmt.Printf("--- Testing Endpoint: %s ---\nURL: %s\n", ep.Name, ep.URL)

		var totalDuration time.Duration
		successCount := 0

		for i := 1; i <= repeatCount; i++ {
			duration, statusCode, err := sendRequest(ep.URL, token)
			if err != nil {
				fmt.Printf("  Iteration #%d: Gagal - %v\n", i, err)
				continue
			}

			totalDuration += duration
			successCount++
			fmt.Printf("  Iteration #%d: Status [%d] | Respon Time: %d ms (%.2f ms)\n",
				i, statusCode, duration.Milliseconds(), float64(duration.Microseconds())/1000.0)
		}

		// Hitung Rata-rata per endpoint
		if successCount > 0 {
			avgDuration := totalDuration / time.Duration(successCount)
			fmt.Printf("=> Rata-rata Respon Time [%s]: %d ms (%.2f ms)\n\n",
				ep.Name, avgDuration.Milliseconds(), float64(avgDuration.Microseconds())/1000.0)
		} else {
			fmt.Printf("=> Seluruh iterasi gagal untuk [%s]\n\n", ep.Name)
		}
	}
}

// sendRequest melakukan HTTP GET request dan mengembalikan durasi eksekusi
func sendRequest(url string, token string) (time.Duration, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal mengirim request: %w", err)
	}
	defer resp.Body.Close()

	// Membaca penuh body respon agar perhitungan waktu akurat
	_, _ = io.Copy(io.Discard, resp.Body)

	duration := time.Since(start)

	return duration, resp.StatusCode, nil
}
