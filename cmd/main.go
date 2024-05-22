package main

import (
	"context"
	"fmt"
	"sync"
	"log"
	"os/signal"
	"syscall"

	// "testbanc/pkg/tools"
	// "testbanc/pkg/account"
	pay "testbanc/pkg/paymentSystem"
	"testbanc/pkg/repository"
)





func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup
	db := repository.NewLocalDB()
	// fmt.Println(db)
	ps, err := pay.NewPaymentSystem(ctx, db)
	// fmt.Println(&ps)
	if err != nil{
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ps.CreateAccount(ctx); err != nil{
				log.Println(err)
			}
		}()
	}
	wg.Wait()

	allAcc, err := ps.GetAccounts(ctx)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(allAcc)
	// // Emit money
	// ps.EmitMoney(1000)

	// for i := 0; i < 10; i++ {
	// 	ps.CreateAccount()
	// }


	
	// accounts, err := ps.GetAccounts()
	// if err != nil{
	// 	log.Println(err)
	// }else {
	// 	fmt.Println(accounts)
	// }
	

	// // Transfer money
	// ps.TransferMoney("state_emission", "state_destruction", 500)

	// // Print account details
	// fmt.Println(ps.GetAccounts()
}
