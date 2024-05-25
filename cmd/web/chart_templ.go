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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"chart\"></div><script>\n\t\t\tconst margin = { top: 20, right: 30, bottom: 40, left: 40 };\n\t\t\tconst width = 500 - margin.left - margin.right;\n\t\t\tconst height = 300 - margin.top - margin.bottom;\n\n\t\t\tconst createChart = (data, xField, yField, svgSelector, legendText) => {\n\t\t\t\tconst svg = d3.select(svgSelector).append(\"svg\")\n\t\t\t\t\t.attr(\"width\", width + margin.left + margin.right)\n\t\t\t\t\t.attr(\"height\", height + margin.top + margin.bottom)\n\t\t\t\t\t.append(\"g\")\n\t\t\t\t\t.attr(\"transform\", `translate(${margin.left},${margin.top})`);\n\n\t\t\t\tconst parseTime = d3.isoParse;\n\n\t\t\t\tdata.forEach(d => {\n\t\t\t\t\td.timestamp = parseTime(d.timestamp);\n\t\t\t\t\td[xField] = +d[xField];\n\t\t\t\t});\n\n\t\t\t\tconst x = d3.scaleTime()\n\t\t\t\t\t.domain(d3.extent(data, d => d.timestamp))\n\t\t\t\t\t.range([0, width]);\n\n\t\t\t\tconst y = d3.scaleLinear()\n\t\t\t\t\t.domain([0, d3.max(data, d => d[yField])]).nice()\n\t\t\t\t\t.range([height, 0]);\n\n\t\t\t\tconst xAxis = g => g\n\t\t\t\t\t.attr(\"transform\", `translate(0,${height})`)\n\t\t\t\t\t.call(d3.axisBottom(x).tickSizeOuter(0));\n\n\t\t\t\tconst yAxis = g => g\n\t\t\t\t\t.call(d3.axisLeft(y));\n\n\t\t\t\tsvg.append(\"path\")\n\t\t\t\t\t.datum(data)\n\t\t\t\t\t.attr(\"fill\", \"none\")\n\t\t\t\t\t.attr(\"stroke\", \"steelblue\")\n\t\t\t\t\t.attr(\"stroke-width\", 1.5)\n\t\t\t\t\t.attr(\"d\", d3.line()\n\t\t\t\t\t\t.x(d => x(d.timestamp))\n\t\t\t\t\t\t.y(d => y(d[yField]))\n\t\t\t\t\t);\n\n\t\t\t\tsvg.append(\"g\")\n\t\t\t\t\t.call(xAxis);\n\n\t\t\t\tsvg.append(\"g\")\n\t\t\t\t\t.call(yAxis);\n\n\t\t\t\t// Add legend\n\t\t\t\tsvg.append(\"g\")\n\t\t\t\t\t.attr(\"class\", \"legend\")\n\t\t\t\t\t.attr(\"transform\", `translate(${width - 50},${margin.top})`)\n\t\t\t\t\t.append(\"text\")\n\t\t\t\t\t.attr(\"x\", 0)\n\t\t\t\t\t.attr(\"y\", 0)\n\t\t\t\t\t.attr(\"dy\", \".35em\")\n\t\t\t\t\t.style(\"text-anchor\", \"start\")\n\t\t\t\t\t.text(legendText);\n\n\t\t\t\tsvg.select(\".legend\")\n\t\t\t\t\t.append(\"rect\")\n\t\t\t\t\t.attr(\"x\", -20)\n\t\t\t\t\t.attr(\"width\", 12)\n\t\t\t\t\t.attr(\"height\", 12)\n\t\t\t\t\t.attr(\"fill\", \"steelblue\")\n\t\t\t\t\t.attr(\"stroke\", \"steelblue\");\n\t\t\t};\n\n\t\t\tconst fetchAndCreateChart = (url, xField, yField, svgSelector, legendText) => {\n\t\t\t\tfetch(url)\n\t\t\t\t\t.then(response => response.json())\n\t\t\t\t\t.then(data => createChart(data, xField, yField, svgSelector, legendText))\n\t\t\t\t\t.catch(error => {\n\t\t\t\t\t\tdocument.querySelector(svgSelector).innerHTML = \"Error loading chart.\";\n\t\t\t\t\t\tconsole.error('Error fetching chart data:', error);\n\t\t\t\t\t});\n\t\t\t};\n\n\t\t\tfetchAndCreateChart('/api/data?sensor=dht11', 'timestamp', 'temperature', '#chart', 'Temperature');\n\t\t\tfetchAndCreateChart('/api/data?sensor=dht11', 'timestamp', 'humidity', '#chart', 'Humidity');\n\t\t\tfetchAndCreateChart('/api/data?sensor=bmp180', 'timestamp', 'pressure', '#chart', 'Pressure');\n\t\t\tfetchAndCreateChart('/api/data?sensor=mq135', 'timestamp', 'gas_level', '#chart', 'Gas Level');\n\t\t</script>")
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