package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/prometheus/common/log"
	"io/ioutil"
)

const (
	AliyunAccessKey    = "<your accessKey>"
	AliyunAccessSecret = "<your accessSecret>"
)

type AliyunClient struct {
	*sdk.Client
}

var client *AliyunClient

func GetAliyunClient() *AliyunClient {
	if client != nil {
		return nil
	}
	client, err := sdk.NewClientWithAccessKey("default", AliyunAccessKey, AliyunAccessSecret)
	if err != nil {
		return nil
	}

	return &AliyunClient{
		client,
	}
}

func (client *AliyunClient) CreateCluster(req *CreateRequest) (*CreateResponse, error) {
	// TODO 先create一个vpc
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters"
	request.Headers["Content-Type"] = "application/json"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	body, err := ioutil.ReadFile("driver/config/aliyun.config")
	if err != nil {
		log.Error("readfile error", err)
		return nil, err
	}

	body = bytes.ReplaceAll(body, []byte("#NAME#"), []byte(req.Name))
	body = bytes.ReplaceAll(body, []byte("#REGION#"), []byte(req.Region))
	body = bytes.ReplaceAll(body, []byte("#INSTANCE_TYPE"), []byte(req.InstanceType))

	request.Content = body
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	var respMap map[string]interface{}
	_ = json.Unmarshal(response.GetHttpContentBytes(), &respMap)
	clusterId, ok := respMap["cluster_id"].(string)
	if ok {
		return &CreateResponse{
			ClusterId: clusterId,
			VPCId:     "vpc-bp1tz33v9lv47nptrykbu",
		}, nil
	}
	return nil, fmt.Errorf("create cluster error")
}

func (client *AliyunClient) DeleteCluster(req *CommonRequest) error {
	request := requests.NewCommonRequest()
	request.Method = "DELETE"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters/" + req.ClusterId
	request.Headers["Content-Type"] = "application/json"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	body := `{}`
	request.Content = []byte(body)
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}
	fmt.Print(response.GetHttpContentString())
	return nil
}

func (client *AliyunClient) GetClusterStatus(req *CommonRequest) (*GetStatusResponse, error) {
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters/" + req.ClusterId
	request.Headers["Content-Type"] = "application/json"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	body := `{}`
	request.Content = []byte(body)
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	var respMap map[string]interface{}
	_ = json.Unmarshal(response.GetHttpContentBytes(), &respMap)
	status, ok := respMap["state"].(string)

	//running：集群正在运行的。
	//stopped：集群已经停止运行。
	//deleted：集群已经被删除。
	//delete_failed：集群删除失败。
	//failed：集群创建失败。

	if ok {
		fmt.Println(status)
		return &GetStatusResponse{
			Status: status,
		}, nil
	}
	return nil, fmt.Errorf("get cluster status error")
}

func (client *AliyunClient) Print(text string) {
	fmt.Println(text)
}

func (client *AliyunClient) IsClusterExist(req *CommonRequest) (*ExistResponse, error) {
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters/" + req.ClusterId
	request.Headers["Content-Type"] = "application/json"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	body := `{}`
	request.Content = []byte(body)
	response, _ := client.ProcessCommonRequest(request)
	//if err != nil {
	//
	//	fmt.Println("IsClusterExist Response")
	//	fmt.Println(response, err)
	//	return nil, err
	//}
	var respMap map[string]interface{}
	_ = json.Unmarshal(response.GetHttpContentBytes(), &respMap)
	errCode, ok := respMap["code"].(string)

	//running：集群正在运行的。
	//stopped：集群已经停止运行。
	//deleted：集群已经被删除。
	//delete_failed：集群删除失败。
	//failed：集群创建失败。

	if ok && errCode == "ErrorQueryCluster" {
		fmt.Println(errCode)
		return &ExistResponse{
			Exist: false,
		}, nil
	}

	_, ok = respMap["state"].(string)
	if ok {
		return &ExistResponse{
			Exist: true,
		}, nil
	}
	return nil, fmt.Errorf("cluster exists error")

}
