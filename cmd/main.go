package main

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/mohsanabbas/authorizer/commons/commandhandler/basic"
	"github.com/mohsanabbas/authorizer/commons/eventsource"
	"github.com/mohsanabbas/authorizer/commons/utils"
	"github.com/mohsanabbas/authorizer/config"
	acc "github.com/mohsanabbas/authorizer/internal/account"
)

func main() {
	flag.Parse()
	commandBus, err := getConfig()
	if err != nil {
		glog.Infoln(err)
		os.Exit(1)
	}
	end := make(chan bool)
	//Create Account
	for i := 0; i < 3; i++ {
		go func() {
			uuid, err := utils.UUID()
			if err != nil {
				return
			}
			//1) Create an account
			var account acc.CreateAccount
			account.AggregateID = uuid
			account.ActiveCard = true
			account.AvailableLimit = 108

			commandBus.HandleCommand(account)
			glog.Infof("account %s - account created", uuid)

			//2) Perform a transaction
			time.Sleep(time.Millisecond * 100)
			transaction := acc.PerformTransaction{
				Amount:   100,
				Time:     time.Now().UTC(),
				Merchant: "burger-king",
			}

			transaction.AggregateID = uuid
			transaction.Version = 1

			commandBus.HandleCommand(transaction)
			glog.Infof("account %s - transaction performed", uuid)

			// 3) Perform another transaction it is in goroutine so we have to
			// make sure sleep time between two goroutines
			time.Sleep(time.Millisecond * 110)

			transaction = acc.PerformTransaction{
				Amount:   104,
				Time:     time.Now().UTC(),
				Merchant: "uber-eats",
			}
			transaction.AggregateID = uuid
			transaction.Version = 2

			commandBus.HandleCommand(transaction)
			glog.Infof("account %s - transaction performed", uuid)
		}()
	}
	<-end
}


// getConfig will setups app configs like eventbus ex:nats eventstore ex mongo
// which will give us a commandBus instance
func getConfig() (eventsource.CommandBus, error) {
	//register events
	reg := eventsource.NewEventRegister()
	reg.Set(acc.AccountCreated{})
	reg.Set(acc.TransactionPerformed{})

	//eventsourcing configs
	return config.NewClient(
		config.Mongo("localhost", 27017, "account"), // event store
		config.Nats("nats://localhost:4222", false), // event bus
		config.AsyncCommandBus(30),                  // command bus
		config.WireCommands(
			&acc.Account{},           // aggregate
			basic.NewCommandHandler,   // command handler
			"account",                 // event store bucket
			"account",                 // event store subset
			acc.CreateAccount{},      // command
			acc.PerformTransaction{}, // command
		),
	)
}
