package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Data struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMoID     int    `json:"indicator_to_mo_id"`
	IndicatorToMoFactID int    `json:"indicator_to_mo_fact_id"`
	Value               int    `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              int    `json:"is_plan"`
	AuthUserID          int    `json:"auth_user_id"`
	Comment             string `json:"comment"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString() string {
	b := make([]byte, 10) // Длина строки
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomTimestamp() time.Time {
	// Генерируем случайное количество секунд в прошлом относительно 1 января 1970 года
	randomSeconds := rand.Int63n(time.Now().Unix()-9466) + 9466
	return time.Unix(randomSeconds, 0)
}

func main() {
	dataChan := make(chan Data, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		newWord := generateRandomString()
		dataChan <- Data{
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
			go func(d Data) {
				defer wg.Done()
				sendData(d, count)
			}(data) // Передаем копию переменной, а не ссылку
		}
	}()

	wg.Wait()
	close(dataChan) // Закрытие канала происходит после ожидания завершения всех горутин
}

// Остальная часть кода остается без изменений...

func sendData(data Data, count int) {
	url_api := "https://development.kpi-drive.ru/_api/facts/save_fact"
	token := "48ab34464a5573519725deb5865cc74c"

	// Преобразование структуры Data в map[string]string
	formData := url.Values{}
	formData.Set("period_start", data.PeriodStart)
	formData.Set("period_end", data.PeriodEnd)
	formData.Set("period_key", data.PeriodKey)
	formData.Set("indicator_to_mo_id", strconv.Itoa(data.IndicatorToMoID))
	formData.Set("indicator_to_mo_fact_id", strconv.Itoa(data.IndicatorToMoFactID))
	formData.Set("value", strconv.Itoa(data.Value))
	formData.Set("fact_time", data.FactTime)
	formData.Set("is_plan", strconv.Itoa(data.IsPlan))
	formData.Set("auth_user_id", strconv.Itoa(data.AuthUserID))
	formData.Set("comment", data.Comment)

	// Создание тела запроса в формате x-www-form-urlencoded
	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", url_api, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(bodyBytes))
	fmt.Println(count)
	fmt.Println(strconv.Itoa(count))
}
