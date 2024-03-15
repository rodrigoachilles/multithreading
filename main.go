package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	service  string
	url      string
	response string
}

const brasilapiUrl = "https://brasilapi.com.br/api/cep/v1/{cep}"
const viacepUrl = "http://viacep.com.br/ws/{cep}/json/"

func main() {
	channel := make(chan *Result)

	go findCepBrasilapi("01153000", channel)
	go findCepViacep("01153000", channel)

	result := <-channel
	fmt.Println("O servico '" + result.service + "' com a url '" + result.url + "' retornou : ")
	fmt.Println(result.response)
}

func findCepBrasilapi(cep string, channel chan<- *Result) {
	//fmt.Println("Calling brasilapi...")
	url := strings.Replace(brasilapiUrl, "{cep}", cep, 1)
	//fmt.Println("URL : " + url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	_ = res.Body.Close()

	result := &Result{
		service:  "brasilapi",
		url:      url,
		response: string(body),
	}
	channel <- result
}

func findCepViacep(cep string, channel chan<- *Result) {
	//fmt.Println("Calling viacep...")
	url := strings.Replace(viacepUrl, "{cep}", cep, 1)
	//fmt.Println("URL : " + url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	_ = res.Body.Close()

	result := &Result{
		service:  "viacep",
		url:      url,
		response: string(body),
	}
	channel <- result
}
