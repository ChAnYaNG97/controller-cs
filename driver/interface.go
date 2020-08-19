package driver

type CreateRequest struct {
	Name         string
	Region       string
	InstanceType string
}
type CreateResponse struct {
	ClusterId string
	VPCId     string
}

type CommonRequest struct {
	ClusterId string
}

type GetStatusResponse struct {
	Status string
}

type ExistResponse struct {
	Exist bool
}

type ClientInterface interface {
	CreateCluster(*CreateRequest) (*CreateResponse, error)
	DeleteCluster(*CommonRequest) error
	GetClusterStatus(*CommonRequest) (*GetStatusResponse, error)
	IsClusterExist(*CommonRequest) (*ExistResponse, error)

	//for test
	Print(text string)
}
