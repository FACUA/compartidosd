package main

import (
	"fmt"
	"time"
)

func initializeIndex() {
	if err := updateIndex(); err != nil {
		panic(err)
	}

	fmt.Println("Index initialized successfully")

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := updateIndex(); err == nil {
					fmt.Println("Index updated successfully")
				} else {
					fmt.Println("Index update failed!")
					fmt.Println(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
