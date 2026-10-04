package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func fetchMovie(ctx context.Context, id int, timeout time.Duration) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)
	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "id %d: %v\n", id, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "id %d: %v\n", id, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "id %d: status %d\n", id, resp.StatusCode)
		return
	}
	var m Movie
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		fmt.Fprintf(os.Stderr, "id %d: bad json\n", id)
		return
	}
	fmt.Printf("%d — %s — %d — %s\n", m.ID, m.Title, m.Year, m.Director)
}
func main() {
	from := flag.Int("from", -1, "")
	to := flag.Int("to", -1, "")
	workers := flag.Int("workers", 10, "")
	timeout := flag.Duration("timeout", 5*time.Second, "")
	flag.Parse()
	if *from == -1 || *to == -1 || *from > *to || *workers <= 0 || *timeout <= 0 {
		fmt.Fprintf(os.Stderr, "bad flags\n")
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	jobs := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				if ctx.Err() != nil {
					return
				}
				fetchMovie(ctx, id, *timeout)
			}
		}()
	}
	for id := *from; id <= *to; id++ {
		if ctx.Err() != nil {
			break
		}
		jobs <- id
	}
	close(jobs)
	wg.Wait()
}
