package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func sanitizeString(dirtyString string, maxlen int) string {
	var re = regexp.MustCompile(`(?i)[^a-z0-9_\.]`)
	if len(re.ReplaceAllString(dirtyString, "")) > maxlen {
		return re.ReplaceAllString(dirtyString, "")[:maxlen]
	}
	return re.ReplaceAllString(dirtyString, "")
}

func announce(message string) {
	fmt.Printf(
		"\n" + nc + green + "---------------------------------------------\n   " + brightcyan + "+ " + message + green + "\n---------------------------------------------\n\n" + nc)
}

// fileExists Does this file exist?
func fileExists(filename string) bool {
	referencedFile, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !referencedFile.IsDir()
}

// fileContainsString This is a utility to see if a string in a file exists.
func fileContainsString(str, filepath string) bool {
	accused, _ := os.ReadFile(filepath)
	isExist, _ := regexp.Match(str, accused)
	return isExist
}

func timeStamp() string {
	current := time.Now()
	return current.Format("2006-01-02 15:04:05")
}

func unixTimeStampNano() string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return timestamp
}

func createDirIfItDontExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		handle("Could not create directory: ", err)
	}
}

// handle Ye Olde Error Handler takes a message and an error code
func handle(msg string, err error) {
	if err != nil {
		fmt.Printf(brightred+"\n%s: %s"+white, msg, err)
	}
}

// createFile Generic file handler
func createFile(filename string) {
	var _, err = os.Stat(filename)
	if os.IsNotExist(err) {
		var file, err = os.Create(filename)
		handle("", err)
		defer file.Close()
	}
}

func appendFile(filename, textToWrite string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		handle("error opening file to append: ", err)

	}
	if _, err := f.WriteString(textToWrite); err != nil {
		handle("error appending to file: ", err)

	}
	if err := f.Close(); err != nil {
		handle("error closing file after appending: ", err)

	}

}

// writeFile Generic file handler
func writeFile(filename, textToWrite string) {
	if !fileExists(filename) {
		createFile(filename)
	}
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.WriteString(textToWrite)
	handle("error writing file: ", err)
	err = file.Sync()
	handle("", err)
}

// writeFileBytes Generic file handler
func writeFileBytes(filename string, bytesToWrite []byte) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.Write(bytesToWrite)
	handle("error writing file bytes: ", err)

	err = file.Sync()
	handle("", err)
}

// readFile Generic file handler
func readFile(filename string) string {
	text, err := os.ReadFile(filename)
	handle("Couldnt read the file: ", err)
	return string(text)
}

func readFileBytes(filename string) []byte {
	text, err := os.ReadFile(filename)
	handle("Couldnt read the file: ", err)
	return text
}

// deleteFile Generic file handler
func deleteFile(filename string) {
	err := os.Remove(filename)
	handle("Problem deleting file: ", err)
}

func validJSON(stringToValidate string) bool {
	return json.Valid([]byte(stringToValidate))
}

func splitUserPassStr(userpass string) []string {
	return strings.Split(userpass, ",")
}

func stringExistsInFile(thisString string) bool {
	f, err := os.Open("admin/users.list")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if thisString == scanner.Text() {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return false
}
