package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getFiles(dir string, ext string) ([]string, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}
	fileNames := []string{}
	for _, dirEntry := range dirEntries {
		if filepath.Ext(dirEntry.Name()) == ext {
			fileNames = append(fileNames, dirEntry.Name())
		}
	}
	return fileNames, nil
}

func generateTable(data [][]string) []byte {
	line := data[0]
	columns := fmt.Sprintf("|%s|%s|%s|%s|%s|\n|----|-----|----|---|---|\n", line[0], line[1], line[2], line[3], line[4])
	rows := []string{}
	for _, line := range data[1:] {
		row := fmt.Sprintf("|%s|[%s](%s)|%s|%s|%s|\n", line[0], line[1], line[5], line[2], line[3], line[4])
		rows = append(rows, row)
	}
	return []byte(columns + strings.Join(rows, ""))
}

func writeToMD(records [][]string, fileName string) error {
	table := generateTable(records)
	file, err := os.Create(fileName + ".md")
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(table)
	return err
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return [][]string{}, (err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, (err)
	}
	return records, nil
}

func generateName(fileName string) string {
	parts := strings.Split(fileName, "_")
	for i, p := range parts {
		parts[i] = strings.Title(p)
	}
	return strings.Join(parts, " ")
}

func generateREADME(dir string, fileNames []string) error {
	rows := []string{}
	cols := "|No|Company|\n|-----|----|\n"
	for i, fileName := range fileNames {
		name := generateName(fileName[:len(fileName)-3])
		row := fmt.Sprintf("|%d|[%s](%s)|\n", i+1, name, dir+fileName)
		rows = append(rows, row)
	}
	table := cols + strings.Join(rows, "")
	file, err := os.Create("../README.md")
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString(table)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fileNames, err := getFiles("../csv/", ".csv")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range fileNames {
		data, err := readCSV("../csv/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		writeToMD(data, "../markdown/"+fileName[:len(fileName)-4])
	}
	markdownFileNames, err := getFiles("../markdown/", ".md")
	err = generateREADME("./markdown/", markdownFileNames)
	if err != nil {
		log.Fatal(err)
	}
}
