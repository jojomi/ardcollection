package main

// ARD fake feed
type Feed struct {
	Result struct {
		Episodes []struct {
			Enclosure struct {
				DownloadURL string `json:"download_url"`
			} `json:"enclosure"`
			Title string `json:"title"`
		} `json:"episodes"`
	} `json:"result"`
}
