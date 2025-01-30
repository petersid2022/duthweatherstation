# duthweatherstation

This repo contains a proposed implementation of an Integrated Automatic Weather Station, centered around the FireBeetle 2 ESP32-E microcontroller, that's based on the ESP-WROOM-32E, chosen for its ultra-low power consumption, on-board charging circuit and compatibility with a wide range of sensors. The sensor array includes a DHT11 for humidity, a BMP180 for temperature and barometric pressure, and an MQ135 for air quality monitoring.

Included in this repo:
* The ESP32 firmware
* The Go microservice that handles the data stream from the ESP32 and various SQL transactions with the MySQL database
* Our Go/TEMPL frontend app
* The final board schematic, designed using KiCad
* Various 3D renderings of the PCB generated using KiCad
* The paper that describes in detail our implementation

~~You can view the live production app on: https://duthweatherstation.azurewebsites.net/~~

![PXL_20240612_170053991](https://github.com/user-attachments/assets/9658f100-8c8f-493f-9273-241401b907ba)
![PXL_20240612_170105711](https://github.com/user-attachments/assets/78f0174d-7ca2-4f56-a04f-34e039f4a54b)
![image](https://github.com/user-attachments/assets/43d6b4a0-cb35-49cb-b132-8db5b7945403)

## License
This project is licensed under the MIT License.

(The MIT License)
Copyright (c) 2024 Peter Sideris petersid2022@gmail.com
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
