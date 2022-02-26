package main

import (
    "bufio"
    "fmt"
    "net/http"
	"sync"
)

func main() {

	var wg sync.WaitGroup	

	for i := 0 ; i < 10; i++ {
		wg.Add(1)
		go func(){
			resp, err := http.Get("http://localhost:5000/")
			if err != nil {
				panic(err)
			}
			defer  func(){
				resp.Body.Close()
				wg.Done()
			}()
			fmt.Println("Response status:", resp.Status)
	
			scanner := bufio.NewScanner(resp.Body)
			for i := 0; scanner.Scan() && i < 5; i++ {
				fmt.Println(scanner.Text())
			}
		
			if err := scanner.Err(); err != nil {
				panic(err)
			}
		}()
	}
	for i := 0 ; i < 10; i++ {
		wg.Add(1)
		go func(idx int){
			fmt.Print("compute ")
			fmt.Println(idx)
			resp, err := http.Get("http://localhost:5000/compute")
			if err != nil {
				panic(err)
			}
			defer  func(){
				resp.Body.Close()
				wg.Done()
			}()
			fmt.Println("Response status:", resp.Status)
	
			scanner := bufio.NewScanner(resp.Body)
			for i := 0; scanner.Scan() && i < 5; i++ {
				fmt.Println(scanner.Text())
			}
		
			if err := scanner.Err(); err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()

}