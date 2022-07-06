package main

import (
	"encoding/json"
	"fmt"
	v1 "grpcdemo/api/v1"

	"github.com/gogo/protobuf/proto"
)

func main() {
	Info := &v1.TestInfo{
		Name: "test",
		DemoInfo: &v1.TestInfo_DemoInfo{
			Name:   "demo",
			Age:    10,
			Groups: []*v1.TestInfo_DemoInfo_GroupInfo{},
		},
	}
	infoBytes, _ := proto.Marshal(Info)
	NewInfo := &v1.TestInfo{}

	_ = proto.Unmarshal(infoBytes, NewInfo)
	// NewInfo.DemoInfo.Groups = []*v1.TestInfo_DemoInfo_GroupInfo{}
	//NewInfo.DemoInfo.Groups = append(NewInfo.DemoInfo.Groups, &v1.TestInfo_DemoInfo_GroupInfo{Title: "haha"})

	jsonBytes, _ := json.Marshal(NewInfo)

	fmt.Printf("%s", string(jsonBytes))
	// proto.Marshal slice 会被序列化成nil，而json的序列化对于结构体中slice 为nil和[]是两种不同的结果前者null，后者[]

}
