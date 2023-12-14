package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func getExchangeRate(ctx context.Context) (any, error) {
	client := &http.Client{
		Timeout: 3000 * time.Millisecond,
	}

	resp, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	print(string(body))
	if err != nil {
		return 0, err
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result["bid"], nil
}

func saveToFile(value any) error {
	content := fmt.Sprintf("Dólar: %f\n", value)
	return os.WriteFile("cotacao.txt", []byte(content), 0644)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		fmt.Println("Erro ao obter cotação:", err)
		return
	}

	err = saveToFile(exchangeRate)
	if err != nil {
		fmt.Println("Erro ao salvar no arquivo:", err)
		return
	}

	fmt.Println("Cotação do dólar salva com sucesso.")
}
