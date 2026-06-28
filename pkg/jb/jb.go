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
	StreamTime = 5 * time.Second

	HamletURL   = "https://www.gutenberg.org/cache/epub/1524/pg1524.txt"
	MobyDickURL = "https://www.gutenberg.org/cache/epub/2701/pg2701.txt"
)

var (
	WorksDatabase = map[string]string{
		"hamlet":    filepath.FromSlash("works/pg1524.txt.bmx"),
		"mobi_dick": filepath.FromSlash("works/pg2701.txt.bmx"),
	}

	TheatreCleaningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?m)^[A-Z0-9_\s]{2,15}[.:]\s*`),
		regexp.MustCompile(`(?m)[\[(].*?[\])]`),
	}

	NovelCleaningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^(CHAPTER|C_H_A_P_T_E_R)\s+[IVX0-9]+.*`),
	}
)

func FetchAndParseVerses(url string, patterns []*regexp.Regexp) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("jb: fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("jb: read %s: %w", url, err)
	}

	return ParseVerses(string(body), patterns), nil
}

func LoadWork(name string, patterns []*regexp.Regexp) ([]string, error) {
	rel, ok := WorksDatabase[name]
	if !ok {
		return nil, fmt.Errorf("jb: unknown work %q", name)
	}

	pathsToTry := []string{
		rel,
		filepath.Join("jb", rel),
		filepath.Join("..", rel),
	}

	var data []byte
	var errs []string
	for _, p := range pathsToTry {
		d, e := os.ReadFile(p)
		if e == nil {
			data = d
			break
		}
		errs = append(errs, fmt.Sprintf("%s: %v", p, e))
	}
	if data == nil {
		return nil, fmt.Errorf("jb: load work %q: %s", name, strings.Join(errs, "; "))
	}

	unsealed, err := bmx.UnsealText(string(data))
	if err != nil {
		return nil, fmt.Errorf("jb: unseal %s: %w", name, err)
	}

	return ParseVerses(unsealed, patterns), nil
}

func ParseVerses(rawText string, patterns []*regexp.Regexp) []string {
	rawLines := strings.Split(rawText, "\n")

	var lines []string
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line != "" {
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
			endIdx = len(lines) - (len(lines) - i)
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
		if result != "" && regexp.MustCompile(`[a-zA-Z]`).MatchString(result) {
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

func PrintHamletVerses(dur time.Duration) {
	lines, err := LoadWork("hamlet", TheatreCleaningPatterns)
	if err != nil {
		var fetchErr error
		lines, fetchErr = FetchAndParseVerses(HamletURL, TheatreCleaningPatterns)
		if fetchErr != nil {
			fmt.Fprintf(os.Stderr, "jb: load hamlet: %v (fetch: %v)\n", err, fetchErr)
			return
		}
	}

	deadline := time.Now().Add(dur)
	for verse := range StreamBlueVerses(lines, 15) {
		if time.Now().After(deadline) {
			break
		}
		fmt.Fprintln(os.Stdout, verse)
	}
}

func PrintMobyDickVerses(dur time.Duration) {
	lines, err := LoadWork("mobi_dick", NovelCleaningPatterns)
	if err != nil {
		var fetchErr error
		lines, fetchErr = FetchAndParseVerses(MobyDickURL, NovelCleaningPatterns)
		if fetchErr != nil {
			fmt.Fprintf(os.Stderr, "jb: load moby_dick: %v (fetch: %v)\n", err, fetchErr)
			return
		}
	}

	deadline := time.Now().Add(dur)
	for verse := range StreamBlueVerses(lines, 25) {
		if time.Now().After(deadline) {
			break
		}
		fmt.Fprintln(os.Stdout, verse)
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
	PrintHamletVerses(StreamTime)
	return 0
}
