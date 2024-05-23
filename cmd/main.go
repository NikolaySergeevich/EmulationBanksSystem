package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	pay "testbanc/pkg/paymentSystem"
	"testbanc/pkg/repository"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup
	db := repository.NewLocalDB()
	ps, err := pay.NewPaymentSystem(ctx, db)

	if err != nil {
		log.Fatal(err)
	}

	//Создание десяти счетов
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ps.CreateAccount(ctx); err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()

	// Демонстрация получения номеров счетов страны и уничтожения денег
	ibanCountry, err := ps.GetCountryIBAN(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ibanCountry)
	ibanDestruction, err := ps.GetDestructionIBAN(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ibanDestruction)
	// Эмиссия
	if err := ps.EmitMoney(ctx, 2300); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Эмиссия успешно проведена")
	// Уничтожение денег
	if err := ps.DestroyMoney(ctx, ibanCountry, 300); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Указанная сумма успешна уничтожена")
	//Пример неудачного уничтожения денег
	if err := ps.DestroyMoney(ctx, "BY29JYEI0873459522488868084694", 300); err != nil {
		log.Println(err)
	}
	if err := ps.DestroyMoney(ctx, ibanCountry, 300000); err != nil {
		log.Println(err)
	}
	// Перевод денег с одного счёта на другой
	if err := ps.TransferMoney(ctx, ibanCountry, ibanDestruction, 200); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Указанная сумма успешна переведена первым способом")

	// Перевод денег с одного счёта на другой с JSON функциями
	if err := ps.TransferMoneyJSON(ctx, ibanCountry, ibanDestruction, 100); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Указанная сумма успешна переведена вторым способом")

	// Неудачные переводы денег с одгого счёта на другой
	if err := ps.TransferMoney(ctx, "BY29JYEI0873459522488868084694", ibanDestruction, 200); err != nil {
		log.Println(err)
	}
	if err := ps.TransferMoney(ctx, ibanCountry, "BY29JYEI0873459522488868084694", 200); err != nil {
		log.Println(err)
	}
	if err := ps.TransferMoneyJSON(ctx, ibanCountry, ibanDestruction, 554545); err != nil {
		log.Println(err)
	}
	// Полуение всех счетов, включая служебыне
	allAcc, err := ps.GetAccounts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(allAcc)
}
