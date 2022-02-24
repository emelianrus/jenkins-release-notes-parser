package releaseParser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SendRequest(url string) []byte {

	fmt.Println("get URL:" + url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return body
}
