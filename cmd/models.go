package run

import "fmt"

// ServiceTransport ...
type ServiceTransport struct {
	IP              string `json:"ip"`
	Port            string `json:"port"`
	BalanceType     string `json:"balanceType"`
	RoutingType     string `json:"routingType"`
	Protocol        string `json:"protocol"`
	HealthcheckType string `json:"healthcheckType"`
	HelloTimer      string `json:"helloTimer"`
	ResponseTimer   string `json:"responseTimer"`
	AliveThreshold  string `json:"aliveThreshold"`
	DeadThreshold   string `json:"deadThreshold"`
	Quorum          string `json:"quorum"`
	// Hysteresis      string `json:"hysteresis"`
	ApplicationServersTransport []*ApplicationServerTransport `json:"-"`
}

// ApplicationServerTransport ...
type ApplicationServerTransport struct {
	IP                 string `json:"ip"`
	Port               string `json:"port"`
	HealthcheckAddress string `json:"healthcheckAddress"`
}

// Release stringer interface for print/log data in []*ServiceTransport
func (serviceTransport *ServiceTransport) String() string {
	return fmt.Sprintf("IP: %v, Port: %v, BalanceType: %v, RoutingType: %v, Protocol: %v, HealthcheckType: %v, HelloTimer: %v, ResponseTimer: %v, AliveThreshold: %v, DeadThreshold: %v, Quorum: %v, ApplicationServersTransport: %v",
		serviceTransport.IP,
		serviceTransport.Port,
		serviceTransport.BalanceType,
		serviceTransport.RoutingType,
		serviceTransport.Protocol,
		serviceTransport.HealthcheckType,
		serviceTransport.HelloTimer,
		serviceTransport.ResponseTimer,
		serviceTransport.AliveThreshold,
		serviceTransport.DeadThreshold,
		serviceTransport.Quorum,
		serviceTransport.ApplicationServersTransport)
}

// Release stringer interface for print/log data in []*ApplicationServerTransport
func (applicationServerTransport *ApplicationServerTransport) String() string {
	return fmt.Sprintf("IP: %v, Port: %v, HealthcheckAddress: %v",
		applicationServerTransport.IP,
		applicationServerTransport.Port,
		applicationServerTransport.HealthcheckAddress)
}
