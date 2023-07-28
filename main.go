package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/eiannone/keyboard"
)

func getInputLines() []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\033[32mPaste your JSON text and press Enter.")
	fmt.Println("To convert and get the output, type '.' on a single line and press Enter.\033[0m")
	fmt.Println("--------------------------------------------------------------------------------")

	var inputLines []string

	for {
		// Read the user's input line-by-line
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		// Check if the user wants to exit (single ".")
		if strings.TrimSpace(input) == "." {
			break
		}

		// Add the input to the slice
		inputLines = append(inputLines, input)
	}

	return inputLines
}

func convertToEscapedJSON(inputLines []string) (string, error) {
	// Combine all lines into a single string
	input := strings.Join(inputLines, "")
	input = strings.TrimSpace(input)
	// Check if the user wants to exit
	if input == "." {
		return "", fmt.Errorf("exiting the program")
	}

	// Convert JSON to escaped JSON format
	escapedJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("error converting JSON to escaped JSON: %v", err)
	}

	// Remove the first and last " characters
	escapedJSONString := string(escapedJSON)[1 : len(escapedJSON)-1]

	return escapedJSONString, nil
}

func convertToUnescapedJSON(inputLines []string) (string, error) {
	// Combine all lines into a single string
	input := strings.Join(inputLines, "")

	// Trim the input
	input = strings.TrimSpace(input)
	if len(input) >= 1 {
		if input[0] != '"' || input[len(input)-1] != '"' {
			input = fmt.Sprintf(`"%s"`, input)
		}

	}

	// Check if the user wants to exit
	if input == "." {
		return "", fmt.Errorf("exiting the program")
	}

	// Convert escaped JSON to regular JSON format
	var unescapedJSON string
	err := json.Unmarshal([]byte(input), &unescapedJSON)
	if err != nil {
		return "", fmt.Errorf("error converting escaped JSON to regular JSON: %v", err)
	}

	return unescapedJSON, nil
}

func copyToClipboard(JSONString string) error {
	// Send the escaped JSON to the clipboard
	if err := clipboard.WriteAll(JSONString); err != nil {
		return fmt.Errorf("error sending the JSON to clipboard: %v", err)
	}

	return nil
}

func escapeJson() {
	inputLines := getInputLines()

	escapedJSONString, err := convertToEscapedJSON(inputLines)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\033[35mEscaped JSON:\033[0m")
	fmt.Println(escapedJSONString)

	if err := copyToClipboard(escapedJSONString); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\033[32mEscaped JSON has been copied to the clipboard.\033[0m")
}

func unescapeJson() {
	inputLines := getInputLines()

	unescapedJSONString, err := convertToUnescapedJSON(inputLines)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\033[35mUnescaped JSON:\033[0m")
	fmt.Println(unescapedJSONString)

	if err := copyToClipboard(unescapedJSONString); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\033[32mUnescaped JSON has been copied to the clipboard.\033[0m")
}

func maxLength(strings []string) int {
	max := 0

	for _, str := range strings {
		length := len(str)
		if length > max {
			max = length
		}
	}

	return max
}

func printOptions() {
	options := []string{
		"Press 1 to escape JSON",
		"Press 2 to unescape (regular) JSON",
		"Press ESC or 'q' to quit",
	}
	width := maxLength(options) + 1
	height := len(options) + 1

	// Top line of the box
	fmt.Print("+")
	for i := 0; i <= width; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Lines inside the box
	for j := 0; j <= height-2; j++ {
		if j == 2 {
			fmt.Print("| \033[34m", options[j])
		} else if j == 1 {
			fmt.Print("| \033[36m", options[j])
		} else {
			fmt.Print("| \033[34m", options[j])
		}
		for i := 0; i <= width-len(options[j])-1; i++ {
			fmt.Print(" ")
		}
		fmt.Println("\033[0m|")
	}

	// Bottom line of the box
	fmt.Print("+")
	for i := 0; i <= width; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func main() {
	exit := make(chan struct{})

	printOptions()

	defer func() {
		_ = keyboard.Close()
	}()

	go func() {
		defer close(exit)

		for {
			rn, kb, keyErr := keyboard.GetSingleKey()
			if keyErr != nil {
				fmt.Println("Error reading key:", keyErr)
				return
			}
			if kb == keyboard.KeyEsc || rn == 'q' {
				return
			}
			if rn == '1' {
				escapeJson()
				printOptions()
			}
			if rn == '2' {
				unescapeJson()
				printOptions()

			}
		}
	}()

	<-exit
}
