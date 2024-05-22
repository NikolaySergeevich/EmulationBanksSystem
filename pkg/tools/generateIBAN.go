package tools

import (
	"math/rand"
	"sync"
	"time"
)

type Iban struct {
	Iban string
}

// Генерирует уникальный номер IBAN с указанным форматом.
func GenerateIBAN() Iban {
	const prefix = "BY"
	const length = 28

	// Сгенерирует случайные символы для оставшейся части IBAN
	letter := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number := "0123456789"

	var wg sync.WaitGroup
	var mu sync.Mutex
	rand.Seed(time.Now().UnixNano())
	var ibanOne []byte
	var ibanTwo []byte
	var ibanThree []byte
	for i := 0; i < length-len(prefix); i++ {
		wg.Add(1)
		go func(letter string, number string, i int) {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			if i >= 2 && i <= 3 {
				ibanOne = append(ibanOne, number[rand.Intn(len(number))])
			}
			if i >= 4 && i <= 7 {
				ibanTwo = append(ibanTwo, letter[rand.Intn(len(letter))])
			} else {
				ibanThree = append(ibanThree, number[rand.Intn(len(number))])
			}
		}(letter, number, i)
	}
	wg.Wait()
	iban := prefix + string(ibanOne) + string(ibanTwo) + string(ibanThree)
	return Iban{Iban: iban}
}
