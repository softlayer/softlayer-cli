package utils

type Access struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	PrivateIPAddress string `json:"private_ip_address,omitempty"`
	SourceSubnet     string `json:"source_subnet,omitempty"`
	HostIQN          string `json:"host_iqn,omitempty"`
	UserName         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	AllowedHostID    string `json:"allowed_host_id,omitempty"`
}

type AccessByID []Access

func (a AccessByID) Len() int {
	return len(a)
}
func (a AccessByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByID) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}

type AccessByName []Access

func (a AccessByName) Len() int {
	return len(a)
}
func (a AccessByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

type AccessByType []Access

func (a AccessByType) Len() int {
	return len(a)
}
func (a AccessByType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByType) Less(i, j int) bool {
	return a[i].Type < a[j].Type
}

type AccessByPrivateIPAddress []Access

func (a AccessByPrivateIPAddress) Len() int {
	return len(a)
}
func (a AccessByPrivateIPAddress) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByPrivateIPAddress) Less(i, j int) bool {
	return a[i].PrivateIPAddress < a[j].PrivateIPAddress
}

type AccessBySourceSubnet []Access

func (a AccessBySourceSubnet) Len() int {
	return len(a)
}
func (a AccessBySourceSubnet) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessBySourceSubnet) Less(i, j int) bool {
	return a[i].SourceSubnet < a[j].SourceSubnet
}

type AccessByHostIQN []Access

func (a AccessByHostIQN) Len() int {
	return len(a)
}
func (a AccessByHostIQN) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByHostIQN) Less(i, j int) bool {
	return a[i].HostIQN < a[j].HostIQN
}

type AccessByUserName []Access

func (a AccessByUserName) Len() int {
	return len(a)
}
func (a AccessByUserName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByUserName) Less(i, j int) bool {
	return a[i].UserName < a[j].UserName
}

type AccessByPassword []Access

func (a AccessByPassword) Len() int {
	return len(a)
}
func (a AccessByPassword) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByPassword) Less(i, j int) bool {
	return a[i].Password < a[j].Password
}

type AccessByAllowedHostID []Access

func (a AccessByAllowedHostID) Len() int {
	return len(a)
}
func (a AccessByAllowedHostID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AccessByAllowedHostID) Less(i, j int) bool {
	return a[i].AllowedHostID < a[j].AllowedHostID
}
