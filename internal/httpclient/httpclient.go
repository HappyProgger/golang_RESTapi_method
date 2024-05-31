package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func SendData(data Data, count int) {
	url_api := "https://development.kpi-drive.ru/_api/facts/save_fact"
	token := "48ab34464a5573519725deb5865cc74c"

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

	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", url_api, body)
	if err!= nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err!= nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err!= nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(bodyBytes))
	fmt.Println(count)
	fmt.Println(strconv.Itoa(count))
}
