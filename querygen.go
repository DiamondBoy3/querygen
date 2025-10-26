package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rix4uni/querygen/banner"
	"github.com/spf13/pflag"
)

const (
	baseURL = "https://raw.githubusercontent.com/rix4uni/nucleihubquery/refs/heads/main"
)

var (
	engine   = pflag.StringP("engine", "e", "shodan", "Search engine: shodan, google, censys, fofa, hunter, zoomeye")
	severity = pflag.StringP("severity", "s", "all", "Severity level: all, critical, high, medium, low, info (comma-separated)")
	silent   = pflag.Bool("silent", false, "Silent mode.")
	version  = pflag.Bool("version", false, "Print the version of the tool and exit.")
)

type EngineConfig struct {
	Name     string
	Path     string
	Prefix   string
	UseGrepA bool // For binary files that need grep -a
}

func getEngineConfigs(severities []string) []EngineConfig {
	var configs []EngineConfig
	
	// Define all available engines with their prefixes and grep requirements
	engines := map[string]struct {
		prefix  string
		useGrepA bool
	}{
		"google":  {"intitle:", true},
		"censys":  {"services.http.response.html_title=", false},
		"shodan":  {"http.title:", false},
		"fofa":    {"title=", false},
		"hunter":  {"web.title=", false},
		"zoomeye": {"title=", false},
	}

	// Generate configs for each severity and each engine
	for _, severity := range severities {
		for engineName, engineConfig := range engines {
			configs = append(configs, EngineConfig{
				Name:     engineName,
				Path:     fmt.Sprintf("%s/%s-query-%s.txt", engineName, engineName, severity),
				Prefix:   engineConfig.prefix,
				UseGrepA: engineConfig.useGrepA,
			})
		}
	}
	
	return configs
}

func main() {
	pflag.Parse()

    if *version {
        banner.PrintBanner()
        banner.PrintVersion()
        os.Exit(0)
    }

    if !*silent {
        banner.PrintBanner()
    }

	// Validate engine
	validEngines := map[string]bool{
		"shodan": true, "google": true, "censys": true, 
		"fofa": true, "hunter": true, "zoomeye": true,
	}
	if !validEngines[*engine] {
		fmt.Printf("Error: engine must be one of: shodan, google, censys, fofa, hunter, zoomeye\n")
		pflag.Usage()
		os.Exit(1)
	}

	// Parse and validate severities
	severities := parseSeverities(*severity)
	if len(severities) == 0 {
		fmt.Printf("Error: severity must be one or more of: all, critical, high, medium, low, info\n")
		pflag.Usage()
		os.Exit(1)
	}

	targetFormat := getTargetFormat(*engine)
	engines := getEngineConfigs(severities)
	
	// Generate queries for ALL engines (like original code)
	for _, engine := range engines {
		if err := generateQueries(engine, targetFormat); err != nil {
			log.Printf("Error generating %s queries: %v", engine.Name, err)
		}
	}
}

func parseSeverities(severityStr string) []string {
	validSeverities := map[string]bool{
		"all": true, "critical": true, "high": true, 
		"medium": true, "low": true, "info": true,
	}
	
	var result []string
	parts := strings.Split(severityStr, ",")
	
	for _, part := range parts {
		severity := strings.TrimSpace(part)
		if validSeverities[severity] {
			result = append(result, severity)
		} else {
			fmt.Printf("Warning: invalid severity '%s', skipping\n", severity)
		}
	}
	
	return result
}

func getTargetFormat(engine string) string {
	switch engine {
	case "shodan":
		return "http.title:"
	case "google":
		return "intitle:"
	case "censys":
		return "services.http.response.html_title="
	case "fofa":
		return "title="
	case "hunter":
		return "web.title="
	case "zoomeye":
		return "title="
	default:
		return "http.title:" // default to shodan format
	}
}

func generateQueries(engine EngineConfig, targetFormat string) error {
	// Fetch the data from GitHub
	url := fmt.Sprintf("%s/%s", baseURL, engine.Path)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s data: %v", engine.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Skip if file doesn't exist for this severity
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("failed to fetch %s data: HTTP %d", engine.Name, resp.StatusCode)
	}

	// Process the data
	if engine.UseGrepA {
		// Read as raw bytes to handle binary data
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read %s data: %v", engine.Name, err)
		}
		
		// Convert to string and split by lines
		content := string(body)
		lines := strings.Split(content, "\n")
		
		for _, line := range lines {
			processLine(line, engine.Prefix, targetFormat)
		}
	} else {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			processLine(line, engine.Prefix, targetFormat)
		}
		
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading %s response: %v", engine.Name, err)
		}
	}

	return nil
}

func processLine(line, sourcePrefix, targetFormat string) {
	// Filter lines that start with the source prefix and convert to target format
	if strings.HasPrefix(line, sourcePrefix) {
		converted := strings.Replace(line, sourcePrefix, targetFormat, 1)
		fmt.Println(converted)
	}
}