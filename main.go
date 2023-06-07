package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/exec"
  "sort"
  "strings"
)

type Certificate struct {
  NameValue string `json:"name_value"`
}

func main() {
  // Parse command-line arguments
  domain := flag.String("d", "", "Domain to query")
  outputFile := flag.String("o", "", "Output file")
  flag.Parse()

  if *domain == "" || *outputFile == "" {
    fmt.Println("Both -d (domain) and -o (output) arguments are required")
    flag.PrintDefaults()
    os.Exit(1)
  }

  // Send GET request to the specified URL
  url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", *domain)
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  // Read the response body
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Parse the JSON response
  var certificates []Certificate
  err = json.Unmarshal(body, &certificates)
  if err != nil {
    log.Fatal(err)
  }

  // Extract the name_value keys
  var result []string
  for _, cert := range certificates {
    result = append(result, cert.NameValue)
  }

  // Remove duplicates from the result
  result = uniqueItems(result)

  // Sort the result
  sort.Strings(result)

  // Remove asterisks using the sed command
  sedCmd := exec.Command("sed", "s/\\*//g")
  sedCmd.Stdin = strings.NewReader(strings.Join(result, "\n"))
  sedOutput, err := sedCmd.Output()
  if err != nil {
    log.Fatal(err)
  }

  // Write the result to the output file
  err = ioutil.WriteFile(*outputFile, sedOutput, 0644)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Results written to %s\n", *outputFile)
}

// Returns a slice with unique items
func uniqueItems(items []string) []string {
  encountered := map[string]bool{}
  uniqueItems := []string{}

  for _, item := range items {
    if !encountered[item] {
      encountered[item] = true
      uniqueItems = append(uniqueItems, item)
    }
  }

  return uniqueItems
}
