package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const CEP_FILE_NAME = "cep.txt"

func main() {
	for _, cep := range os.Args[1:] {
		data, err := fetchViaCepData(cep)
		if err != nil {
			panic(err)
		}

		createCepFile(data)
	}
}

func fetchViaCepData(cep string) (*ViaCEP, error) {
	preparedUrl := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.Get(preparedUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na requisição: %v\n", err)
		return nil, err
	}

	defer req.Body.Close()

	bodyContent, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler o body: %v\n", err)
		return nil, err
	}

	data, err := parseJsonToObject(bodyContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse do json: %v\n", err)
		return &data, err
	}

	return &data, nil
}

func parseJsonToObject(bodyContent []byte) (ViaCEP, error) {
	var data ViaCEP
	err := json.Unmarshal(bodyContent, &data)
	return data, err
}

func createCepFile(data *ViaCEP) {
	file, err := os.Create(CEP_FILE_NAME)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("CEP: %s\nLogradouro: %s\nBairro: %s\nLocalidade: %s\nUF: %s\n", data.Cep, data.Logradouro, data.Bairro, data.Localidade, data.Uf))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v\n", err)
	}

	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Cidade:", data.Localidade)
}
