package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Error struct {
	Message string
}

type Result struct {
	Service  string
	Url      string
	Response string
	Error    *Error
}

const brasilapiUrl = "https://brasilapi.com.br/api/cep/v1/{cep}"
const viacepUrl = "http://viacep.com.br/ws/{cep}/json/"

func main() {
	channel := make(chan *Result)

	go findCepBrasilapi("01153000", channel)
	go findCepViacep("01153000", channel)

	result := <-channel
	if result.Error != nil {
		fmt.Println("O servico '" + result.Service + "' com a url '" + result.Url + "' retornou com erro : ")
		fmt.Println(result.Error.Message)
		return
	}
	fmt.Println("O servico '" + result.Service + "' com a url '" + result.Url + "' retornou com sucesso : ")
	fmt.Println(result.Response)
}

func findCepBrasilapi(cep string, channel chan<- *Result) {
	//fmt.Println("Calling brasilapi...")
	url := strings.Replace(brasilapiUrl, "{cep}", cep, 1)
	//fmt.Println("URL : " + url)

	result := &Result{
		Service: "brasilapi",
		Url:     url,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	_ = res.Body.Close()

	result.Response = string(body)
	channel <- result
}

func findCepViacep(cep string, channel chan<- *Result) {
	//fmt.Println("Calling viacep...")
	url := strings.Replace(viacepUrl, "{cep}", cep, 1)
	//fmt.Println("URL : " + url)

	result := &Result{
		Service: "viacep",
		Url:     url,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		errorHandler(result, err, channel)
		return
	}
	_ = res.Body.Close()

	result.Response = string(body)
	channel <- result
}

func errorHandler(result *Result, err error, channel chan<- *Result) {
	result.Error = &Error{
		Message: err.Error(),
	}
	channel <- result
}
