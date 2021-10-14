package getMe

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetMe(token string){
	//cmd_flag := flag.NewFlagSet("get-me", flag.PanicOnError)
	if resp,err := http.Get(fmt.Sprintf(`https://api.telegram.org/bot%s/getMe`, token)); err!=nil {
		panic(err)
	} else if body, err := ioutil.ReadAll(resp.Body)  ; err != nil {
		panic(err)
	} else {
		print(string(body))
	}


}
