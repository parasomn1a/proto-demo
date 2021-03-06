package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proto-demo/app/pb"
	"time"
)

func main() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	go server()
	go client()
	<-sign
}

func server() {
	r := gin.New()
	r.POST("/login", func(c *gin.Context) {
		var (
			loginReq pb.LoginReq
			loginRes pb.LoginRes
			status   int
			bytes    []byte
			err      error
		)
		if err = c.ShouldBindWith(&loginReq, binding.ProtoBuf); err != nil {
			status = http.StatusUnavailableForLegalReasons
			loginRes = pb.LoginRes{
				Code: 1,
				Msg:  err.Error(),
			}
		} else {
			status = http.StatusOK
			loginRes = pb.LoginRes{
				Code: 0,
				Msg:  fmt.Sprintf("hello %s, your password is %s", loginReq.Username, loginReq.Password),
			}
		}
		if bytes, err = proto.Marshal(&loginRes); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("bytes", bytes)
		c.ProtoBuf(status, &loginRes)
	})
	r.Run(":8080")
}

func client() {
	var (
		client = &http.Client{}
		req    *http.Request
		reqPb  = &pb.LoginReq{
			Username: "admin",
			Password: "jsangorno",
		}
		err       error
		data      []byte
		resp      *http.Response
		respBytes []byte
		respPb    = &pb.LoginRes{}
	)

	data, err = proto.Marshal(reqPb)
	if err != nil {
		log.Fatalln(err.Error())
	}
	req, err = http.NewRequest("POST", "http://127.0.0.1:8080/login", bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/protobuf")
	for true {
		time.Sleep(time.Second)
		if resp, err = client.Do(req); err != nil {
			log.Fatalln(err.Error())
		}
		if respBytes, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Fatalln(err.Error())
		}
		if err = proto.Unmarshal(respBytes, respPb); err != nil {
			log.Fatalln(err.Error())
		}
		log.Println("client", "response get: ", respPb)
	}
}
