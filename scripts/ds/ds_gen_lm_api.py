import json
import xml.etree.ElementTree as eT

filehandler = open("argus_lm_api_tmpl.xml", "r")
raw_data = eT.parse(filehandler)
data_root = raw_data.getroot()
filehandler.close()

jsonfile = open("config.json", "r")
jsonNode = json.load(jsonfile)
jsonfile.close()

quantiles = [25, 50, 75, 90, 99]
colors = ['gray', 'purple', 'yellow', 'orange', 'red']

tree = eT.ElementTree(data_root)
datapoints = tree.find(".//entry/datapoints")

for graph in datapoints:
    datapoints.remove(graph)
for api in jsonNode['apis']:
    graph = eT.SubElement(datapoints, "datapoint")
    eT.SubElement(graph, "name").text = "argus_lm_api_response_time_" + (
        api['verb'] if api['verb'] != "*" else "all") + "_" + str(
        api['status_code']) + "_rate"
    eT.SubElement(graph, "dataType").text = str(7)
    eT.SubElement(graph, "type").text = str(1)
    eT.SubElement(graph, "postprocessormethod").text = "regex"
    eT.SubElement(graph, "postprocessorparam").text = "code=\"" + str(
        api['status_code']) + "\",method=\"" + api['verb'].upper() + "\",url=\"##WILDVALUE##\""
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
    eT.SubElement(graph, "userparam1").text = "argus_lm_api_response_time_count"
    # et.SubElement(graph, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + api[
    #     'verb'].upper() + "\",url=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + (api[
                                                                                                          'verb'].upper() if
                                                                                                      api[
                                                                                                          'verb'] != "*" else "all") + "\",url=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "sum" if api['verb'] == "*" else "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    graph = eT.SubElement(datapoints, "datapoint")
    eT.SubElement(graph, "name").text = "argus_lm_api_response_time_" + (
        api['verb'] if api['verb'] != "*" else "all") + "_" + str(
        api['status_code']) + "_sum"
    eT.SubElement(graph, "dataType").text = str(7)
    eT.SubElement(graph, "type").text = str(2)
    eT.SubElement(graph, "postprocessormethod").text = "regex"
    eT.SubElement(graph, "postprocessorparam").text = "code=\"" + str(
        api['status_code']) + "\",method=\"" + api['verb'].upper() + "\",url=\"##WILDVALUE##\""
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
    eT.SubElement(graph, "userparam1").text = "argus_lm_api_response_time_sum"
    # et.SubElement(graph, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + api[
    #     'verb'].upper() + "\",url=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + (api[
                                                                                                          'verb'].upper() if
                                                                                                      api[
                                                                                                          'verb'] != "*" else "all") + "\",url=\"##WILDVALUE##\""
    eT.SubElement(graph, "userparam3").text = "sum" if api['verb'] == "*" else "none"
    eT.SubElement(graph, "iscomposite").text = "false"
    eT.SubElement(graph, "rpn")
    eT.SubElement(graph, "alertTransitionIval").text = "0"
    eT.SubElement(graph, "alertClearTransitionIval").text = "0"

    for quantile in quantiles:
        graph = eT.SubElement(datapoints, "datapoint")
        eT.SubElement(graph, "name").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(
            api['status_code']) + "_" + str(quantile)
        eT.SubElement(graph, "dataType").text = str(7)
        eT.SubElement(graph, "type").text = str(2)
        eT.SubElement(graph, "postprocessormethod").text = "regex"
        eT.SubElement(graph, "postprocessorparam").text = "code=\"" + str(
            api['status_code']) + "\",method=\"" + (api['verb'].upper() if api[
                                                                               'verb'] != "*" else "all") + "\",url=\"##WILDVALUE##\",quantile=\"" + str(
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
        eT.SubElement(graph, "userparam1").text = "argus_lm_api_response_time"
        eT.SubElement(graph, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + (
            api['verb'].upper() if api[
                                       'verb'] != "*" else "all") + "\",url=\"##WILDVALUE##\",quantile=\"" + str(
            quantile / 100) + "\""
        eT.SubElement(graph, "userparam3").text = "avg" if api['verb'] == "*" else "none"
        eT.SubElement(graph, "iscomposite").text = "false"
        eT.SubElement(graph, "rpn")
        eT.SubElement(graph, "alertTransitionIval").text = "0"
        eT.SubElement(graph, "alertClearTransitionIval").text = "0"

        # use virtual datapoint wherever necessary
        # ## ms
        # dp = et.SubElement(child, "datapoint")
        # et.SubElement(dp, "name").text = "argus_lm_api_response_time_" + (
        #     api['verb'] if api['verb'] != "*" else "all") + "_" + str(
        #     api['status_code']) + "_" + str(quantile) + "_ms"
        # et.SubElement(dp, "dataType").text = str(7)
        # et.SubElement(dp, "type").text = str(2)
        # et.SubElement(dp, "postprocessormethod").text = "expression"
        # et.SubElement(dp, "postprocessorparam").text = "argus_lm_api_response_time_" + (
        #     api['verb'] if api['verb'] != "*" else "all") + "_" + str(
        #     api['status_code']) + "_" + str(quantile) + "/1000000"
        # et.SubElement(dp, "usevalue")
        # et.SubElement(dp, "alertexpr")
        # et.SubElement(dp, "alertmissing").text = str(1)
        # et.SubElement(dp, "alertsubject")
        # et.SubElement(dp, "alertbody")
        # et.SubElement(dp, "enableanomalyalertsuppression")
        # et.SubElement(dp, "adadvsettingenabled").text = "false"
        # et.SubElement(dp, "warnadadvsetting")
        # et.SubElement(dp, "erroradadvsetting")
        # et.SubElement(dp, "criticaladadvsetting")
        # et.SubElement(dp, "description")
        # et.SubElement(dp, "maxvalue")
        # et.SubElement(dp, "minvalue")
        # et.SubElement(dp, "userparam1").text = "argus_lm_api_response_time"
        # et.SubElement(dp, "userparam2").text = "code=\"" + str(api['status_code']) + "\",method=\"" + api[
        #     'verb'].upper() + "\",url=\"##WILDVALUE##\",quantile=\"" + str(
        #     quantile / 100) + "\""
        # et.SubElement(dp, "userparam3").text = "none"
        # et.SubElement(dp, "iscomposite").text = "false"
        # et.SubElement(dp, "rpn")
        # et.SubElement(dp, "alertTransitionIval").text = "0"
        # et.SubElement(dp, "alertClearTransitionIval").text = "0"

graphs = tree.find(".//entry/graphs")

for graph in graphs:
    graphs.remove(graph)
for api in jsonNode['graphs']:
    graph = eT.SubElement(graphs, "graph")
    # eT.SubElement(graph, "name").text = (api['verb'].upper() if api['verb'] != "*" else "ALL") + " " + str(
    #     api['status_code'])
    # eT.SubElement(graph, "title").text = (api['verb'].upper() if api['verb'] != "*" else "ALL") + " " + str(
    #     api['status_code'])

    eT.SubElement(graph, "name").text = api['title']
    eT.SubElement(graph, "title").text = api['title']

    eT.SubElement(graph, "verticallabel").text = "milliseconds"
    eT.SubElement(graph, "rigid").text = "false"
    eT.SubElement(graph, "maxvalue").text = "NaN"
    eT.SubElement(graph, "minvalue").text = "NaN"
    eT.SubElement(graph, "displayprio").text = api['disp_prio']
    eT.SubElement(graph, "timescale").text = "1day"
    eT.SubElement(graph, "base1024").text = "false"
    gps = eT.SubElement(graph, "graphdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(gps, "graphdatapoint")
        eT.SubElement(dp, "name").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(
            api['status_code']) + "_" + str(quantile)
        eT.SubElement(dp, "datapointname").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(
            api['status_code']) + "_" + str(quantile)
        eT.SubElement(dp, "cf").text = "1"
    vgps = eT.SubElement(graph, "graphvirtualdatapoints")
    for quantile in quantiles:
        dp = eT.SubElement(vgps, "graphvirtualdatapoint")
        eT.SubElement(dp, "name").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(api[
                                                                          'status_code']) + "_" + str(
            quantile) + "_ms"
        eT.SubElement(dp, "rpn").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(api[
                                                                          'status_code']) + "_" + str(
            quantile) + "/1000000"
    gdatas = eT.SubElement(graph, "graphdatas")
    i = 0
    for quantile in quantiles:
        dp = eT.SubElement(gdatas, "graphdata")
        eT.SubElement(dp, "datapointname").text = "argus_lm_api_response_time_" + (
            api['verb'] if api['verb'] != "*" else "all") + "_" + str(
            api[
                'status_code']) + "_" + str(quantile) + "_ms"
        eT.SubElement(dp, "legend").text = (
            api['verb'] if api['verb'] != "*" else "all").upper() + " " + str(api[
                                                                          'status_code']) + " p" + str(
            quantile)
        eT.SubElement(dp, "type").text = "1"
        eT.SubElement(dp, "color").text = colors[i]
        i = i + 1
        eT.SubElement(dp, "isvirtualdatapoint").text = "true"

ographs = tree.find(".//entry/overviewgraphs")
ograph = eT.SubElement(ographs, "overviewgraph")
eT.SubElement(ograph, "name").text = "API Requests HTTP Stats"
eT.SubElement(ograph, "title").text = "API Requests HTTP Stats"
eT.SubElement(ograph, "verticallabel").text = "requests rate"
eT.SubElement(ograph, "rigid").text = "false"
eT.SubElement(ograph, "maxvalue").text = "NaN"
eT.SubElement(ograph, "minvalue").text = "0"
eT.SubElement(ograph, "displayprio").text = "1"
eT.SubElement(ograph, "timescale").text = "1day"
eT.SubElement(ograph, "base1024").text = "false"
eT.SubElement(ograph, "aggregated").text = "false"


odps = eT.SubElement(ograph, "datapoints")
ovdps = eT.SubElement(ograph, "virtualdatapoints")
olines = eT.SubElement(ograph, "lines")
for dp in datapoints:
    for c in dp:
        if c.tag == "name" and c.text.startswith("argus_lm_api_response_time_") and c.text.endswith("_rate"):
            name = c.text[len("argus_lm_api_response_time_"):]
            name = name[:-len("_rate")]
            if not name.startswith("all"):
                odp = eT.SubElement(odps, "overviewgraphdatapoint")
                eT.SubElement(odp, "name").text = c.text
                eT.SubElement(odp, "datapointname").text = c.text
                eT.SubElement(odp, "cf").text = "1"
                eT.SubElement(odp, "aggregateMethod").text = "avg"

                oline = eT.SubElement(olines, "overviewgraphline")
                eT.SubElement(oline, "type").text = "1"
                eT.SubElement(oline, "datapointname").text = c.text
                eT.SubElement(oline, "isvirtualdatapoint").text = "false"
                eT.SubElement(oline, "color").text = "silver"
                eT.SubElement(oline, "legend").text = name.replace("_", " ").upper() + " ##INSTANCE##"

tree.write("argus_lm_api_gen.xml")
