package requests

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_request(t *testing.T)  {
	var (
		params=make(map[string][]string,0)
		//data=make(map[string]interface{})
	)

	params["namespace"]=[]string{"xmzl-default"}
	params["app_name"]=[]string{"devops-a4"}
	//get方法
	result :=Get("http://127.0.0.1:777/v1/k8s/deployments/detail",params,10)

	//post
	//data["url"] = "http://gd-gitlab.dc.servyou-it.com"
	//data["token"] = "7spS9Dy3Q9WUxasqudB9"
	//data["projectId"] = 931
	//data["srcBranchName"] = "master"
	//data["destBranchName"] = "1.0.1"
	//result :=Post("http://127.0.0.1:8088/v1/gitlab/project/branch/create",data,10)


	jsonBytes, err := json.Marshal(result.Json)
	if err!=nil{
		t.Error(err)
	}
	strmes := string(jsonBytes)
	fmt.Println(strmes)

}