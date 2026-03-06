package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github-activity/internal/usecase"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func main() {

	programName := os.Args[0]

	// Determine the output file path, defaulting to the user's config directory if possible
	var outputFilePath string
	configFilePath, err := os.UserConfigDir()
	if err != nil {
		color.Red("Error getting user config directory: " + err.Error())
		outputFilePath = "github-activity.json"
		color.Red("Falling back to current directory: " + outputFilePath)
	} else {
		appDir := filepath.Join(configFilePath, "github-activity")
		err = os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			color.Red("Error creating app directory: " + err.Error())
			outputFilePath = "github-activity.json"
			color.Red("Falling back to current directory: " + outputFilePath)
		} else {
			outputFilePath = filepath.Join(appDir, "github-activity.json")
			color.Green("Output file path: " + outputFilePath)
		}
	}

	// Variables for command-line arguments
	userPtr := flag.String("username", "jkligel", "GitHub username to fetch activity for")

	ghTokenPtr := flag.String("token", "", "GitHub API token")

	getNewestPtr := flag.Bool("newest", false, "Only print the newest event instead of the entire activity feed")

	// Parse command-line arguments
	flag.Parse()

	urlStr := fmt.Sprintf("https://api.github.com/users/%s/events", *userPtr)

	if *ghTokenPtr == "" {
		if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
			*ghTokenPtr = envToken
		} else {
			color.Red("Please provide a GitHub API token using the -token flag or by setting the GITHUB_TOKEN environment variable.")
			os.Exit(1)
		}
	}

	authorizationString := fmt.Sprintf("bearer %s", *ghTokenPtr)

	// Print output
	color.Blue("Program name: " + programName)
	color.Blue("GitHub username: " + *userPtr)
	color.Blue("API URL: " + urlStr)

	// Make the API request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		color.Red("Error creating request: " + err.Error())
		os.Exit(1)
	}
	req.Header.Set("Authorization", authorizationString)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error making request: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		color.Red(fmt.Sprintf("Error: received status code %d", resp.StatusCode))
		os.Exit(1)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Error reading response body: " + err.Error())
		os.Exit(1)
	}

	// Process the response body in separate goroutines and write to file if there are changes
	var wg sync.WaitGroup

	// Goroutine 1: Write to file
	wg.Add(1)
	go func() {
		defer wg.Done()
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			color.Red("Error marshaling JSON: " + err.Error())
			return
		}
		fileBytes, err := os.ReadFile(outputFilePath)
		if err != nil {
			color.Red("Nothing was written to " + outputFilePath)
			color.Red("Error reading file: " + err.Error())
			color.Green("Writing to file.")
			err = os.WriteFile(outputFilePath, prettyJSON.Bytes(), os.FileMode(os.O_CREATE))
			if err != nil {
				color.Red("Error writing to file: " + err.Error())
				return
			}
			return
		}
		if bytes.Equal(fileBytes, prettyJSON.Bytes()) {
			color.Blue("No changes detected, not writing to file.")
			return
		}
		color.Green("Changes detected, writing to file.")
		err = os.WriteFile(outputFilePath, prettyJSON.Bytes(), os.FileMode(os.O_CREATE))
		if err != nil {
			color.Red("Error writing to file: " + err.Error())
			return
		}
	}()

	// Goroutine 2: Parse and print output
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Unmarshal the JSON response into a slice of maps and parse the actions into a slice of strings
		var data []map[string]any
		var githubActivityOutput []string
		err = json.Unmarshal(body, &data)
		if err != nil {
			color.Red("Error unmarshaling JSON: " + err.Error())
			os.Exit(1)
		}
		githubActivityOutput = usecase.ParseActions(data)

		// Print the output to the console, either the entire activity feed or just the newest event

		color.Magenta("Output:")

		if *getNewestPtr {
			fmt.Println(githubActivityOutput[0])
		} else {
			fmt.Println(strings.Join(githubActivityOutput, "\n"))
		}
	}()

	// Wait for the goroutine to complete
	wg.Wait()
}
