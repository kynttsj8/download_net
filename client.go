package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	var key string
	for {
		// Authentication
		fmt.Print("Enter username: ")
		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(c, username)

		fmt.Print("Enter password: ")
		password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(c, password)

		// Read the response from the server
		response, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("Server connection: " + response)
		c.Write([]byte("Client Connected\n"))
		if strings.Contains(response, "successful") {
			key, err = bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("You are now connected to the server. Your key is %s \n", key)
			break
		}

	}

	for {
		fmt.Println("1. Guess the number")
		fmt.Println("2. Download a file")
		fmt.Println("Enter your choice: ")
		choice, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Guess the number: ")
			text, _ := bufio.NewReader(os.Stdin).ReadString('\n') // Read input from the user
			fmt.Fprintf(c, text)                                  // Send the input to the server

			message, _ := bufio.NewReader(c).ReadString('\n') // Read the response from the server
			key = strings.TrimSpace(key)
			fmt.Print("Server: " + message)
			if strings.TrimSpace(message) == key+"_"+"Correct" {
				fmt.Println("Congratulations! You've guessed the number correctly.\nAutomatically play with another guessing number? (PRESS -1 to end the game)")
			} else if strings.TrimSpace(message) == key+"_"+"GameOver" {
				fmt.Println("Game Over! You've quit the game.")
				break
			}
		case "2":
			// Request file download
			fmt.Print("Enter the file name to download: ")
			filename, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Fprintf(c, filename)

			// Receive file content
			fileContent, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			// Remove newline characters from the file content
			fileContent = strings.TrimSpace(fileContent)

			// Prompt the user for the destination file name
			fmt.Print("Enter the destination file name: ")
			destFileName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			destFileName = strings.TrimSpace(destFileName)

			// Create the destination file
			destFile, err := os.Create(destFileName)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer destFile.Close()

			// Write the file content to the destination file
			_, err = destFile.WriteString(fileContent)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("File '%s' downloaded successfully!\n", destFileName)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
