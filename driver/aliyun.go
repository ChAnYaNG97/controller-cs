package driver

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"strings"
)

const (
	AliyunAccessKey    = "<your access key>"
	AliyunAccessSecret = "<your access secret>"
)

type AliyunClient struct {
	*sdk.Client
}

func NewAliyunClient() (*AliyunClient, error) {
	client, err := sdk.NewClientWithAccessKey("default", AliyunAccessKey, AliyunAccessSecret)
	if err != nil {
		return nil, err
	}
	return &AliyunClient{
		client,
	}, err
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
	body := `{
	   "cluster_type": "Kubernetes",
	   "name": "#NAME#",
	   "region_id": "cn-hangzhou",
	   "disable_rollback": true,
	   "timeout_mins": 60,
	   "kubernetes_version": "1.16.9-aliyun.1",
	   "snat_entry": false,
	   "endpoint_public_access": false,
	   "cloud_monitor_flags": false,
	   "deletion_protection": false,
	   "node_cidr_mask": "26",
	   "proxy_mode": "ipvs",
	   "tags": [],
	   "addons": [
	       {
	           "name": "flannel"
	       },
	       {
	           "name": "csi-plugin"
	       },
	       {
	           "name": "csi-provisioner"
	       },
	       {
	           "name": "nginx-ingress-controller",
	           "config": "{\"IngressSlbNetworkType\":\"internet\"}"
	       }
	   ],
	   "os_type": "Linux",
	   "platform": "AliyunLinux",
	   "node_port_range": "30000-32767",
	   "login_password": "YANGchen970617",
	   "cpu_policy": "none",
	   "master_count": 3,
	   "master_vswitch_ids": [
	       "vsw-bp1bgp5v4duhtkykk8t4h",
	       "vsw-bp1bgp5v4duhtkykk8t4h",
	       "vsw-bp1bgp5v4duhtkykk8t4h"
	   ],
	   "master_instance_types": [
	       "ecs.n1.medium",
	       "ecs.n1.medium",
	       "ecs.n1.medium"
	   ],
	   "master_system_disk_category": "cloud_ssd",
	   "master_system_disk_size": 120,
	   "runtime": {
	       "name": "docker",
	       "version": "19.03.5"
	   },
	   "worker_instance_types": [
	       "ecs.n1.medium"
	   ],
	   "num_of_nodes": 3,
	   "worker_system_disk_category": "cloud_efficiency",
	   "worker_system_disk_size": 120,
	   "vpcid": "vpc-bp1tz33v9lv47nptrykbu",
	   "worker_vswitch_ids": [
	       "vsw-bp1bgp5v4duhtkykk8t4h"
	   ],
	   "is_enterprise_security_group": true,
	   "container_cidr": "172.20.0.0/16",
	   "service_cidr": "172.21.0.0/20"
	}`
	body = strings.ReplaceAll(body, "#NAME#", req.Name)

	var bodyMap map[string]interface{}
	_ = json.Unmarshal([]byte(body), &bodyMap)
	bodyMap["name"] = req.Name
	bodyBytes, _ := json.Marshal(bodyMap)
	body = string(bodyBytes)

	request.Content = []byte(body)
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
		}, nil
	}
	return nil, fmt.Errorf("create cluster error")
}

func (client *AliyunClient) DeleteCluster(clusterId string) error {
	request := requests.NewCommonRequest()
	request.Method = "DELETE"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters/" + clusterId
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

func (client *AliyunClient) GetClusterStatus(clusterId string) (*GetStatusReponse, error) {
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https" // https | http
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = "/clusters/" + clusterId
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
	status, ok := respMap["status"].(string)
	if ok {
		return &GetStatusReponse{
			Status: status,
		}, nil
	}
	return nil, fmt.Errorf("get cluster status error")
}

func (client *AliyunClient) Print(text string) {
	fmt.Println(text)
}
