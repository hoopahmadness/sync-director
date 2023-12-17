package main

import (
	"fmt"
	"sync"
)

type ClientStatus int

// const (
// 	NOTFOUND = iota
// 	CONNECTED
// )

type Client struct {
	*Device
	deviceId  string
	apiKey    string
	ipAddress string
	// status    ClientStatus
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

func (client *Client) getConnectedDeviceIPs() map[string]string {
	//  rest/system/connections
	return map[string]string{}
}

func (client *Client) ping() {}

func findIPs(clients map[string]*Client, lock *sync.RWMutex) {
	fmt.Println("Entering findIPs")
	fmt.Println(clients)
	if lock == nil {
		lock = &sync.RWMutex{}
	}

	lock.RLock()
	if len(clients) == 0 {
		lock.RUnlock()
		fmt.Println("length of clients is 0, returning")
		return
	}
	IDs := make([]string, len(clients))
	ii := 0
	for id, _ := range clients {
		IDs[ii] = id
		ii++
	}
	lock.RUnlock()
	fmt.Println("got list of client IDs")
	fmt.Println(IDs)

	for _, id := range IDs {
		lock.Lock()
		fmt.Println("Locking map for client")
		fmt.Println(id)
		client, OK := clients[id]
		if !OK {
			lock.Unlock()
			fmt.Println("client no longer exists in map, moving on")
			continue
		}
		if client.ipAddress == "" {
			lock.Unlock()
			fmt.Println("Client has no IP yet, moving on")
			continue
		}
		delete(clients, id)
		lock.Unlock()
		fmt.Println("Unlocked the map")
		connectedIPs := client.getConnectedDeviceIPs()
		for deviceID, IP := range connectedIPs {
			lock.Lock()
			fmt.Println("Locking map for other device")
			fmt.Println(deviceID)
			otherDeviceClient, OK := clients[deviceID]
			if OK {
				println("Updating the IP for this device")
				otherDeviceClient.ipAddress = IP
			}
			lock.Unlock()
			println("Unlocked, moving on to next returned device ID")
		}
		println("Moving on to the next client ID")
	}
	println("recursing")
	findIPs(clients, lock)
}
