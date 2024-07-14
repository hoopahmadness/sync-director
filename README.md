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