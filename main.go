package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Task1")
	task1()

	fmt.Println("\nTask2")
	task2()
}

func task1() {
	var wg sync.WaitGroup

	reader := bufio.NewReader(os.Stdin)
	numbers := make(chan int)
	squares := make(chan int)
	prompt := make(chan bool)

	// Горутина для вычисления квадрата числа
	wg.Add(1)
	go func() {
		defer wg.Done()
		for n := range numbers {
			squared := n * n
			fmt.Printf("Квадрат: %d\n", squared)

			squares <- squared
		}
		close(squares)
	}()

	// Горутина для удвоения результата
	wg.Add(1)
	go func() {
		defer wg.Done()
		for s := range squares {
			doubled := s * 2
			fmt.Printf("Произведение: %d\n", doubled)

			prompt <- true // Сигнализируем, что можно выводить приглашение к вводу
		}
	}()

	// Горутина для вывода приглашения к вводу после каждого произведения
	go func() {
		for range prompt {
			fmt.Print("Введите число или 'стоп' для выхода: ")
		}
	}()

	// Чтение и обработка ввода пользователя
	prompt <- true // Инициируем первый вывод приглашения к вводу
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения: ", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "стоп" {
			break
		}

		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Введите целое число или 'стоп' для выхода.")
			continue
		}

		numbers <- num
	}
	close(numbers)
	close(prompt)

	wg.Wait()
}

func task2() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина, печатающая квадраты натуральных чисел
	go func() {
		i := 1
		for {
			select {
			case <-stopChan:
				fmt.Println("Выхожу из программы")

				return
			default:
				fmt.Println(i * i)
				i++
				time.Sleep(time.Second)
			}
		}
	}()

	<-stopChan
}
