package manifest

type Manifest struct {
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	Author      string              `json:"author"`
	Repo        string              `json:"repo"`
	Bin         string              `json:"bin"`
	Platforms   map[string]Platform `json:"platforms"`
}

type Platform struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}
