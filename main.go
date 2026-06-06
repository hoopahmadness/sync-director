/*
Pin down tea views for startup, Device screen, Folder screen
Formatted output of device web similar to biofabric https://biofabric.systemsbiology.net/gallery/pages/SuperQuickBioFabric.html
Tea view for above
Potentially rip out echarts but maybe keep as an auxillary view outside of tui
Add folders, add devices
Decide on how we're going to insantiate device list Json file? yaml?
Server mode?
track the display names of folders and devices from the perspective of other devices. Maybe IP address too as backups for the configurued ip
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hoopahmadness/sync-director/v2/director"
	log "github.com/inconshreveable/log15"
)

func main() {
	// deviceConnections = newDeviceWeb()
	logger := log.New()
	newFile, err := os.Create("logs.txt")
	if err != nil {
		panic("Unable to create logging file")
	}
	logger.SetHandler(log.StreamHandler(newFile, log.LogfmtFormat()))
	model := director.InitialState(logger, createMyDevices)
	prog := tea.NewProgram(&model)
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
