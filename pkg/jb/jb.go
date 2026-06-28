package jb

import (
	"fmt"
	"io"
	"iter"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/apalala/jb/pkg/blue"
	"github.com/apalala/jb/pkg/bmx"
)

const (
	JBHeader = "# Johannes Blues - A view into great literary works"

	StreamTime = 5 * time.Second

	GutenbergURL = "https://www.gutenberg.org/cache/epub/%[1]s/pg%[1]s.txt"
)

type Work struct {
	ID         string
	Type       string
	WindowSize int
}

var (
	WorksDatabase = []Work{
		{ID: "pg1524", Type: "T", WindowSize: 15},
		{ID: "pg2701", Type: "N", WindowSize: 25},
		{ID: "pg1508", Type: "T", WindowSize: 15},
		{ID: "pg84", Type: "N", WindowSize: 25},
	}

	TheatreCleaningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?m)^[A-Z0-9_\s]{2,15}[.:]\s*`),
		regexp.MustCompile(`(?m)[\[(].*?[\])]`),
	}

	NovelCleaningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^(CHAPTER|C_H_A_P_T_E_R)\s+[IVX0-9]+.*`),
	}

	CleaningPatterns = map[string][]*regexp.Regexp{
		"T": TheatreCleaningPatterns,
		"N": NovelCleaningPatterns,
	}

	wordRE   = regexp.MustCompile(`\w`)
	letterRE = regexp.MustCompile(`[a-zA-Z]`)
)

func extractID(workID string) string {
	return strings.TrimPrefix(workID, "pg")
}

func findWorkFile(workID string, suffix string) (string, bool) {
	dirs := []string{"works", filepath.Join("jb", "works"), filepath.Join("..", "works")}
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) && strings.Contains(entry.Name(), workID) {
				return filepath.Join(dir, entry.Name()), true
			}
		}
	}
	return "", false
}

func LoadWork(workID string) (string, error) {
	gid := extractID(workID)

	if path, ok := findWorkFile(workID, ".txt.bmx"); ok {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("jb: read %s: %w", path, err)
		}
		unsealed, err := bmx.UnsealText(string(data))
		if err != nil {
			return "", fmt.Errorf("jb: unseal %s: %w", path, err)
		}
		return unsealed, nil
	}

	if path, ok := findWorkFile(workID, ".txt"); ok {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("jb: read %s: %w", path, err)
		}
		return string(data), nil
	}

	url := fmt.Sprintf(GutenbergURL, gid)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("jb: fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("jb: read %s: %w", url, err)
	}
	return string(body), nil
}

func CleanWork(workType string, text string) string {
	pats := CleaningPatterns[workType]
	return strings.Join(ParseVerses(text, pats), "\n")
}

func ParseVerses(rawText string, patterns []*regexp.Regexp) []string {
	var lines []string
	for line := range strings.Lines(rawText) {
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSpace(line)
		if line != "" && wordRE.MatchString(line) {
			lines = append(lines, line)
		}
	}

	startIdx := 0
	for i, line := range lines {
		if i > 500 {
			break
		}
		upper := strings.ToUpper(line)
		if strings.Contains(upper, "START OF THE PROJECT") || strings.Contains(upper, "START OF THIS PROJECT") {
			startIdx = i + 1
			break
		}
	}

	endIdx := len(lines)
	for i, line := range lines {
		if len(lines)-i > 1000 {
			continue
		}
		upper := strings.ToUpper(line)
		if strings.Contains(upper, "END OF THE PROJECT") || strings.Contains(upper, "END OF THIS PROJECT") {
			endIdx = i
			break
		}
	}

	body := lines[startIdx:endIdx]
	var cleaned []string

	for _, line := range body {
		result := line
		for _, pat := range patterns {
			result = pat.ReplaceAllString(result, "")
		}
		result = strings.TrimSpace(result)
		if result != "" && letterRE.MatchString(result) {
			cleaned = append(cleaned, result)
		}
	}

	return cleaned
}

func StreamBlueVerses(verses []string, windowSize int) iter.Seq[string] {
	return func(yield func(string) bool) {
		if len(verses) == 0 {
			return
		}
		if len(verses) == 1 {
			for {
				if !yield(verses[0]) {
					return
				}
			}
		}

		currentIdx := rand.IntN(len(verses))
		recent := make(map[string]struct{})
		recentWindow := min(windowSize, len(verses)-1)
		if recentWindow < 1 {
			recentWindow = 1
		}
		recentOrder := make([]string, 0, recentWindow)

		for range recentWindow {
			recent[""] = struct{}{}
			recentOrder = append(recentOrder, "")
		}

		for v := range blue.StreamBlueSignal(0.65) {
			jump := int(v * float64(windowSize))
			currentIdx = mod(currentIdx+jump, len(verses))
			verse := verses[currentIdx]

			if _, ok := recent[verse]; ok {
				continue
			}

			delete(recent, recentOrder[0])
			recentOrder = recentOrder[1:]

			recent[verse] = struct{}{}
			recentOrder = append(recentOrder, verse)

			if !yield(verse) {
				return
			}
		}
	}
}

func PrintWork(text string, windowSize int) {
	verses := strings.Split(text, "\n")
	var nonEmpty []string
	for _, v := range verses {
		if v != "" {
			nonEmpty = append(nonEmpty, v)
		}
	}
	if len(nonEmpty) == 0 {
		return
	}

	out := os.Stdout
	if stat, _ := os.Stdout.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
		out = os.Stderr
	}

	fmt.Println(JBHeader)
	deadline := time.Now().Add(StreamTime)
	for verse := range StreamBlueVerses(nonEmpty, windowSize) {
		if time.Now().After(deadline) {
			break
		}
		fmt.Fprintln(out, verse)
	}
}

func mod(a, n int) int {
	r := a % n
	if r < 0 {
		r += n
	}
	return r
}

func Main() int {
	work := WorksDatabase[rand.IntN(len(WorksDatabase))]
	raw, err := LoadWork(work.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "jb: %v\n", err)
		return 1
	}
	text := CleanWork(work.Type, raw)
	PrintWork(text, work.WindowSize)
	return 0
}
