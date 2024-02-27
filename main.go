package main

import (
	"bufio"
	"dz_int_p/modules/creaters"
	"dz_int_p/modules/generators"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	fileFolder = "./_files"
	file1      = fileFolder + "/messages.txt"
	file2      = fileFolder + "/readonly.txt"
	file3      = fileFolder + "/messages_ioutil.txt"

	timeFormat = "2006-01-02 15:04:05"
)

func main() {
	if err := creaters.FolderIfNotExists(fileFolder); err != nil {
		log.Panic("ERROR: ", err)
	}

	if err := task1(file1, timeFormat); err != nil {
		log.Println("ERROR: ", err)
	}
	printTaskSeparator()

	if err := task2(file1); err != nil {
		log.Println("ERROR: ", err)
	}
	printTaskSeparator()

	if err := task3(file2); err != nil {
		log.Println("ERROR: ", err)
	}
	printTaskSeparator()

	if err := task41(file3, timeFormat); err != nil {
		log.Println("ERROR: ", err)
	}
	printTaskSeparator()

	if err := task42(file3); err != nil {
		log.Println("ERROR: ", err)
	}
	printTaskSeparator()

	task5()
}

func task1(filePatch, timeFormat string) error {
	file, err := os.OpenFile(filePatch, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 1
	for {
		fmt.Print("Введите сообщение (или 'exit' для выхода): ")
		scanner.Scan()
		text := scanner.Text()

		if text == "exit" {
			break
		}

		timestamp := time.Now().Format(timeFormat)
		line := fmt.Sprintf("%d %s %s\n", lineNumber, timestamp, text)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("ошибка при записи в файл: %s", err)
		}

		lineNumber++
	}

	return nil
}

func task2(filePatch string) error {
	file, err := os.Open(filePatch)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %s", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ошибка при получении информации о файле: %s", err)
	}

	if stat.Size() == 0 {
		return fmt.Errorf("файл %s пуст", filePatch)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка при чтении файла: %s", err)
	}

	return nil
}

func task3(filePatch string) error {
	file, err := os.Create(filePatch)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %s", err)
	}
	defer file.Close()

	err = os.Chmod(filePatch, 0444) // Устанавливаем права только на чтение
	if err != nil {
		return fmt.Errorf("ошибка при изменении прав файла: %s", err)

	}

	file, err = os.OpenFile("readonly.txt", os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("файл %s открыт только для чтения", filePatch)
	}

	_, err = file.WriteString("Тест")
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл:", err)
	} else {
		fmt.Println("Данные успешно записаны")
	}

	return nil
}

func task41(filePatch, timeFormat string) error {
	lineNumber := 1

	for {
		fmt.Print("Введите сообщение (или 'exit' для выхода): ")
		var text string
		fmt.Scanln(&text)

		if text == "exit" {
			break
		}

		timestamp := time.Now().Format(timeFormat)
		newLine := fmt.Sprintf("%d %s %s\n", lineNumber, timestamp, text)

		// Считываем текущее содержимое файла
		content, err := ioutil.ReadFile(filePatch)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("ошибка при чтении файла: %s", err)
		}

		// Добавляем новую строку к существующему содержимому
		content = append(content, newLine...)

		// Перезаписываем файл с новым содержимым
		err = ioutil.WriteFile(filePatch, content, 0644)
		if err != nil {
			return fmt.Errorf("ошибка при записи в файл: %s", err)
		}

		lineNumber++
	}

	return nil
}

func task42(filePatch string) error {
	content, err := ioutil.ReadFile(filePatch)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %s", err)

	}

	if len(content) == 0 {
		return fmt.Errorf("файл %s пуст", filePatch)
	}

	fmt.Println(string(content))

	return nil
}

func task5() {
	var n int

	fmt.Print("Введите количество пар скобок: ")
	fmt.Scan(&n)

	combinations := generators.ParenthesisString(n)
	fmt.Println(combinations)
}

func printTaskSeparator() {
	fmt.Println("\n-----------------------------------------------\n")
}
