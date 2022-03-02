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
    eT.SubElement(graph, "name").text = "ArgusWatcherProcessingTime" + event.title() + "Rate"
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
    eT.SubElement(graph, "description").text = "Rate of total " + event + " events processed"
    eT.SubElement(graph, "maxvalue")
    eT.SubElement(graph, "minvalue").text = "0"
    eT.SubElement(graph, "userparam1").text = "argus_watcher_processing_time"
    eT.SubElement(graph, "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    for quantile in quantiles:
        graph = eT.SubElement(datapoints, "datapoint")
        eT.SubElement(graph, "name").text = "ArgusWatcherProcessingTime" + event.title() + str(quantile)
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
        eT.SubElement(graph, "description").text = "Average time taken to complete " +str(quantile/100) + " quantile of "  + event + " event in last 1 minute"
        eT.SubElement(graph, "maxvalue")
        eT.SubElement(graph, "minvalue").text = "0"
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
    eT.SubElement(graph, "name").text = "ArgusWatcherEventLatency" + event.title() + "Rate"
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
    eT.SubElement(graph, "description").text = "Rate of total " + event + " events completed"
    eT.SubElement(graph, "maxvalue")
    eT.SubElement(graph, "minvalue").text = "0"
    eT.SubElement(graph, "userparam1").text = "argus_watcher_event_latency"
    eT.SubElement(graph, "userparam2").text = "event=\"" + event + "\",resource=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    for quantile in quantiles:
        graph = eT.SubElement(datapoints, "datapoint")
        eT.SubElement(graph, "name").text = "ArgusWatcherEventLatency" + event.title() + str(quantile)
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
        eT.SubElement(graph, "description").text = "Completed " +str(quantile/100) + " quantile of " + event + " events in last 1 minute"
        eT.SubElement(graph, "maxvalue")
        eT.SubElement(graph, "minvalue").text = "0"
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
    eT.SubElement(graph, "minvalue").text = "0.0"
    eT.SubElement(graph, "displayprio").text = "1"
    eT.SubElement(graph, "timescale").text = "1day"
    eT.SubElement(graph, "base1024").text = "false"
    gps = eT.SubElement(graph, "graphdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(gps, "graphdatapoint")
        eT.SubElement(dp, "name").text = "ArgusWatcherProcessingTime" + event.title() + str(quantile)
        eT.SubElement(dp, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + str(
            quantile)
    eT.SubElement(graph, "cf").text = "1"
    vgps = eT.SubElement(graph, "graphvirtualdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(vgps, "graphvirtualdatapoint")
        eT.SubElement(dp, "name").text = "ArgusWatcherProcessingTime" + event.title() + str(
            quantile) + "_ms"
        eT.SubElement(dp, "rpn").text = "ArgusWatcherProcessingTime" + event.title() + str(
            quantile) + "/1000000"
    gdatas = eT.SubElement(graph, "graphdatas")
    i = 0
    for quantile in quantiles:
        dp = eT.SubElement(gdatas, "graphdata")
        eT.SubElement(dp, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + str(
            quantile) + "_ms"
        eT.SubElement(dp, "legend").text = "P" + str(
            quantile)
        eT.SubElement(dp, "type").text = "1"
        eT.SubElement(dp, "color").text = colors[i]
        i = i + 1
        eT.SubElement(dp, "isvirtualdatapoint").text = "true"

for event in latency_events:
    graph = eT.SubElement(graphs, "graph")
    eT.SubElement(graph, "name").text = event.upper() + " Event Latency"
    eT.SubElement(graph, "title").text = event.upper() + " Event Latency"
    eT.SubElement(graph, "verticallabel").text = "milliseconds"
    eT.SubElement(graph, "rigid").text = "false"
    eT.SubElement(graph, "maxvalue").text = "NaN"
    eT.SubElement(graph, "minvalue").text = "0"
    eT.SubElement(graph, "displayprio").text = "100"
    eT.SubElement(graph, "timescale").text = "1day"
    eT.SubElement(graph, "base1024").text = "false"
    gps = eT.SubElement(graph, "graphdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(gps, "graphdatapoint")
        eT.SubElement(dp, "name").text = "ArgusWatcherEventLatency" + event.title() + str(quantile)
        eT.SubElement(dp, "datapointname").text = "ArgusWatcherEventLatency" + event.title() + str(
            quantile)
    eT.SubElement(graph, "cf").text = "1"
    vgps = eT.SubElement(graph, "graphvirtualdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(vgps, "graphvirtualdatapoint")
        eT.SubElement(dp, "name").text = "ArgusWatcherEventLatency" + event.title() + str(
            quantile) + "_ms"
        eT.SubElement(dp, "rpn").text = "ArgusWatcherEventLatency" + event.title() + str(
            quantile) + "/1000000"
    gdatas = eT.SubElement(graph, "graphdatas")
    i = 0
    for quantile in quantiles:
        dp = eT.SubElement(gdatas, "graphdata")
        eT.SubElement(dp, "datapointname").text = "ArgusWatcherEventLatency" + event.title() + str(
            quantile) + "_ms"
        eT.SubElement(dp, "legend").text = "P" + str(
            quantile)
        eT.SubElement(dp, "type").text = "1"
        eT.SubElement(dp, "color").text = colors[i]
        i = i + 1
        eT.SubElement(dp, "isvirtualdatapoint").text = "true"

ographs = tree.find(".//entry/overviewgraphs")

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 ADD Events Processing Time"
eT.SubElement(ograph, "title").text = "P99 ADD Events Processing Time"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "0.0"
eT.SubElement(ograph, "displayprio").text = "1"
eT.SubElement(ograph, "timescale").text = "1hour"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in events:
    if event == 'add':
        odp = eT.SubElement(odps, "overviewgraphdatapoint")
        eT.SubElement(odp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "cf").text = "1"
        eT.SubElement(odp, "aggregateMethod").text = "average"

        ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
        eT.SubElement(ovdp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(ovdp, "rpn").text = "ArgusWatcherProcessingTime" + event.title() + "99/1000000"

        line = eT.SubElement(olines, "overviewgraphline")
        eT.SubElement(line, "type").text = "1"
        eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
        eT.SubElement(line, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(line, "isvirtualdatapoint").text = "true"
        eT.SubElement(line, "color").text = "silver"

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 UPDATE Events Processing Time"
eT.SubElement(ograph, "title").text = "P99 UPDATE Events Processing Time"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "0.0"
eT.SubElement(ograph, "displayprio").text = "1"
eT.SubElement(ograph, "timescale").text = "1hour"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in events:
    if event == 'update':
        odp = eT.SubElement(odps, "overviewgraphdatapoint")
        eT.SubElement(odp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "cf").text = "1"
        eT.SubElement(odp, "aggregateMethod").text = "average"

        ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
        eT.SubElement(ovdp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(ovdp, "rpn").text = "ArgusWatcherProcessingTime" + event.title() + "99/1000000"

        line = eT.SubElement(olines, "overviewgraphline")
        eT.SubElement(line, "type").text = "1"
        eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
        eT.SubElement(line, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(line, "isvirtualdatapoint").text = "true"
        eT.SubElement(line, "color").text = "silver"

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 DELETE Events Processing Time"
eT.SubElement(ograph, "title").text = "P99 DELETE Events Processing Time"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "0.0"
eT.SubElement(ograph, "displayprio").text = "1"
eT.SubElement(ograph, "timescale").text = "1hour"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in events:
    if event == 'delete':
        odp = eT.SubElement(odps, "overviewgraphdatapoint")
        eT.SubElement(odp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99"
        eT.SubElement(odp, "cf").text = "1"
        eT.SubElement(odp, "aggregateMethod").text = "average"

        ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
        eT.SubElement(ovdp, "name").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(ovdp, "rpn").text = "ArgusWatcherProcessingTime" + event.title() + "99/1000000"

        line = eT.SubElement(olines, "overviewgraphline")
        eT.SubElement(line, "type").text = "1"
        eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
        eT.SubElement(line, "datapointname").text = "ArgusWatcherProcessingTime" + event.title() + "99_ms"
        eT.SubElement(line, "isvirtualdatapoint").text = "true"
        eT.SubElement(line, "color").text = "silver"

ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "P99 Event Listener Latency"
eT.SubElement(ograph, "title").text = "P99 Event Listener Latency"
eT.SubElement(ograph, "verticallabel").text = "milliseconds"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "0.0"
eT.SubElement(ograph, "displayprio").text = "100"
eT.SubElement(ograph, "timescale").text = "1day"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"

odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")

for event in latency_events:
    odp = eT.SubElement(odps, "overviewgraphdatapoint")
    eT.SubElement(odp, "name").text = "ArgusWatcherEventLatency" + event.title() + "99"
    eT.SubElement(odp, "datapointname").text = "ArgusWatcherEventLatency" + event.title() + "99"
    eT.SubElement(odp, "cf").text = "1"
    eT.SubElement(odp, "aggregateMethod").text = "average"

    ovdp = eT.SubElement(ovdps, "overviewgraphvirtualdatapoint")
    eT.SubElement(ovdp, "name").text = "ArgusWatcherEventLatency" + event.title() + "99_ms"
    eT.SubElement(ovdp, "rpn").text = "ArgusWatcherEventLatency" + event.title() + "99/1000000"

    line = eT.SubElement(olines, "overviewgraphline")
    eT.SubElement(line, "type").text = "1"
    eT.SubElement(line, "legend").text = event.upper() + " ##INSTANCE##"
    eT.SubElement(line, "datapointname").text = "ArgusWatcherEventLatency" + event.title() + "99_ms"
    eT.SubElement(line, "isvirtualdatapoint").text = "true"
    eT.SubElement(line, "color").text = "silver"

tree.write("argus_watcher_gen.xml")
