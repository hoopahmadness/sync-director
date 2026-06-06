OK this is turning into a real mess. I need to take a step back and think through what I'm doing and what I need

I start with a list of devices and access credentials.

For each device:
    I poll its connected/pending devices.
        New devices are also added to the list just to keep track of them
        The connections can probably fit into a DeviceWeb for connections

    I poll its folders
        All folders are added to a map of network Folders
        The folder info returned includes other devices and we can use that to generate device webs for each of those folders
            I WAS going to globalize all those device pair objects in the webs but I don't think that makes sense because each device pair can mark a device as "pending" and that is information specific to a folder (or to the device connections themselves.)
            So instead I think I'll just have completely different device pairs per web. If it blows up later I can come back and optimize this way.
    
    This gives me my list of devices, connections, list of network folders, and the connections for all those folders.

I'm using bubbletea to interact with this information, so I need to think through the different views.

I'm thinking of a schema that divides the pages up into two "spaces" : Folders and Devices. Each space has a List view with multiple items and a singular entity View. Likewise, each of these views can be within the context of a single entity from the other space (so a device list or device view in the context of one folder, or vice versa) or the entire network. So that's (folder vs device) x (item vs list) x (single context vs network) = 8 types of views. 
Folders v devices axis doesn't have a "home" (we can just set default to device but it's easy to flip back and forth between them). But the other two axes can have "list" and "network" as the home base, so let's add a shortcut to quickly get back to list and network views. 
Also, come to think of it, the view for a Device in context of a single Folder should be basically the same as a Folder view in terms of a single device. Likewise, the Device List view for a single folder is essentially the same as the Folder view for the network, and vice versa... So really only 5 views.

Device List View (network)
- m.DeviceList, m.FolderList
- Lists all devices with nickname, status, # connections, # folders. Arrow keys to highlight and scroll, ENTER to select for Device Network View. Perhaps a grid?
- BIOFABRIC: switch to fabric of all devices
- Panel that shows verbose details for selected device
- Allow sorting ASC/DEC by name, status, number of connections, num folders
- Allow filtering by status, hidden, pending actions
- Allow creating new device

Folder List View (network)
- m.DeviceList, m.FolderList
- Lists all folders with nickname, # device pairs. Arrow keys, highlight scroll, ENTER key, Folder View
- No biofabric view
- Panel with verbose details for selected folder
- Sorting by name, num devices, num device pairs
- Filtering of hidden folders
- Allow creating new folder, list target devices with set path, list targets to leave pending.........

Single folder, many devices view
- m.FolderList[1]
- Lists all devices syncing this folder with nickname, status, # connections. Arrow keys to highlight and scroll, ENTER to select for Device View. Perhaps a grid?
- BIOFABRIC: switch to fabric of all devices that have this folder
- Change folder nickname for all devices
- Toggle folder as hidden (in the app)
- Panel that shows verbose details for selected device, including path
- Allow sorting ASC/DEC by name, status, number of connections
- Allow filtering by status, hidden
- Allow syncing this folder to new device 

Single Device, many folders View
- m.DeviceList[1], m.GetFoldersByDevice(dev)
- No biofabric view
- Essentially a recreation of the syncthing web gui
    - Lists various stats and info about device
    - Lists connected devices and folders.
    - Pair with device, either choose from list or input ID
    - Sync new folder, either choose from list or create new folder, choose local path, choose devices to sync with
    - Allows you to highlight and select them to go to the corresponding Device or Folder view. 
    - Allows configuration of unconfigured devices
    - Allow toggling of Hidden (in the app)

Single device, single folder
- m.deviceList[1], m.folderList[2].folders[dev]
