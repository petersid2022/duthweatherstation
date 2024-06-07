// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package web

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Chart() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav id=\"navbar\" class=\"sticky top-0 z-50 flex flex-col sm:flex-row justify-between items-center bg-sky-600 text-white px-4 py-2.5 mb-4 shadow-lg\"><h1 class=\"text-2xl font-bold sm:mb-0 mb-2\">DUTh Weather Station</h1><div class=\"flex flex-col sm:flex-row items-center sm:mb-0 mb-2\"><div class=\"flex flex-row items-center sm:w-auto w-full flex-grow sm:mb-0\"><div class=\"flex sm:flex-row flex-col items-center sm:w-auto w-full flex-grow sm:mb-0 mb-2\"><label for=\"start-date\" class=\"text-lg font-bold text-white mx-2\">Start date</label> <input id=\"start-date\" type=\"text\" class=\"text-center mx-2 block w-32 px-2 py-1 mr-4 bg-white border border-gray-300 rounded-md shadow-sm focus:ring-sky-500 focus:border-sky-500 text-lg text-black\" autocomplete=\"off\"></div><div class=\"flex sm:flex-row flex-col items-center sm:w-auto w-full flex-grow sm:mb-0 mb-2\"><label for=\"end-date\" class=\"text-lg font-bold text-white mx-2\">End date</label> <input id=\"end-date\" type=\"text\" class=\"text-center mx-2 block w-32 px-2 py-1 mr-4 bg-white border border-gray-300 rounded-md shadow-sm focus:ring-sky-500 focus:border-sky-500 text-lg text-black\" autocomplete=\"off\"></div></div><div class=\"flex flex-row justify-center items-center sm:w-auto w-full mt-2 sm:mt-0 sm:mb-0\"><button id=\"update-charts\" class=\"text-lg mx-2 px-4 py-1.5 font-bold text-white bg-sky-800 rounded-md shadow hover:bg-sky-900 focus:outline-none focus:ring-2 focus:ring-sky-500\">Update</button></div></div><div id=\"last-timestamp\" class=\"text-lg sm:pl-4\"><span>Last updated at: </span><span id=\"timestamp-value\"></span></div></nav><div id=\"chart\" class=\"grid grid-cols-1 md:grid-cols-2 gap-0 md:gap-8 justify-items-center\"><div id=\"chart-temperature\"></div><div id=\"chart-humidity\"></div><div id=\"chart-pressure\"></div><div id=\"chart-airquality\"></div></div><div class=\"py-2 text-center\"><p class=\"text-gray-600\">&copy; 2024 Peter Sideris and Fotis Mitsis</p><p class=\"text-gray-600\">This website is open <a href=\"https://github.com/petersid2022/duthweatherstation\" class=\"text-blue-500 underline hover:text-blue-700\" target=\"_blank\" rel=\"noopener noreferrer\">source</a></p></div><script>\n      const margin = { top: 30, right: 40, bottom: 40, left: 60 };\n      const width = 800 - margin.left - margin.right;\n      const height = 450 - margin.top - margin.bottom;\n      let allData = {};\n\n      const createChart = (data, xField, yField, svgSelector, legendText) => {\n        d3.select(svgSelector).select(\"svg\").remove();\n        const svg = d3.select(svgSelector).append(\"svg\")\n          .attr(\"width\", \"100%\")\n          .attr(\"height\", height + margin.top + margin.bottom)\n          .attr(\"viewBox\", `0 0 ${width + margin.left + margin.right} ${height + margin.top + margin.bottom}`)\n          .attr(\"preserveAspectRatio\", \"xMidYMid meet\")\n          .append(\"g\")\n          .attr(\"transform\", `translate(${margin.left},${margin.top})`);\n\n        const parseTime = d3.isoParse;\n\n        data.forEach(d => {\n          d.timestamp = parseTime(d.timestamp);\n          d[xField] = +d[xField];\n        });\n\n        const xExtent = d3.extent(data, d => d.timestamp);\n        const x = d3.scaleTime()\n          .domain(xExtent)\n          .range([0, width]);\n\n        let xAxisTicks = x.ticks(Math.max(width / 80, 2));\n        const maxXValue = x.domain()[1];\n        if (!xAxisTicks.includes(maxXValue)) {\n          xAxisTicks = xAxisTicks.concat(maxXValue);\n        }\n\n        const yExtent = d3.extent(data, d => d[yField]);\n        const yMargin = (yExtent[1] - yExtent[0]) * 0.05;\n        const y = d3.scaleLinear()\n          .domain([yExtent[0] - yMargin, yExtent[1] + yMargin]).nice()\n          .range([height, 0]);\n\n        let yAxisTicks = y.ticks(Math.max(height / 45, 2));\n        const maxYValue = y.domain()[1];\n        if (!yAxisTicks.includes(maxYValue)) {\n          yAxisTicks = yAxisTicks.concat(maxYValue);\n        }\n\n        const xAxis = g => g\n          .attr(\"class\", \"x-axis\")\n          .style(\"font-size\", \"14px\")\n          .attr(\"transform\", `translate(0,${height})`)\n          .call(d3.axisBottom(x).tickValues(xAxisTicks));\n\n        const yAxis = g => g\n          .attr(\"class\", \"y-axis\")\n          .style(\"font-size\", \"14px\")\n          .call(d3.axisLeft(y).tickValues(yAxisTicks));\n\n        const clip = svg.append(\"defs\").append(\"svg:clipPath\")\n          .attr(\"id\", \"clip\")\n          .append(\"svg:rect\")\n          .attr(\"width\", width)\n          .attr(\"height\", height)\n          .attr(\"x\", 0)\n          .attr(\"y\", 0);\n\n        const brush = d3.brushX()\n          .extent([[0, 0], [width, height]])\n          .on(\"end\", updateChart);\n\n        svg.append(\"g\")\n          .attr(\"class\", \"brush\")\n          .call(brush);\n\n        svg.selectAll(\"line.horizontalGrid\")\n          .data(yAxisTicks)\n          .enter()\n          .append(\"line\")\n          .attr(\"class\", \"horizontalGrid\")\n          .attr(\"x1\", 0)\n          .attr(\"x2\", width)\n          .attr(\"y1\", d => y(d))\n          .attr(\"y2\", d => y(d))\n          .attr(\"fill\", \"none\")\n          .attr(\"shape-rendering\", \"crispEdges\")\n          .attr(\"stroke\", \"lightgray\")\n          .attr(\"stroke-width\", \"1px\");\n\n        svg.selectAll(\"line.verticalGrid\")\n          .data(xAxisTicks)\n          .enter()\n          .append(\"line\")\n          .attr(\"class\", \"verticalGrid\")\n          .attr(\"x1\", d => x(d))\n          .attr(\"x2\", d => x(d))\n          .attr(\"y1\", 0)\n          .attr(\"y2\", height)\n          .attr(\"fill\", \"none\")\n          .attr(\"shape-rendering\", \"crispEdges\")\n          .attr(\"stroke\", \"lightgray\")\n          .attr(\"stroke-width\", \"1px\");\n\n        const line = svg.append(\"path\")\n          .datum(data)\n          .attr(\"class\", \"line\")\n          .attr(\"clip-path\", \"url(#clip)\")\n          .attr(\"fill\", \"none\")\n          .attr(\"stroke\", \"steelblue\")\n          .attr(\"stroke-width\", 1.5)\n          .attr(\"d\", d3.line()\n            .curve(d3.curveBasis)\n            .x(d => x(d.timestamp))\n            .y(d => y(d[yField])));\n\n        let idleTimeout;\n        function idled() { idleTimeout = null; }\n\n        function updateChart(event) {\n          const extent = event.selection;\n          if (!extent) {\n            if (!idleTimeout) return idleTimeout = setTimeout(idled, 350);\n            x.domain(d3.extent(data, d => d.timestamp));\n          } else {\n            x.domain([x.invert(extent[0]), x.invert(extent[1])]);\n            svg.select(\".brush\").call(brush.move, null);\n          }\n\n          svg.select(\".x-axis\").transition().duration(1000).call(d3.axisBottom(x));\n          line\n            .datum(data)\n            .transition()\n            .duration(1000)\n            .attr(\"d\", d3.line()\n            .curve(d3.curveBasis)\n            .x(d => x(d.timestamp))\n            .y(d => y(d[yField])));\n        }\n\n        svg.on(\"dblclick\", function() {\n          x.domain(d3.extent(data, d => d.timestamp));\n          svg.select(\".x-axis\").transition().call(d3.axisBottom(x));\n          line\n            .datum(data)\n            .transition()\n            .attr(\"d\", d3.line()\n            .curve(d3.curveBasis)\n            .x(d => x(d.timestamp))\n            .y(d => y(d[yField])));\n        });\n\n        svg.append(\"g\")\n          .call(xAxis);\n\n        svg.append(\"g\")\n          .call(yAxis);\n\n        const lastValue = data[0][yField];\n\n        const legend = svg.append(\"g\")\n          .attr(\"class\", \"legend\")\n          .attr(\"transform\", `translate(${width / 2},${margin.top})`);\n\n        legend.append(\"rect\")\n          .attr(\"x\", -70)\n          .attr(\"y\", -52)\n          .attr(\"width\", 14)\n          .attr(\"height\", 14)\n          .attr(\"fill\", \"steelblue\")\n          .attr(\"stroke\", \"steelblue\");\n\n        legend.append(\"text\")\n          .attr(\"x\", -50)\n          .attr(\"y\", -45)\n          .attr(\"dy\", \".35em\")\n          .style(\"text-anchor\", \"start\")\n          .style(\"font-size\", \"18px\")\n          .text(`${legendText}: ${lastValue}`);\n      };\n\n      const fetchData = (url) => {\n        return fetch(url, {\n          method: 'GET',\n          headers: {\n            'Content-Type': 'application/json',\n          }\n        }).then(response => response.json());\n      };\n\n      const filterDataByDate = (data, startDate, endDate) => {\n        const parseTime = d3.isoParse;\n        return data.filter(d => {\n         const date = parseTime(d.timestamp);\n          if (startDate.getTime() === endDate.getTime()) {\n            return date.toDateString() === startDate.toDateString();\n          } else {\n            return date >= startDate && date <= endDate;\n          }\n        });\n      };\n\n      const initializeDatePickers = (data) => {\n        const parseTime = d3.isoParse;\n        const dateExtent = d3.extent(data, d => parseTime(d.timestamp));\n        \n        const startDatePicker = flatpickr(\"#start-date\", {\n          defaultDate: dateExtent[0],\n          minDate: dateExtent[0],\n          maxDate: dateExtent[1],\n          disableMobile: \"true\",\n        });\n        \n        const endDatePicker = flatpickr(\"#end-date\", {\n          defaultDate: dateExtent[1],\n          minDate: dateExtent[0],\n          maxDate: dateExtent[1],\n          disableMobile: \"true\",\n        });\n\n        document.getElementById(\"start-date\").classList.remove(\"hidden\");\n        document.getElementById(\"end-date\").classList.remove(\"hidden\");\n      };\n\n      const updateCharts = () => {\n        const startDate = new Date(document.getElementById('start-date').value);\n        const endDate = new Date(document.getElementById('end-date').value);\n\n        if (startDate > endDate) {\n          console.error(\"Start date cannot be after end date\");\n          return;\n        }\n\n        const temperatureData = filterDataByDate(allData.temperature, startDate, endDate);\n        const humidityData = filterDataByDate(allData.humidity, startDate, endDate);\n        const pressureData = filterDataByDate(allData.pressure, startDate, endDate);\n        const airQualityData = filterDataByDate(allData.airQuality, startDate, endDate);\n\n        if (\n            temperatureData.length === 0 ||\n            humidityData.length === 0 ||\n            pressureData.length === 0 ||\n            airQualityData.length === 0\n           ) {\n          console.error(\"No data available within the selected date range\");\n          return;\n        }\n\n        createChart(temperatureData, 'timestamp', 'temperature', '#chart-temperature', 'Temperature');\n        createChart(humidityData, 'timestamp', 'humidity', '#chart-humidity', 'Humidity');\n        createChart(pressureData, 'timestamp', 'pressure', '#chart-pressure', 'Pressure');\n        createChart(airQualityData, 'timestamp', 'gas_level', '#chart-airquality', 'Air Quality');\n      };\n\n      document.getElementById('update-charts').addEventListener('click', updateCharts);\n\n      Promise.all([\n        fetchData('/api/data?sensor=dht11'),\n        fetchData('/api/data?sensor=bmp180'),\n        fetchData('/api/data?sensor=mq135')\n      ]).then(([dht11Data, bmp180Data, mq135Data]) => {\n        allData.temperature = dht11Data;\n        allData.humidity = dht11Data;\n        allData.pressure = bmp180Data;\n        allData.airQuality = mq135Data;\n            \n        initializeDatePickers(dht11Data.concat(bmp180Data, mq135Data));\n        updateCharts();\n            \n        const lastTimestamp = dht11Data[0].timestamp;\n        document.getElementById('timestamp-value').textContent = new Date(lastTimestamp).toLocaleString();\n      }).catch(error => {\n        console.error('Error fetching chart data:', error);\n      });\n    </script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
