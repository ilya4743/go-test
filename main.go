package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

type Item struct {
	Name string `json:"name"`
}

type Person struct {
	Name string       `json:"name"`
	Age  int          `json:"age"`
	Bag  map[int]Item `json:"bag"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	person_json, err := json.Marshal(Person{"Ilya", 24, map[int]Item{1: Item{"Карандаш"}, 2: Item{"Ручка"}}})
	if err != nil {
		fmt.Println("Ошибка при парсинге JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(person_json)
}

func main() {
	// Устанавливаем максимальное количество потоков, которые могут быть запущены одновременно
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Создаем многопоточный веб-сервер
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", handler)
		http.HandleFunc("/Person", jsonHandler)
		http.ListenAndServe(":8080", nil)
	}()

	// Бесконечный цикл для вывода информации о текущем состоянии сервера
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		fmt.Println("Number of goroutines:", runtime.NumGoroutine())
		fmt.Println("Number of CPUs:", runtime.NumCPU())
		fmt.Println("Number of threads:", runtime.GOMAXPROCS(0))
		fmt.Println("Memory usage:", memStats.Alloc/1024/1024, "MB")
		fmt.Println("-----------------------")
	}
	wg.Wait()
}
