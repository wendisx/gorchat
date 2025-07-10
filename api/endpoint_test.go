package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/wendisx/gorchat/internal/constant"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestSignupApi(t *testing.T) {
	t.Skip()
	for i := range 5 {
		payload := map[string]any{
			"userName":     fmt.Sprintf("test%dall", i),
			"userPassword": fmt.Sprintf("test%dlovevm", i),
			"userEmail":    fmt.Sprintf("test%d@good.com", i),
		}
		buf, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("%v", err)
		}
		res, err := http.Post("http://127.0.0.1:3000/user/signup", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Fatalf("%v", err)
		}
		body, _ := io.ReadAll(res.Body)
		log.Printf("%v\n", string(body))
		res.Body.Close()
	}
}
func TestLoginApi(t *testing.T) {
	t.Skip()
	for i := range 5 {
		payload := map[string]any{
			"userEmail":    fmt.Sprintf("test%d@good.com", i),
			"userName":     fmt.Sprintf("test%dall", i),
			"userPassword": fmt.Sprintf("test%dlovevm", i),
		}
		buf, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("%v", err)
		}
		res, err := http.Post("http://127.0.0.1:3000/user/login", "application/json", bytes.NewBuffer(buf))
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == constant.SESSION_KEY {
				log.Printf("%s\n", cookie.Value)
			}
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		body, _ := io.ReadAll(res.Body)
		log.Printf("%v\n", string(body))
		res.Body.Close()
	}
}
func TestUpdateInfoApi(t *testing.T) {
	t.Skip()
	for i := range 5 {
		payload := map[string]any{
			"userId":       6 + i,
			"userEmail":    fmt.Sprintf("test%d@god.com", i),
			"userName":     fmt.Sprintf("test%vvim", i),
			"userPassword": fmt.Sprintf("test%dlovevm", i),
		}
		buf, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("%v", err)
		}
		req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:3000/user/update", bytes.NewBuffer(buf))
		if err != nil {
			log.Fatalf("%v", err)
		}
		res, err := http.Post("http://127.0.0.1:3000/user/login", "application/json", bytes.NewBuffer(buf))
		cookies := res.Cookies()
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, cookie := range cookies {
			if cookie.Name == constant.SESSION_KEY {
				log.Printf("%s\n", cookie.Value)
				req.Header.Set("Content-Type", "application/json")
				req.AddCookie(&http.Cookie{
					Name:  constant.SESSION_KEY,
					Value: cookie.Value,
				})
			}
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		log.Printf("%d\n", res.StatusCode)
		resp.Body.Close()
		res.Body.Close()
	}
}

func TestDeleteApi(t *testing.T) {
	t.Skip()
	for i := range 5 {
		payload := map[string]any{
			"userId":       6 + i,
			"userEmail":    fmt.Sprintf("test%d@god.com", i),
			"userName":     fmt.Sprintf("test%vvim", i),
			"userPassword": fmt.Sprintf("test%dlovevm", i),
		}
		buf, err := json.Marshal(payload)
		reqUrl := fmt.Sprintf("http://127.0.0.1:3000/user/delete?userId=%d", 6+i)
		req, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		if err != nil {
			log.Fatalf("%v", err)
		}
		res, err := http.Post("http://127.0.0.1:3000/user/login", "application/json", bytes.NewBuffer(buf))
		cookies := res.Cookies()
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, cookie := range cookies {
			if cookie.Name == constant.SESSION_KEY {
				log.Printf("%s\n", cookie.Value)
				req.Header.Set("Content-Type", "application/json")
				req.AddCookie(&http.Cookie{
					Name:  constant.SESSION_KEY,
					Value: cookie.Value,
				})
			}
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		log.Printf("%d\n", res.StatusCode)
		resp.Body.Close()
		res.Body.Close()
	}
}
