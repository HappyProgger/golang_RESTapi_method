package main

import (
	"fmt"
	"strconv"
	"sync"

	"golang_kpi-drive/data"
	"golang_kpi-drive/httpclient"
	"golang_kpi-drive/utils"
)

func main() {
	dataChan := make(chan data.Data, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		newWord := utils.GenerateRandomString()
		dataChan <- data.Data{
			PeriodStart:         "2024-06-10",
			PeriodEnd:           "2024-06-11",
			PeriodKey:           "month",
			IndicatorToMoID:     227373,
			IndicatorToMoFactID: 0,
			Value:               1,
			FactTime:            "2024-06-10",
			IsPlan:              0,
			AuthUserID:          40,
			Comment:             fmt.Sprintf("buffer Last_name %s", newWord),
		}
	}

	wg.Add(1)
	count := 0
	go func() {
		defer wg.Done()

		for data := range dataChan {
			count++
			fmt.Println(strconv.Itoa(count))
			fmt.Println("Received data:", data)
			go func(d data.Data) {
				defer wg.Done()
				httpclient.SendData(d, count)
			}(data) // Передаем копию переменной, а не ссылку
		}
	}()

	wg.Wait()
	close(dataChan) // Закрытие канала происходит после ожидания завершения всех горутин
}
