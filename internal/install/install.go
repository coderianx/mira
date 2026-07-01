package install

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/coderianx/mira/internal/manifest"
	"github.com/coderianx/mira/internal/output"
)

type Record struct {
	Repo        string `json:"repo"`
	Name        string `json:"name"`
	Bin         string `json:"bin"`
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at"`
}

func stateDir() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "mira")
}

func statePath() string {
	return filepath.Join(stateDir(), "state.json")
}

func loadState() (map[string]Record, error) {
	data, err := os.ReadFile(statePath())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Record), nil
		}
		return nil, err
	}
	var records map[string]Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func saveState(records map[string]Record) error {
	if err := os.MkdirAll(stateDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(), data, 0644)
}

func Install(repo string) error {
	owner, name, tag, err := parseRepo(repo)
	if err != nil {
		return fmt.Errorf("invalid repo: %w", err)
	}

	output.Header(repo, "install")

	m, err := fetchManifest(owner, name, tag)
	if err != nil {
		return fmt.Errorf("fetch manifest: %w", err)
	}

	platformKey := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	platform, ok := m.Platforms[platformKey]
	if !ok {
		return fmt.Errorf("no binary for %s (supported: %s)", platformKey, supportedList(m.Platforms))
	}

	output.Infof("Downloading %s %s ...", m.Name, m.Version)
	data, err := download(platform.URL)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	output.Successf("Downloaded %s (%s)", m.Name, humanSize(len(data)))

	if err := verifyChecksum(data, platform.SHA256); err != nil {
		return fmt.Errorf("checksum mismatch: %w", err)
	}
	output.Success("SHA256 checksum verified")

	dir := binDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}

	binPath := filepath.Join(dir, m.Bin)
	if err := os.WriteFile(binPath, data, 0755); err != nil {
		return fmt.Errorf("write %s: %w", binPath, err)
	}

	records, err := loadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	records[repo] = Record{
		Repo:        repo,
		Name:        m.Name,
		Bin:         m.Bin,
		Version:     m.Version,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := saveState(records); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	output.Successf("Installed %s to %s", m.Name, binPath)
	output.Dim("  Platform: %s", platformKey)
	output.Dim("  Version: %s", m.Version)
	warnPath(dir)
	return nil
}

func parseRepo(repo string) (owner, name, tag string, err error) {
	repo = strings.TrimSpace(repo)
	if strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "http://") {
		return "", "", "", fmt.Errorf("use github.com/user/repo format, not a full URL")
	}

	tag = ""
	if idx := strings.Index(repo, "@"); idx != -1 {
		tag = repo[idx+1:]
		repo = repo[:idx]
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 3 || parts[0] != "github.com" {
		return "", "", "", fmt.Errorf("expected github.com/user/repo, got %q", repo)
	}
	if parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("owner and repo name cannot be empty")
	}
	return parts[1], parts[2], tag, nil
}

func fetchManifest(owner, name, tag string) (*manifest.Manifest, error) {
	ref := tag
	if ref == "" {
		ref = "main"
	}
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/mira.json", owner, name, ref)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: %s", url, resp.Status)
	}

	var m manifest.Manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("decode mira.json: %w", err)
	}
	return &m, nil
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func verifyChecksum(data []byte, expectedHex string) error {
	if expectedHex == "" {
		return fmt.Errorf("no sha256 provided in manifest")
	}
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if !strings.EqualFold(got, expectedHex) {
		return fmt.Errorf("expected %s, got %s", expectedHex, got)
	}
	return nil
}

func humanSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func supportedList(platforms map[string]manifest.Platform) string {
	keys := make([]string, 0, len(platforms))
	for k := range platforms {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func Remove(repo string) error {
	output.Header(repo, "uninstall")

	records, err := loadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	rec, ok := records[repo]
	if !ok {
		return fmt.Errorf("%s is not installed", repo)
	}

	binPath := filepath.Join(binDir(), rec.Bin)
	if err := os.Remove(binPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", binPath, err)
	}

	delete(records, repo)
	if err := saveState(records); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	output.Successf("Removed %s (%s)", rec.Bin, rec.Repo)
	output.Dim("  Binary deleted: %s", binPath)
	return nil
}

func List() error {
	output.Header("", "list")

	records, err := loadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if len(records) == 0 {
		output.Warning("No packages installed")
		return nil
	}

	repos := make([]string, 0, len(records))
	for r := range records {
		repos = append(repos, r)
	}
	sort.Strings(repos)

	rows := make([][]string, 0, len(repos))
	for _, r := range repos {
		rec := records[r]
		rows = append(rows, []string{rec.Repo, rec.Name, rec.Version, rec.InstalledAt})
	}

	output.Table([]string{"Package", "Name", "Version", "Installed At"}, rows)
	return nil
}

func Update(repo string) error {
	records, err := loadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	key := findStateKey(records, repo)
	if key == "" {
		return fmt.Errorf("%s is not installed", repo)
	}

	output.Header(key, "update")
	return Install(key)
}

func Info(repo string) error {
	records, err := loadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	key := findStateKey(records, repo)
	if key == "" {
		return fmt.Errorf("%s is not installed", repo)
	}

	rec := records[key]

	owner, name, tag, err := parseRepo(key)
	if err != nil {
		return err
	}

	m, err := fetchManifest(owner, name, tag)
	if err != nil {
		output.Warningf("Could not fetch manifest: %v", err)
	}

	output.Header(rec.Repo, rec.Version)

	output.Subtitle("Installed")
	output.KeyValue("Name", rec.Name)
	output.KeyValue("Binary", rec.Bin)
	output.KeyValue("Installed", rec.InstalledAt)

	if m != nil {
		output.Subtitle("Manifest")
		output.KeyValue("Version", m.Version)
		output.KeyValue("Description", m.Description)
		output.KeyValue("Author", m.Author)
		output.KeyValue("Repository", m.Repo)

		output.Subtitle("Platforms")
		rows := make([][]string, 0, len(m.Platforms))
		for p := range m.Platforms {
			mark := " "
			if p == fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) {
				mark = "●"
			}
			rows = append(rows, []string{mark, p})
		}
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][1] < rows[j][1]
		})
		output.Table([]string{"", "Platform"}, rows)
		output.Dim("  ● = current system")
	}

	return nil
}

func findStateKey(records map[string]Record, repo string) string {
	if _, ok := records[repo]; ok {
		return repo
	}
	prefix := repo + "@"
	for k := range records {
		if strings.HasPrefix(k, prefix) {
			return k
		}
	}
	return ""
}

func binDir() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "bin")
}

func warnPath(dir string) {
	path := os.Getenv("PATH")
	if !strings.Contains(path, dir) {
		output.Warningf("%s is not in your PATH", dir)
		output.Dim("  Add this to your shell config:")
		output.Dim("  export PATH=\"%s:$PATH\"", dir)
	}
}
