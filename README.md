# duthweatherstation

This repo contains our proposed implementation of an Integrated Automatic Weather Station, centered around the FireBeetle 2 ESP32-E microcontroller, that's based on the ESP-WROOM-32E, chosen for its ultra-low power consumption, on-board charging circuit and compatibility with a wide range of sensors. The sensor array includes a DHT11 for humidity, a BMP180 for temperature and barometric pressure, and an MQ135 for air quality monitoring.

Included in this repo:
* The ESP32 firmware
* The Go microservice that handles the data stream from the ESP32 and various SQL transactions with the MySQL database
* Our Go/TEMPL frontend app
* The final board schematic, designed using KiCad
* Various 3D renderings of the PCB generated using KiCad

You can view the live production app on: https://duthweatherstation.azurewebsites.net/
