package real_debrid

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"qdebrid/logger"
	"sync"
)

func AddTorrent(torrents []io.Reader) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(torrents))

	for _, torrent := range torrents {
		wg.Add(1)
		go func(torrent io.Reader) {
			defer wg.Done()
			url, _ := url.Parse(apiHost)
			url.Path += apiPath + "/torrents/addTorrent"

			req, err := http.NewRequest("PUT", url.String(), torrent)
			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer response.Body.Close()

			if response.StatusCode != 201 {
				errChan <- fmt.Errorf("failed to add torrent: %v", response.Status)
			}
		}(torrent)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
