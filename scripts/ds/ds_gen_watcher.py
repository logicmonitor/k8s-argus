import xml.etree.ElementTree as eT

filehandler = open("argus_watcher_tmpl.xml", "r")
raw_data = eT.parse(filehandler)
data_root = raw_data.getroot()
filehandler.close()

tree = eT.ElementTree(data_root)

datapoints = tree.find(".//entry/datapoints")

events = ['add', 'update', 'delete']
latency_events = ['delete']
quantiles = [25, 50, 75, 90, 99]
colors = ['gray', 'purple', 'yellow', 'orange', 'red']

for event in events:
    graph = eT.SubElement(datapoints, "datapoint")
    eT.SubElement(graph, "name").text = "argus_watcher_processing_time_" + event + "_rate"
    eT.SubElement(graph, "dataType").text = str(7)
    eT.SubElement(graph, "type").text = str(1)
    eT.SubElement(graph, "postprocessormethod").text = "regex"
    eT.SubElement(graph, "postprocessorparam").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "usevalue").text = "body"
    eT.SubElement(graph, "alertexpr")
    eT.SubElement(graph, "alertmissing").text = str(1)
    eT.SubElement(graph, "alertsubject")
    eT.SubElement(graph, "alertbody")
    eT.SubElement(graph, "enableanomalyalertsuppression")
    eT.SubElement(graph, "adadvsettingenabled").text = "false"
    eT.SubElement(graph, "warnadadvsetting")
    eT.SubElement(graph, "erroradadvsetting")
    eT.SubElement(graph, "criticaladadvsetting")
    eT.SubElement(graph, "description")
    eT.SubElement(graph, "maxvalue")
    eT.SubElement(graph, "minvalue")
    eT.SubElement(graph, "userparam1").text = "argus_watcher_processing_time"
    eT.SubElement(graph, "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    for quantile in quantiles:
        graph = eT.SubElement(datapoints, "datapoint")
        eT.SubElement(graph, "name").text = "argus_watcher_processing_time_" + event + "_" + str(quantile)
        eT.SubElement(graph, "dataType").text = str(7)
        eT.SubElement(graph, "type").text = str(2)
        eT.SubElement(graph, "postprocessormethod").text = "regex"
        eT.SubElement(graph,
                      "postprocessorparam").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\",quantile=\"" + str(
            quantile / 100) + "\""
        eT.SubElement(graph, "usevalue").text = "body"
        eT.SubElement(graph, "alertexpr")
        eT.SubElement(graph, "alertmissing").text = str(1)
        eT.SubElement(graph, "alertsubject")
        eT.SubElement(graph, "alertbody")
        eT.SubElement(graph, "enableanomalyalertsuppression")
        eT.SubElement(graph, "adadvsettingenabled").text = "false"
        eT.SubElement(graph, "warnadadvsetting")
        eT.SubElement(graph, "erroradadvsetting")
        eT.SubElement(graph, "criticaladadvsetting")
        eT.SubElement(graph, "description")
        eT.SubElement(graph, "maxvalue")
        eT.SubElement(graph, "minvalue")
        eT.SubElement(graph, "userparam1").text = "argus_watcher_processing_time"
        eT.SubElement(graph,
                      "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\",quantile=\"" + str(
            quantile / 100) + "\""
        eT.SubElement(graph, "userparam3").text = "none"
        eT.SubElement(graph, "iscomposite").text = "false"
        eT.SubElement(graph, "rpn")
        eT.SubElement(graph, "alertTransitionIval").text = "0"
        eT.SubElement(graph, "alertClearTransitionIval").text = "0"

for event in latency_events:
    graph = eT.SubElement(datapoints, "datapoint")
    eT.SubElement(graph, "name").text = "argus_watcher_event_latency_" + event + "_rate"
    eT.SubElement(graph, "dataType").text = str(7)
    eT.SubElement(graph, "type").text = str(1)
    eT.SubElement(graph, "postprocessormethod").text = "regex"
    eT.SubElement(graph, "postprocessorparam").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "usevalue").text = "body"
    eT.SubElement(graph, "alertexpr")
    eT.SubElement(graph, "alertmissing").text = str(1)
    eT.SubElement(graph, "alertsubject")
    eT.SubElement(graph, "alertbody")
    eT.SubElement(graph, "enableanomalyalertsuppression")
    eT.SubElement(graph, "adadvsettingenabled").text = "false"
    eT.SubElement(graph, "warnadadvsetting")
    eT.SubElement(graph, "erroradadvsetting")
    eT.SubElement(graph, "criticaladadvsetting")
    eT.SubElement(graph, "description")
    eT.SubElement(graph, "maxvalue")
    eT.SubElement(graph, "minvalue")
    eT.SubElement(graph, "userparam1").text = "argus_watcher_event_latency"
    eT.SubElement(graph, "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    for quantile in quantiles:
        graph = eT.SubElement(datapoints, "datapoint")
        eT.SubElement(graph, "name").text = "argus_watcher_event_latency_" + event + "_" + str(quantile)
        eT.SubElement(graph, "dataType").text = str(7)
        eT.SubElement(graph, "type").text = str(2)
        eT.SubElement(graph, "postprocessormethod").text = "regex"
        eT.SubElement(graph,
                      "postprocessorparam").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\",quantile=\"" + str(
            quantile / 100) + "\""
        eT.SubElement(graph, "usevalue").text = "body"
        eT.SubElement(graph, "alertexpr")
        eT.SubElement(graph, "alertmissing").text = str(1)
        eT.SubElement(graph, "alertsubject")
        eT.SubElement(graph, "alertbody")
        eT.SubElement(graph, "enableanomalyalertsuppression")
        eT.SubElement(graph, "adadvsettingenabled").text = "false"
        eT.SubElement(graph, "warnadadvsetting")
        eT.SubElement(graph, "erroradadvsetting")
        eT.SubElement(graph, "criticaladadvsetting")
        eT.SubElement(graph, "description")
        eT.SubElement(graph, "maxvalue")
        eT.SubElement(graph, "minvalue")
        eT.SubElement(graph, "userparam1").text = "argus_watcher_event_latency"
        eT.SubElement(graph,
                      "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\",quantile=\"" + str(
            quantile / 100) + "\""
        eT.SubElement(graph, "userparam3").text = "none"
        eT.SubElement(graph, "iscomposite").text = "false"
        eT.SubElement(graph, "rpn")
        eT.SubElement(graph, "alertTransitionIval").text = "0"
        eT.SubElement(graph, "alertClearTransitionIval").text = "0"

graphs = tree.find(".//entry/graphs")

for graph in graphs:
    graphs.remove(graph)
for event in events:
    graph = eT.SubElement(graphs, "graph")
    eT.SubElement(graph, "name").text = event.upper() + " Processing Time"
    eT.SubElement(graph, "title").text = event.upper() + " Processing Time"
    eT.SubElement(graph, "verticallabel").text = "milliseconds"
    eT.SubElement(graph, "rigid").text = "false"
    eT.SubElement(graph, "maxvalue").text = "NaN"
    eT.SubElement(graph, "minvalue").text = "NaN"
    eT.SubElement(graph, "displayprio").text = "1"
    eT.SubElement(graph, "timescale").text = "1day"
    eT.SubElement(graph, "base1024").text = "false"
    gps = eT.SubElement(graph, "graphdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(gps, "graphdatapoint")
        eT.SubElement(dp, "name").text = "argus_watcher_processing_time_" + event + "_" + str(quantile)
        eT.SubElement(dp, "datapointname").text = "argus_watcher_processing_time_" + event + "_" + str(
            quantile)
    eT.SubElement(graph, "cf").text = "1"
    vgps = eT.SubElement(graph, "graphvirtualdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(vgps, "graphvirtualdatapoint")
        eT.SubElement(dp, "name").text = "argus_watcher_processing_time_" + event + "_" + str(
            quantile) + "_ms"
        eT.SubElement(dp, "rpn").text = "argus_watcher_processing_time_" + event + "_" + str(
            quantile) + "/1000000"
    gdatas = eT.SubElement(graph, "graphdatas")
    i = 0
    for quantile in quantiles:
        dp = eT.SubElement(gdatas, "graphdata")
        eT.SubElement(dp, "datapointname").text = "argus_watcher_processing_time_" + event + "_" + str(
            quantile) + "_ms"
        eT.SubElement(dp, "legend").text = "P" + str(
            quantile)
        eT.SubElement(dp, "type").text = "1"
        eT.SubElement(dp, "color").text = colors[i]
        i = i + 1
        eT.SubElement(dp, "isvirtualdatapoint").text = "true"

for event in latency_events:
    graph = eT.SubElement(graphs, "graph")
    eT.SubElement(graph, "name").text = event.upper() + " event latency"
    eT.SubElement(graph, "title").text = event.upper() + " event latency"
    eT.SubElement(graph, "verticallabel").text = "milliseconds"
    eT.SubElement(graph, "rigid").text = "false"
    eT.SubElement(graph, "maxvalue").text = "NaN"
    eT.SubElement(graph, "minvalue").text = "NaN"
    eT.SubElement(graph, "displayprio").text = "100"
    eT.SubElement(graph, "timescale").text = "1day"
    eT.SubElement(graph, "base1024").text = "false"
    gps = eT.SubElement(graph, "graphdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(gps, "graphdatapoint")
        eT.SubElement(dp, "name").text = "argus_watcher_event_latency_" + event + "_" + str(quantile)
        eT.SubElement(dp, "datapointname").text = "argus_watcher_event_latency_" + event + "_" + str(
            quantile)
    eT.SubElement(graph, "cf").text = "1"
    vgps = eT.SubElement(graph, "graphvirtualdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(vgps, "graphvirtualdatapoint")
        eT.SubElement(dp, "name").text = "argus_watcher_event_latency_" + event + "_" + str(
            quantile) + "_ms"
        eT.SubElement(dp, "rpn").text = "argus_watcher_event_latency_" + event + "_" + str(
            quantile) + "/1000000"
    gdatas = eT.SubElement(graph, "graphdatas")
    i = 0
    for quantile in quantiles:
        dp = eT.SubElement(gdatas, "graphdata")
        eT.SubElement(dp, "datapointname").text = "argus_watcher_event_latency_" + event + "_" + str(
            quantile) + "_ms"
        eT.SubElement(dp, "legend").text = "P" + str(
            quantile)
        eT.SubElement(dp, "type").text = "1"
        eT.SubElement(dp, "color").text = colors[i]
        i = i + 1
        eT.SubElement(dp, "isvirtualdatapoint").text = "true"

ographs = tree.find(".//entry/overviewgraphs")

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 Processing Time"
eT.SubElement(ograph, "title").text = "P99 Processing Time"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "NaN"
eT.SubElement(ograph, "displayprio").text = "1"
eT.SubElement(ograph, "timescale").text = "1day"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in events:
    odp = eT.SubElement(odps, "overviewgraphdatapoint")
    eT.SubElement(odp, "name").text = "argus_watcher_processing_time_" + event + "_99"
    eT.SubElement(odp, "datapointname").text = "argus_watcher_processing_time_" + event + "_99"
    eT.SubElement(odp, "cf").text = "1"
    eT.SubElement(odp, "aggregateMethod").text = "average"

    ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
    eT.SubElement(ovdp, "name").text = "argus_watcher_processing_time_" + event + "_99_ms"
    eT.SubElement(ovdp, "rpn").text = "argus_watcher_processing_time_" + event + "_99/1000000"

    line = eT.SubElement(olines, "overviewgraphline")
    eT.SubElement(line, "type").text = "1"
    eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
    eT.SubElement(line, "datapointname").text = "argus_watcher_processing_time_" + event + "_99_ms"
    eT.SubElement(line, "isvirtualdatapoint").text = "true"
    eT.SubElement(line, "color").text = "silver"

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 Event Listener Latency"
eT.SubElement(ograph, "title").text = "P99 Event Listener Latency"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "NaN"
eT.SubElement(ograph, "displayprio").text = "100"
eT.SubElement(ograph, "timescale").text = "1day"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in latency_events:
    odp = eT.SubElement(odps, "overviewgraphdatapoint")
    eT.SubElement(odp, "name").text = "argus_watcher_event_latency_" + event + "_99"
    eT.SubElement(odp, "datapointname").text = "argus_watcher_event_latency_" + event + "_99"
    eT.SubElement(odp, "cf").text = "1"
    eT.SubElement(odp, "aggregateMethod").text = "average"

    ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
    eT.SubElement(ovdp, "name").text = "argus_watcher_event_latency_" + event + "_99_ms"
    eT.SubElement(ovdp, "rpn").text = "argus_watcher_event_latency_" + event + "_99/1000000"

    line = eT.SubElement(olines, "overviewgraphline")
    eT.SubElement(line, "type").text = "1"
    eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
    eT.SubElement(line, "datapointname").text = "argus_watcher_event_latency_" + event + "_99_ms"
    eT.SubElement(line, "isvirtualdatapoint").text = "true"
    eT.SubElement(line, "color").text = "silver"

tree.write("argus_watcher_gen.xml")
