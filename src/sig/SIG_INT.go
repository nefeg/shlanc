package sig

import (
	"os"
	"syscall"
	"fmt"
	"os/signal"
	"errors"
)

var ErrSigINT = errors.New("SIG_INT")

func SIG_INT(callback func()){

	go func() {

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGINT)

		<-sig

		fmt.Println("\nReceived termination signal")

		callback()

		os.Exit(0)
	}()
}