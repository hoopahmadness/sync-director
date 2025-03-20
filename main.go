/*
Pin down tea views for startup, Device screen, Folder screen
Formatted output of device web similar to biofabric https://biofabric.systemsbiology.net/gallery/pages/SuperQuickBioFabric.html
Tea view for above
Potentially rip out echarts but maybe keep as an auxillary view outside of tui
Add folders, add devices
Decide on how we're going to insantiate device list Json file? yaml?
Server mode?
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/inconshreveable/log15"
	// "github.com/go-echarts/go-echarts/v2/charts"
	// "github.com/go-echarts/go-echarts/v2/components"
	// "github.com/go-echarts/go-echarts/v2/opts"
)

// var devicesById = map[string]*Device{}
// var deviceConnections *DeviceWeb
// var netFoldersById = map[string]*NetworkFolder{}

func main() {
	// deviceConnections = newDeviceWeb()
	logger := log.New()
	newFile, err := os.Create("logs.txt")
	if err != nil {
		panic("Unable to create logging file")
	}
	logger.SetHandler(log.StreamHandler(newFile, log.LogfmtFormat()))
	model := initialState(logger)
	if _, err := tea.NewProgram(&model).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
	// page := components.NewPage()
	// // page.AddCharts(graphDeviceWeb(devicesById, deviceConnections, netFoldersById))
	// // page.SetLayout(components.PageFlexLayout)
	// f, err := os.Create("ui")
	// if err != nil {
	// 	panic(err)
	// }
	// page.Render(io.MultiWriter(f))
}

// type webChart struct {
// 	name  string
// 	links []opts.GraphLink
// }

// func createNodesAndLinksForDevices(
// 	deviceMap map[string]*Device,
// 	web *DeviceWeb,
// 	netFolders map[string]*NetworkFolder,
// ) ([]opts.GraphNode, []webChart) {
// 	nodeArr := []opts.GraphNode{}
// 	for _, dev := range deviceMap {
// 		symbol := ""
// 		switch dev.Status {
// 		case UNKNOWN:
// 			symbol = "diamond"
// 		case CONNECTED:
// 			symbol = "circle"
// 		case OFFLINE:
// 			symbol = "rect"
// 		case OUTOFNETWORK:
// 			symbol = "triangle"
// 		default:
// 			symbol = "none"
// 		}
// 		newNode := opts.GraphNode{
// 			Name:       dev.Nickname,
// 			Symbol:     symbol,
// 			SymbolSize: 8,
// 			Tooltip: &opts.Tooltip{
// 				Show:           opts.Bool(true),
// 				Trigger:        "item",
// 				TriggerOn:      "mousemove",
// 				ValueFormatter: dev.DeviceId,
// 			},
// 		}
// 		dev.GraphNode = &newNode
// 		nodeArr = append(nodeArr, newNode)
// 	}
// 	chartArr := []webChart{
// 		{
// 			name:  "All Devices",
// 			links: extractLinksFromWeb(web),
// 		},
// 	}
// 	for id, netFolder := range netFolders {
// 		chartArr = append(chartArr, webChart{
// 			name:  id,
// 			links: extractLinksFromWeb(netFolder.DeviceWeb),
// 		})
// 	}

// 	return nodeArr, chartArr
// }

// func extractLinksFromWeb(web *DeviceWeb) []opts.GraphLink {
// 	linkArr := []opts.GraphLink{}
// 	for pair, _ := range web.AllPairs {
// 		if pair.GraphLink == nil {
// 			newLink := opts.GraphLink{
// 				Source: pair.Dev1.GraphNode.Name,
// 				Target: pair.DevA.GraphNode.Name,
// 			}
// 			pair.GraphLink = &newLink
// 			linkArr = append(linkArr, newLink)
// 		}
// 	}
// 	return linkArr
// }

// func graphDeviceWeb(
// 	deviceMap map[string]*Device,
// 	web *DeviceWeb,
// 	netFoldersByID map[string]*NetworkFolder,
// ) *charts.Graph {
// 	graph := charts.NewGraph()
// 	graph.SetGlobalOptions(
// 		charts.WithTitleOpts(opts.Title{Title: "All Devices"}),
// 	)
// 	nodes, webChart := createNodesAndLinksForDevices(deviceMap, web, netFoldersByID)
// 	for _, chart := range webChart {
// 		graph.AddSeries(chart.name, nodes, chart.links,
// 			charts.WithGraphChartOpts(
// 				opts.GraphChart{
// 					Force: &opts.GraphForce{Repulsion: 2000},
// 				},
// 			))
// 	}
// 	return graph
// }
