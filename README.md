# Goby
Goby is an automatic snapshot solution for DigitalOcean droplets.

## Requirements
	GO 1.12+

## Install
`go build -o Goby`

### Set the environment variables
	export DO_API_TOKEN=<API-TOKEN>
	export DO_DROPLET_ID=<DROPLET-ID>


## Usage
|   Arguments | Description  |
| :------------ | :------------ |
| freq  |  How often to perform snapshots in minutes |
| keep |  Amount of days to keep snapshots |

### Example
`Goby -freq 120 -keep 14`

This example will snapshot the droplet every 2 hours and delete any snapshots older than 14 days of the current time. 
## To-Do

- [ ] Add the ability to watch multiple droplets