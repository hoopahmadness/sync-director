package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var devicesById = map[string]*Device{}
var deviceConnections *DeviceWeb
var netFoldersById = map[string]*NetworkFolder{}

func main() {
	// start with map of device IDs to Devices
	deviceConnections = newDeviceWeb()
	myDevices := createMyDevices()
	for _, dev := range myDevices {
		// populate my device map
		devicesById[dev.DeviceId] = dev
	}
	for _, dev := range myDevices {
		// get connected devices and register each connection in the web
		connectedDevs, err := dev.GetConnectedDevices()
		if err != nil {
			fmt.Println("Get ConnectedDevices crashed")
			panic("Get ConnectedDevices crashed")
		}
		for _, connectedDev := range connectedDevs {
			deviceConnections.NewDeviceConnection(dev, connectedDev)
		}
	}
	// go back through each device and get the folders
	for _, device := range devicesById {
		folders, _ := device.GetFolders()
		for _, folder := range folders {
			// turn each folder into a network folder if it isn't already
			netFolder, OK := netFoldersById[folder.Id]
			if !OK {
				netFolder = newNetworkFolder(folder.Id)
				netFoldersById[folder.Id] = netFolder
			}
			netFolder.IngestFolder(folder)
		}
	}
	page := components.NewPage()
	page.AddCharts(graphDeviceWeb(devicesById, deviceConnections))
	page.SetLayout(components.PageFlexLayout)
	f, err := os.Create("ui")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func createNodesAndLinksForDevices(deviceMap map[string]*Device, web *DeviceWeb) ([]opts.GraphNode, []opts.GraphLink) {
	nodeArr := []opts.GraphNode{}
	linkArr := []opts.GraphLink{}
	for _, dev := range deviceMap {
		symbol := ""
		switch dev.Status {
		case UNKNOWN:
			symbol = "diamond"
		case CONNECTED:
			symbol = "circle"
		case OFFLINE:
			symbol = "rect"
		case OUTOFNETWORK:
			symbol = "triangle"
		default:
			symbol = "none"
		}
		newNode := opts.GraphNode{
			Name:       dev.Nickname,
			Symbol:     symbol,
			SymbolSize: 8,
			Tooltip: &opts.Tooltip{
				Show:           opts.Bool(true),
				Trigger:        "item",
				TriggerOn:      "mousemove",
				ValueFormatter: dev.DeviceId,
			},
		}
		dev.GraphNode = &newNode
		nodeArr = append(nodeArr, newNode)
	}
	for pair, _ := range web.AllPairs {
		if pair.GraphLink == nil {
			newLink := opts.GraphLink{
				Source: pair.Dev1.GraphNode.Name,
				Target: pair.DevA.GraphNode.Name,
			}
			pair.GraphLink = &newLink
			linkArr = append(linkArr, newLink)
		}
	}

	return nodeArr, linkArr
}

func graphDeviceWeb(deviceMap map[string]*Device, web *DeviceWeb) *charts.Graph {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "All Devices"}),
	)
	nodes, links := createNodesAndLinksForDevices(deviceMap, web)
	graph.AddSeries("Device Web", nodes, links,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Force: &opts.GraphForce{Repulsion: 3000},
			},
		))
	return graph
}
