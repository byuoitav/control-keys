# control-keys
This service is responsible for generating control keys for each camera. The control keys are used to authenticate control to the cameras.

## Endpoints
* <mark>GET</mark> `/:controlKey/getPreset`
    * Returns the preset for the camera with the control key
    * ex: `http://localhost:8029/871540/getPreset`
* <mark>GET</mark> `/:preset/getControlKey`
    * Returns the control key for the camera with the preset
    * ex: `http://localhost:8029/JET-1106 JET 1106/getControlKey`
* <mark>GET</mark> `/:room/refresh`
    * Refreshes the control keys for the room
* <mark>GET</mark> `/status`
    * Returns "Healthy!" if the service is running

## Environment Variables
* DB_ADDRESS
  * The address of couch
* DB_USERNAME
  * The username for couch
* DB_PASSWORD
  * The password for couch

