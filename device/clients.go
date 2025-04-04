package device

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type ClientStatus string

const (
	UNKNOWN      = "UNKNOWN"
	CONNECTED    = "CONNECTED"
	OFFLINE      = "OFFLINE"
	OUTOFNETWORK = "OUTOFNETWORK"
)

var OrderedStatuses = map[ClientStatus]int{
	CONNECTED:    10,
	OUTOFNETWORK: 20,
	OFFLINE:      30,
	UNKNOWN:      40,
}

type ConnectedDevicesResponse struct {
	Connections map[string]struct {
		Address string
	}
}

// top level is map of folder IDs to OfferedByObject
// OfferedBy is a map of the other device ID to some extraneous details
// Label is ???
type GetPendingFoldersResponse map[string]struct {
	OfferedBy map[string]struct {
		Label string
	}
}

type GetFolderResponse []*struct {
	Id      string
	Label   string
	Path    string
	Type    string
	Devices []struct {
		DeviceID string
	}
}

// Clients live inside of a Device and do all of the HTTP related dirty work.
// the parentDevice field is just for convenience because we want Clients to be able to set
// pointers to its parent in folders, etc
type Client struct {
	DeviceId     string
	apiKey       string
	IpAddress    string
	Nickname     string
	client       *http.Client
	parentDevice *Device
	Status       ClientStatus
}

func newClient(device *Device, nickname string) *Client {
	c := &Client{
		Nickname:     nickname,
		Status:       OUTOFNETWORK,
		parentDevice: device,
	}
	return c
}

func (c *Client) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (client *Client) addToNetwork(deviceID, apikey, ipAddress string) {
	client.DeviceId = deviceID
	client.apiKey = apikey
	client.IpAddress = ipAddress
	client.Status = OFFLINE
	client.ping()
}

func (client *Client) querySyncedFolders() (GetFolderResponse, error) {
	/*	/rest/config/folders

		GET returns all folders respectively devices as an array. PUT takes an array and POST a single object. In both cases if a given folder/device already exists, it’s replaced, otherwise a new one is added.
	*/
	if client.Status == OUTOFNETWORK || client.Status == OFFLINE {
		return GetFolderResponse{}, nil
	}
	message, err := client.get(client.generateURL("/rest/config/folders"))
	if err != nil {
		return nil, err
	}
	response := GetFolderResponse{}
	err = json.Unmarshal(message, &response)
	return response, err
}

func (client *Client) queryPendingFolders() (GetPendingFoldersResponse, error) {
	// rest/cluster/pending/folders
	if client.Status == OUTOFNETWORK || client.Status == OFFLINE {
		return GetPendingFoldersResponse{}, nil
	}
	message, err := client.get(client.generateURL("/rest/cluster/pending/folders"))
	if err != nil {
		return nil, err
	}
	response := GetPendingFoldersResponse{}
	err = json.Unmarshal(message, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *Client) addFolder(name, id, path string) {
	/*
				/rest/config/folders

				GET returns all folders respectively devices as an array. PUT takes an array and POST a single object. In both cases if a given folder/device already exists, it’s replaced, otherwise a new one is added.

				/rest/config/folders/*id*, /rest/config/devices/*id*
		Put the desired folder- respectively device-ID in place of *id*. GET returns the folder/device for the given ID, PUT replaces the entire config, PATCH replaces only the given child objects and DELETE removes the folder/device.


	*/
}

func (client *Client) addDevice(name, id string) {
	/*
		rest/config/devices

		GET returns all folders respectively devices as an array. PUT takes an array and POST a single object. In both cases if a given folder/device already exists, it’s replaced, otherwise a new one is added.
	*/
}

func (client *Client) queryConnectedDevices() (*ConnectedDevicesResponse, error) {
	//  rest/system/connections
	if client.Status == OUTOFNETWORK || client.Status == OFFLINE {
		return nil, nil
	}
	message, err := client.get(client.generateURL("/rest/system/connections"))
	if err != nil {
		return nil, err
	}
	response := &ConnectedDevicesResponse{}
	err = json.Unmarshal(message, &response)
	return response, err
}

func (client *Client) generateURL(endpoint string) string {
	return "https://" + client.IpAddress + endpoint
}

func (client *Client) ping() {
	/*
		POST /rest/system/ping
		Returns a {"ping": "pong"} object.
	*/
	// fmt.Println("Pinging client")
	// fmt.Println(client.parentDevice.Nickname)
	if client.Status == OUTOFNETWORK {
		return
	}
	_, err := client.get(client.generateURL("/rest/system/ping"))
	if err != nil {
		client.Status = OFFLINE
		return
	}
	client.Status = CONNECTED
}

func (client *Client) get(endpoint string) (json.RawMessage, error) {
	client.initHttp()
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Key", client.apiKey)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	obj := json.RawMessage(body)
	return obj, nil
}

func (client *Client) initHttp() {
	if client.client == nil {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.client = &http.Client{Transport: customTransport, Timeout: 3 * time.Second}
	}
}
