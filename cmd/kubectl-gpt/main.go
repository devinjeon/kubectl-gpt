package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	gpt "github.com/devinjeon/kubectl-gpt/pkg/gpt"
)

const (
	DefaultMaxTokens   = 300
	DefaultTemperature = 0.2
	DefaultModel       = "gpt-3.5-turbo"
	DefaultApiUrl      = "https://api.openai.com/v1/chat/completions"
	systemMessage      = "Translate the given text to a kubectl command. Show only generated kubectl command without any description, code block."
)

var (
	version     string
	yesFlag     = flag.Bool("yes", false, "Execute the generated command without asking for confirmation")
	noFlag      = flag.Bool("no", false, "Print the generated command without executing it")
	helpFlag    = flag.Bool("help", false, "Show usage information")
	versionFlag = flag.Bool("version", false, "Show the version")
)

func main() {
	// Parse flags
	flag.Parse()

	// Handle help and version flags
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	apiUrl, apiKey, model, temperature, maxTokens := getOpenAIConfigFromEnv()
	if apiKey == "" {
		fmt.Println("Please set the environment variable: \"OPENAI_API_KEY\".")
		fmt.Println("You can add the following line to your .zshrc or .bashrc file:")
		fmt.Println("export OPENAI_API_KEY=<your-key>")
		fmt.Println()
		fmt.Println("If you don't have an OpenAI API Key, you can get one at this link: https://platform.openai.com/account/api-keys.")
		os.Exit(1)
	}

	query := strings.Join(flag.Args(), " ")
	query = strings.TrimSpace(query)
	if query == "" {
		fmt.Println("Please input a query.")
		fmt.Println("Usage: kubectl-gpt [OPTIONS] QUERY")
		os.Exit(1)
	}
	request := gpt.NewOpenAIRequest(model, temperature, maxTokens, systemMessage, query)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go printLoadingMessage(&wg)

	response, err := gpt.RequestChatGptAPI(apiUrl, request, apiKey)

	wg.Done()
	fmt.Printf("\033[2K\r") // Clear loading message after completion
	if err != nil {
		fmt.Printf("Failed to call OpenAI API at %s:\n%v\n", apiUrl, err)
		os.Exit(1)
	}

	kubectlCommand := extractCommand(response)
	fmt.Println("\033[1;31m❗[WARNING] Please verify the generated commands before executing them on your k8s cluster,",
		"especially `update` and `patch` commands, as GPT-generated commands may be inaccurate.\033[0m")
	fmt.Printf("\033[1;34m[Generated Command]\033[0m\n%s\n", kubectlCommand)

	if *noFlag {
		os.Exit(0)
	}
	if !*yesFlag {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\033[1;34m⎈ Do you really want to execute this command?\033[0m [y/N]: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.ToLower(text) != "y" {
			os.Exit(0)
		}
	}

	cmd := exec.Command("/bin/sh", "-c", kubectlCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// getOpenAIConfigFromEnv retrieves OpenAI related variables from environment
func getOpenAIConfigFromEnv() (apiUrl, apiKey string, model string, temperature float64, maxTokens int) {
	apiUrl = os.Getenv("OPENAI_API_URL")
	if apiUrl == "" {
		apiUrl = DefaultApiUrl
	}

	apiKey = os.Getenv("OPENAI_API_KEY")

	model = os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = DefaultModel
	}

	temperatureStr := os.Getenv("OPENAI_TEMPERATURE")
	temperature = DefaultTemperature
	if temperatureStr != "" {
		var err error
		temperature, err = strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			fmt.Println("Failed to parse OPENAI_TEMPERATURE. Using default temperature.")
			temperature = DefaultTemperature
		}
	}

	maxTokensStr := os.Getenv("OPENAI_MAX_TOKENS")
	maxTokens = DefaultMaxTokens
	if maxTokensStr != "" {
		var err error
		maxTokens, err = strconv.Atoi(maxTokensStr)
		if err != nil {
			fmt.Println("Failed to parse OPENAI_MAX_TOKENS. Using default max tokens.")
			maxTokens = DefaultMaxTokens
		}
	}

	return apiUrl, apiKey, model, temperature, maxTokens
}

func printLoadingMessage(wg *sync.WaitGroup) {
	loading := "Getting kubectl command from GPT API "
	i := 0

	loadingSymbol := "🚶"

	loadingTicker := time.NewTicker(time.Millisecond * 500)

	go func() {
		for {
			select {
			case <-loadingTicker.C:
				if i%2 == 0 {
					loadingSymbol = "🚶"
				} else {
					loadingSymbol = "🏃"
				}
				fmt.Printf("\033[2K\r%s %s%s", loadingSymbol, loading, strings.Repeat(".", i%30))
				i++
			}
		}
	}()

	wg.Wait()
	loadingTicker.Stop()
}

func printHelp() {
	fmt.Println("Usage: kubectl-gpt [OPTIONS] QUERY")
	fmt.Println("Translate the given query to a kubectl command using OpenAI GPT API.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --yes           Execute the generated command without asking for confirmation")
	fmt.Println("  --no            Print the generated command without executing it")
	fmt.Println("  --help          Show this message and exit")
	fmt.Println("  --version       Show the version")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  OPENAI_API_URL        OpenAI API URL (default is https://api.openai.com/v1/chat/completions)")
	fmt.Println("  OPENAI_API_KEY        OpenAI API Key")
	fmt.Println("  OPENAI_MODEL          OpenAI Model to use (default is gpt-3.5-turbo)")
	fmt.Println("  OPENAI_TEMPERATURE    Temperature for the OpenAI request (default is 0.2)")
	fmt.Println("  OPENAI_MAX_TOKENS     Max tokens for the OpenAI request (default is 300)")
}

func extractCommand(response gpt.OpenAIResponse) string {
	s := strings.Trim(response.Choices[0].Message.Content, "`")
	return strings.TrimSpace(s)
}
