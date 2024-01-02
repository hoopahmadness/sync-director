package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClientStatus int

const (
	UNKNOWN = iota
	CONNECTED
	OFFLINE
	OUTOFNETWORK
)

type ConnectedDevicesResponse struct {
	Connections map[string]struct {
		Address string
	}
}

// top level is map of folder IDs to OfferedByObject
// OfferedBy is a map of the other device ID to some extraneous details
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
	// SharedDevices map[string]struct {
	// 	Pending bool
	// }

}

// Clients live inside of a Device and do all of the HTTP related dirty work.
// the parentDevice field is just for convenience because we want Clients to be able to set
// pointers to its parent in folders, etc
type Client struct {
	deviceId     string
	apiKey       string
	ipAddress    string
	nickname     string
	client       *http.Client
	parentDevice *Device
	status       ClientStatus
}

func (client *Client) querySyncedFolders() (GetFolderResponse, error) {
	/*	/rest/config/folders

		GET returns all folders respectively devices as an array. PUT takes an array and POST a single object. In both cases if a given folder/device already exists, it’s replaced, otherwise a new one is added.
	*/
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
	fmt.Println("Getting folders")
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

// func (client *Client) getConnectedDeviceIPs() (map[string]string, error) {
// 	fmt.Println("Getting connected devices")
// 	//  rest/system/connections
// 	message, err := client.get(client.generateURL("/rest/system/connections"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	response := &ConnectedDevicesResponse{}
// 	connections := map[string]string{}
// 	err = json.Unmarshal(message, &response)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(response.Connections["LANYALJ-PJIJXP2-MXWY52F-G7GUU5W-QS6PTAS-K26GPCZ-DFP3GSF-DKDXSAL"].Address)
// 	fmt.Println("yo")
// 	for id, devInfo := range response.Connections {
// 		if devInfo.Address != "" {
// 			fmt.Println("Got an IP!")
// 			fmt.Println(id)
// 			fmt.Println(devInfo.Address)
// 			connections[id] = devInfo.Address
// 		}
// 	}
// 	return connections, err
// }

func (client *Client) generateURL(endpoint string) string {
	return "https://" + client.ipAddress + endpoint
}

func (client *Client) ping() {}

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
		client.client = &http.Client{Transport: customTransport}
	}
}

// func findIPs(clients map[string]*Client, lock *sync.RWMutex) {
// 	fmt.Println("Entering findIPs")
// 	fmt.Println(clients)
// 	if lock == nil {
// 		lock = &sync.RWMutex{}
// 	}
// 	cameInWithOne := len(clients) == 1

// 	lock.RLock()
// 	if len(clients) == 0 {
// 		lock.RUnlock()
// 		fmt.Println("length of clients is 0, returning")
// 		return
// 	}
// 	IDs := make([]string, len(clients))
// 	ii := 0
// 	for id, _ := range clients {
// 		IDs[ii] = id
// 		ii++
// 	}
// 	lock.RUnlock()
// 	fmt.Println("got list of client IDs")
// 	fmt.Println(IDs)

// 	for _, id := range IDs {
// 		lock.Lock()
// 		fmt.Println("Locking map for client in loop:")
// 		fmt.Println(id)
// 		client, OK := clients[id]
// 		if !OK {
// 			lock.Unlock()
// 			fmt.Println("client no longer exists in map, moving on")
// 			continue
// 		}
// 		if client.ipAddress == "" {
// 			lock.Unlock()
// 			fmt.Println("Client has no IP yet, moving on")
// 			continue
// 		}
// 		delete(clients, id)
// 		lock.Unlock()
// 		fmt.Println("Unlocked the map")
// 		connectedIPs, err := client.getConnectedDeviceIPs()
// 		fmt.Println("Connected devices with IPs:")
// 		fmt.Println(connectedIPs)
// 		if err != nil {
// 			fmt.Println(err)
// 			fmt.Println("Got an error getting connected devices; moving on")
// 			continue
// 		}
// 		for deviceID, IP := range connectedIPs {
// 			lock.Lock()
// 			fmt.Println("Locking map for other device")
// 			fmt.Println(deviceID)
// 			otherDeviceClient, OK := clients[deviceID]
// 			if OK {
// 				println("Updating the IP for this device")
// 				otherDeviceClient.ipAddress = IP
// 			}
// 			lock.Unlock()
// 			println("Unlocked, moving on to next returned device ID")
// 		}
// 		println("Moving on to the next client ID")
// 	}
// 	if len(clients) == 1 && cameInWithOne {
// 		fmt.Println("we haven't improved the list any, we should stop recursiving")
// 		return
// 	}
// 	println("recursing")
// 	findIPs(clients, lock)
// }
