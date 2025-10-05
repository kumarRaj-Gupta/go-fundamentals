package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Goal: Learn to read a file line-by-line efficiently, a skill critical
// for handling your JSON-RPC server's STDIN stream.
func main() {
	fmt.Println("--- Project 2: File and Stream I/O (`os` & `bufio`) ---")

	// 1. Open the file. This returns the file handler and an error.
	file, err := os.Open("data.txt")
	if err != nil {
		// If the file doesn't exist or permissions are wrong, we exit.
		log.Fatalf("Error opening file: %v", err)
	}

	// 2. IMPORTANT: Defer closing the file.
	// The 'defer' keyword ensures file.Close() is called when the surrounding
	// function (main) exits, regardless of whether it exits normally or due to a panic.
	defer file.Close()

	// 3. Create a new scanner to read the file stream.
	// bufio.NewScanner takes an io.Reader (like our opened file).
	/*
		How would you modify this function to read from STDIN, knowing that os.Stdin is an already-opened, globally available *os.File that implements the io.Reader interface required by bufio.NewScanner?
	*/
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	lineCount := 0

	// fmt.Println("\n[Reading File Content Line by Line]")

	// 4. Start the reading loop.
	// scanner.Scan() advances the scanner to the next token (by default, a line).
	// It returns false when the input is exhausted or an error occurs.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Line %d: %v\n", lineCount+1, line)
		var response string = "ACK: packet received"
		fmt.Fprintf(writer, "%v\n", response)
		lineCount++
		err := writer.Flush()
		if err != nil {
			log.Fatalf("Error flushing writer: %v", err)
		}
	}

	// 5. Check for errors that occurred during scanning.
	// This is crucial because the loop only exits on exhausted input OR an error.
	// So this is for when the exit happens due to an error. GOT IT.
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during file scanning: %v", err)
	}

	fmt.Printf("\nSuccessfully read %d lines.\n", lineCount)
}
