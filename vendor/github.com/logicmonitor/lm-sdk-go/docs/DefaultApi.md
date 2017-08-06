# \DefaultApi

All URIs are relative to *https://localhost/santaba/rest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AckAlertById**](DefaultApi.md#AckAlertById) | **Post** /alert/alerts/{id}/ack | ack alert by id
[**AckCollectorDownAlertById**](DefaultApi.md#AckCollectorDownAlertById) | **Post** /setting/collectors/{id}/ackdown | ack collector down alert
[**AddAdmin**](DefaultApi.md#AddAdmin) | **Post** /setting/admins | add admin
[**AddAlertNoteById**](DefaultApi.md#AddAlertNoteById) | **Post** /alert/alerts/{id}/note | add alert note
[**AddAlertRule**](DefaultApi.md#AddAlertRule) | **Post** /setting/alert/rules | add alert rule
[**AddApiTokenByAdminId**](DefaultApi.md#AddApiTokenByAdminId) | **Post** /setting/admins/{adminId}/apitokens | add apiToken by admin
[**AddCollector**](DefaultApi.md#AddCollector) | **Post** /setting/collectors | add collector
[**AddCollectorGroup**](DefaultApi.md#AddCollectorGroup) | **Post** /setting/collectors/groups | add collector group
[**AddDashboard**](DefaultApi.md#AddDashboard) | **Post** /dashboard/dashboards | add dashboard
[**AddDashboardGroup**](DefaultApi.md#AddDashboardGroup) | **Post** /dashboard/groups | add dashboard group
[**AddDevice**](DefaultApi.md#AddDevice) | **Post** /device/devices | add a new device
[**AddDeviceDatasourceInstance**](DefaultApi.md#AddDeviceDatasourceInstance) | **Post** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances | add device instance 
[**AddDeviceDatasourceInstanceGroup**](DefaultApi.md#AddDeviceDatasourceInstanceGroup) | **Post** /device/devices/{deviceId}/devicedatasources/{deviceDsId}/groups | add device datasource instance group 
[**AddDeviceGroup**](DefaultApi.md#AddDeviceGroup) | **Post** /device/groups | add device group
[**AddDeviceGroupProperty**](DefaultApi.md#AddDeviceGroupProperty) | **Post** /device/groups/{gid}/properties | add device group property
[**AddDeviceProperty**](DefaultApi.md#AddDeviceProperty) | **Post** /device/devices/{deviceId}/properties | add device property
[**AddRole**](DefaultApi.md#AddRole) | **Post** /setting/roles | add role
[**AddServiceGroup**](DefaultApi.md#AddServiceGroup) | **Post** /service/groups | add service group
[**DeleteAdminById**](DefaultApi.md#DeleteAdminById) | **Delete** /setting/admins/{id} | delete admin
[**DeleteAlertRuleById**](DefaultApi.md#DeleteAlertRuleById) | **Delete** /setting/alert/rules/{id} | delete alert rule
[**DeleteApiTokenById**](DefaultApi.md#DeleteApiTokenById) | **Delete** /setting/admins/{adminId}/apitokens/{apitokenId} | delete apiToken 
[**DeleteCollectorById**](DefaultApi.md#DeleteCollectorById) | **Delete** /setting/collectors/{id} | delete collector
[**DeleteCollectorGroupById**](DefaultApi.md#DeleteCollectorGroupById) | **Delete** /setting/collectors/groups/{id} | delete collector group
[**DeleteDashboardById**](DefaultApi.md#DeleteDashboardById) | **Delete** /dashboard/dashboards/{id} | delete dashboard
[**DeleteDashboardGroupById**](DefaultApi.md#DeleteDashboardGroupById) | **Delete** /dashboard/groups/{id} | delete dashboard group
[**DeleteDevice**](DefaultApi.md#DeleteDevice) | **Delete** /device/devices/{id} | delete a device
[**DeleteDeviceGroupById**](DefaultApi.md#DeleteDeviceGroupById) | **Delete** /device/groups/{id} | delete device group
[**DeleteDeviceGroupPropertyByName**](DefaultApi.md#DeleteDeviceGroupPropertyByName) | **Delete** /device/groups/{gid}/properties/{name} | delete device group property
[**DeleteDevicePropertyByName**](DefaultApi.md#DeleteDevicePropertyByName) | **Delete** /device/devices/{deviceId}/properties/{name} | delete device  property
[**DeleteRoleById**](DefaultApi.md#DeleteRoleById) | **Delete** /setting/roles/{id} | delete role
[**DeleteServiceGroupById**](DefaultApi.md#DeleteServiceGroupById) | **Delete** /service/groups/{id} | delete service group
[**GetAdminById**](DefaultApi.md#GetAdminById) | **Get** /setting/admins/{id} | get admin
[**GetAdminList**](DefaultApi.md#GetAdminList) | **Get** /setting/admins | get admin list
[**GetAlertById**](DefaultApi.md#GetAlertById) | **Get** /alert/alerts/{id} | get alert
[**GetAlertList**](DefaultApi.md#GetAlertList) | **Get** /alert/alerts | get alert list
[**GetAlertListByDeviceGroupId**](DefaultApi.md#GetAlertListByDeviceGroupId) | **Get** /device/groups/{id}/alerts | get device group alerts
[**GetAlertListByDeviceId**](DefaultApi.md#GetAlertListByDeviceId) | **Get** /device/devices/{id}/alerts | get alerts
[**GetAlertRuleById**](DefaultApi.md#GetAlertRuleById) | **Get** /setting/alert/rules/{id} | get alert rule by id
[**GetAlertRuleList**](DefaultApi.md#GetAlertRuleList) | **Get** /setting/alert/rules | get alert rule list
[**GetApiTokenListByAdminId**](DefaultApi.md#GetApiTokenListByAdminId) | **Get** /setting/admins/{adminId}/apitokens | get apiToken by admin
[**GetCollectorById**](DefaultApi.md#GetCollectorById) | **Get** /setting/collectors/{id} | get collector
[**GetCollectorGroupById**](DefaultApi.md#GetCollectorGroupById) | **Get** /setting/collectors/groups/{id} | get collector group
[**GetCollectorGroupList**](DefaultApi.md#GetCollectorGroupList) | **Get** /setting/collectors/groups | get collector group list
[**GetCollectorList**](DefaultApi.md#GetCollectorList) | **Get** /setting/collectors | get collector list
[**GetDashboardById**](DefaultApi.md#GetDashboardById) | **Get** /dashboard/dashboards/{id} | get dashboard
[**GetDashboardGroupById**](DefaultApi.md#GetDashboardGroupById) | **Get** /dashboard/groups/{id} | get dashboard group
[**GetDashboardGroupList**](DefaultApi.md#GetDashboardGroupList) | **Get** /dashboard/groups | get dashboard group list
[**GetDashboardList**](DefaultApi.md#GetDashboardList) | **Get** /dashboard/dashboards | get dashboard list
[**GetDeviceById**](DefaultApi.md#GetDeviceById) | **Get** /device/devices/{id} | get device by id
[**GetDeviceDatasourceById**](DefaultApi.md#GetDeviceDatasourceById) | **Get** /device/devices/{deviceId}/devicedatasources/{id} | get device datasource 
[**GetDeviceDatasourceDataById**](DefaultApi.md#GetDeviceDatasourceDataById) | **Get** /device/devices/{deviceId}/devicedatasources/{id}/data | get device datasource data 
[**GetDeviceDatasourceInstanceAlertSettingById**](DefaultApi.md#GetDeviceDatasourceInstanceAlertSettingById) | **Get** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{instanceId}/alertsettings/{id} | get device instance alert setting
[**GetDeviceDatasourceInstanceById**](DefaultApi.md#GetDeviceDatasourceInstanceById) | **Get** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{id} | get device instance 
[**GetDeviceDatasourceInstanceData**](DefaultApi.md#GetDeviceDatasourceInstanceData) | **Get** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{id}/data | get device instance data
[**GetDeviceDatasourceInstanceGraphData**](DefaultApi.md#GetDeviceDatasourceInstanceGraphData) | **Get** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{id}/graphs/{graphId}/data | get device instance graph data 
[**GetDeviceDatasourceInstanceGroupById**](DefaultApi.md#GetDeviceDatasourceInstanceGroupById) | **Get** /device/devices/{deviceId}/devicedatasources/{deviceDsId}/groups/{id} | get device datasource instance group 
[**GetDeviceDatasourceInstanceGroupOverviewGraphData**](DefaultApi.md#GetDeviceDatasourceInstanceGroupOverviewGraphData) | **Get** /device/devices/{deviceId}/devicedatasources/{deviceDsId}/groups/{dsigId}/graphs/{ographId}/data | get device instance group overview graph data 
[**GetDeviceDatasourceList**](DefaultApi.md#GetDeviceDatasourceList) | **Get** /device/devices/{deviceId}/devicedatasources | get device datasource list 
[**GetDeviceGroupById**](DefaultApi.md#GetDeviceGroupById) | **Get** /device/groups/{id} | get device group
[**GetDeviceGroupDatasourceAlertSetting**](DefaultApi.md#GetDeviceGroupDatasourceAlertSetting) | **Get** /device/groups/{deviceGroupId}/datasources/{dsId}/alertsettings | get device group datasource alert setting 
[**GetDeviceGroupDatasourceById**](DefaultApi.md#GetDeviceGroupDatasourceById) | **Get** /device/groups/{deviceGroupId}/datasources/{id} | get device group datasource
[**GetDeviceGroupDatasourceList**](DefaultApi.md#GetDeviceGroupDatasourceList) | **Get** /device/groups/{deviceGroupId}/datasources | get device group datasource list
[**GetDeviceGroupList**](DefaultApi.md#GetDeviceGroupList) | **Get** /device/groups | get device group list
[**GetDeviceGroupProperties**](DefaultApi.md#GetDeviceGroupProperties) | **Get** /device/groups/{gid}/properties | get device group properties
[**GetDeviceGroupPropertyByName**](DefaultApi.md#GetDeviceGroupPropertyByName) | **Get** /device/groups/{gid}/properties/{name} | get device group property by name
[**GetDeviceInstanceGraphDataOnlyByInstanceId**](DefaultApi.md#GetDeviceInstanceGraphDataOnlyByInstanceId) | **Get** /device/devicedatasourceinstances/{instanceId}/graphs/{graphId}/data | get device instance graphData
[**GetDeviceList**](DefaultApi.md#GetDeviceList) | **Get** /device/devices | get device list
[**GetDevicePropertiesList**](DefaultApi.md#GetDevicePropertiesList) | **Get** /device/devices/{deviceId}/properties | get device properties
[**GetDevicePropertyByName**](DefaultApi.md#GetDevicePropertyByName) | **Get** /device/devices/{deviceId}/properties/{name} | get device property by name
[**GetImmediateDeviceListByDeviceGroupId**](DefaultApi.md#GetImmediateDeviceListByDeviceGroupId) | **Get** /device/groups/{id}/devices | get immediate devices under group
[**GetRoleById**](DefaultApi.md#GetRoleById) | **Get** /setting/roles/{id} | get role by id
[**GetRoleList**](DefaultApi.md#GetRoleList) | **Get** /setting/roles | get role list
[**GetServiceGraphData**](DefaultApi.md#GetServiceGraphData) | **Get** /service/services/{serviceId}/checkpoints/{checkpointId}/graphs/{graphName}/data | get service graph data
[**GetServiceGroupById**](DefaultApi.md#GetServiceGroupById) | **Get** /service/groups/{id} | get service group
[**GetServiceGroupList**](DefaultApi.md#GetServiceGroupList) | **Get** /service/groups | get service group list
[**GetSiteMonitorCheckPointList**](DefaultApi.md#GetSiteMonitorCheckPointList) | **Get** /service/smcheckpoints | get site monitor checkpoint list
[**GetUnmonitoredDeviceList**](DefaultApi.md#GetUnmonitoredDeviceList) | **Get** /device/unmonitoreddevices | get unmonitored device list
[**InstallCollector**](DefaultApi.md#InstallCollector) | **Get** /setting/collectors/{collectorId}/installers/{osAndArch} | install collector
[**PatchDeviceById**](DefaultApi.md#PatchDeviceById) | **Patch** /device/devices/{id} | patch a device
[**PatchDeviceGroupById**](DefaultApi.md#PatchDeviceGroupById) | **Patch** /device/groups/{id} | patch device group
[**PatchServiceGroupById**](DefaultApi.md#PatchServiceGroupById) | **Patch** /service/groups/{id} | patch service group
[**ScheduleAutoDiscoveryByDeviceId**](DefaultApi.md#ScheduleAutoDiscoveryByDeviceId) | **Post** /device/devices/{id}/scheduleAutoDiscovery | schedule auto discovery for a host
[**UpdateAdminById**](DefaultApi.md#UpdateAdminById) | **Put** /setting/admins/{id} | update admin
[**UpdateAlertRuleById**](DefaultApi.md#UpdateAlertRuleById) | **Put** /setting/alert/rules/{id} | update alert rule
[**UpdateApiTokenByAdminId**](DefaultApi.md#UpdateApiTokenByAdminId) | **Put** /setting/admins/{adminId}/apitokens/{apitokenId} | update apiToken by admin
[**UpdateCollectorById**](DefaultApi.md#UpdateCollectorById) | **Put** /setting/collectors/{id} | update collector
[**UpdateCollectorGroupById**](DefaultApi.md#UpdateCollectorGroupById) | **Put** /setting/collectors/groups/{id} | update collector group
[**UpdateDashboardById**](DefaultApi.md#UpdateDashboardById) | **Put** /dashboard/dashboards/{id} | update dashboard
[**UpdateDashboardGroupById**](DefaultApi.md#UpdateDashboardGroupById) | **Put** /dashboard/groups/{id} | update dashboard group
[**UpdateDevice**](DefaultApi.md#UpdateDevice) | **Put** /device/devices/{id} | update a device
[**UpdateDeviceDatasourceInstanceAlertSettingById**](DefaultApi.md#UpdateDeviceDatasourceInstanceAlertSettingById) | **Put** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{instanceId}/alertsettings/{id} | update device instance alert setting
[**UpdateDeviceDatasourceInstanceById**](DefaultApi.md#UpdateDeviceDatasourceInstanceById) | **Put** /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{id} | update device instance 
[**UpdateDeviceDatasourceInstanceGroupById**](DefaultApi.md#UpdateDeviceDatasourceInstanceGroupById) | **Put** /device/devices/{deviceId}/devicedatasources/{deviceDsId}/groups/{id} | update device datasource instance group 
[**UpdateDeviceGroupById**](DefaultApi.md#UpdateDeviceGroupById) | **Put** /device/groups/{id} | update device group
[**UpdateDeviceGroupDatasourceAlertSetting**](DefaultApi.md#UpdateDeviceGroupDatasourceAlertSetting) | **Put** /device/groups/{deviceGroupId}/datasources/{dsId}/alertsettings | update device group datasource alert setting 
[**UpdateDeviceGroupPropertyByName**](DefaultApi.md#UpdateDeviceGroupPropertyByName) | **Put** /device/groups/{gid}/properties/{name} | update device group property
[**UpdateDevicePropertyByName**](DefaultApi.md#UpdateDevicePropertyByName) | **Put** /device/devices/{deviceId}/properties/{name} | update device  property
[**UpdateRoleById**](DefaultApi.md#UpdateRoleById) | **Put** /setting/roles/{id} | update role
[**UpdateServiceGroupById**](DefaultApi.md#UpdateServiceGroupById) | **Put** /service/groups/{id} | update service group


# **AckAlertById**
> RestNullObjectResponse AckAlertById($body, $id)

ack alert by id




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestAlertAck**](RestAlertAck.md)|  | 
 **id** | **string**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AckCollectorDownAlertById**
> RestNullObjectResponse AckCollectorDownAlertById($id, $body)

ack collector down alert




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestAckCollectorDown**](RestAckCollectorDown.md)|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddAdmin**
> RestAdminResponse AddAdmin($body)

add admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestAdmin**](RestAdmin.md)|  | 

### Return type

[**RestAdminResponse**](RestAdminResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddAlertNoteById**
> RestNullObjectResponse AddAlertNoteById($body, $id)

add alert note




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestAlertAck**](RestAlertAck.md)|  | 
 **id** | **string**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddAlertRule**
> RestAlertRuleResponse AddAlertRule($body)

add alert rule




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestAlertRule**](RestAlertRule.md)|  | 

### Return type

[**RestAlertRuleResponse**](RestAlertRuleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddApiTokenByAdminId**
> RestApiTokenResponse AddApiTokenByAdminId($adminId, $body)

add apiToken by admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminId** | **int32**|  | 
 **body** | [**RestApiToken**](RestApiToken.md)|  | 

### Return type

[**RestApiTokenResponse**](RestApiTokenResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddCollector**
> RestCollectorResponse AddCollector($body)

add collector




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestCollector**](RestCollector.md)|  | 

### Return type

[**RestCollectorResponse**](RestCollectorResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddCollectorGroup**
> RestCollectorGroupResponse AddCollectorGroup($body)

add collector group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestCollectorGroup**](RestCollectorGroup.md)|  | 

### Return type

[**RestCollectorGroupResponse**](RestCollectorGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDashboard**
> RestDashboardV1Response AddDashboard($body)

add dashboard




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDashboardV1**](RestDashboardV1.md)|  | 

### Return type

[**RestDashboardV1Response**](RestDashboardV1Response.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDashboardGroup**
> RestDashboardGroupResponse AddDashboardGroup($body)

add dashboard group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDashboardGroup**](RestDashboardGroup.md)|  | 

### Return type

[**RestDashboardGroupResponse**](RestDashboardGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDevice**
> RestDeviceResponse AddDevice($body, $addFromWizard)

add a new device




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDevice**](RestDevice.md)|  | 
 **addFromWizard** | **bool**|  | [optional] [default to false]

### Return type

[**RestDeviceResponse**](RestDeviceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDeviceDatasourceInstance**
> RestDeviceDataSourceInstanceResponse AddDeviceDatasourceInstance($deviceId, $hdsId, $body)

add device instance 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **body** | [**RestDeviceDataSourceInstance**](RestDeviceDataSourceInstance.md)|  | 

### Return type

[**RestDeviceDataSourceInstanceResponse**](RestDeviceDataSourceInstanceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDeviceDatasourceInstanceGroup**
> RestDeviceDataSourceInstanceGroupResponse AddDeviceDatasourceInstanceGroup($deviceId, $deviceDsId, $body)

add device datasource instance group 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **deviceDsId** | **int32**|  | 
 **body** | [**RestDeviceDataSourceInstanceGroup**](RestDeviceDataSourceInstanceGroup.md)|  | 

### Return type

[**RestDeviceDataSourceInstanceGroupResponse**](RestDeviceDataSourceInstanceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDeviceGroup**
> RestDeviceGroupResponse AddDeviceGroup($body)

add device group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDeviceGroup**](RestDeviceGroup.md)|  | 

### Return type

[**RestDeviceGroupResponse**](RestDeviceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDeviceGroupProperty**
> RestPropertyResponse AddDeviceGroupProperty($gid, $body)

add device group property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gid** | **int32**|  | 
 **body** | [**RestProperty**](RestProperty.md)|  | 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddDeviceProperty**
> RestPropertyResponse AddDeviceProperty($deviceId, $body)

add device property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **body** | [**RestProperty**](RestProperty.md)|  | 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddRole**
> RestRoleResponse AddRole($body)

add role




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestRole**](RestRole.md)|  | 

### Return type

[**RestRoleResponse**](RestRoleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AddServiceGroup**
> RestServiceGroupResponse AddServiceGroup($body)

add service group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestServiceGroup**](RestServiceGroup.md)|  | 

### Return type

[**RestServiceGroupResponse**](RestServiceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteAdminById**
> RestNullObjectResponse DeleteAdminById($id)

delete admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteAlertRuleById**
> RestAlertRuleResponse DeleteAlertRuleById($id)

delete alert rule




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestAlertRuleResponse**](RestAlertRuleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteApiTokenById**
> RestApiTokenResponse DeleteApiTokenById($adminId, $apitokenId)

delete apiToken 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminId** | **int32**|  | 
 **apitokenId** | **int32**|  | 

### Return type

[**RestApiTokenResponse**](RestApiTokenResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteCollectorById**
> RestCollectorResponse DeleteCollectorById($id)

delete collector




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestCollectorResponse**](RestCollectorResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteCollectorGroupById**
> RestCollectorGroupResponse DeleteCollectorGroupById($id)

delete collector group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestCollectorGroupResponse**](RestCollectorGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDashboardById**
> RestNullObjectResponse DeleteDashboardById($id)

delete dashboard




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDashboardGroupById**
> RestNullObjectResponse DeleteDashboardGroupById($id)

delete dashboard group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDevice**
> RestNullObjectResponse DeleteDevice($id)

delete a device




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDeviceGroupById**
> RestNullObjectResponse DeleteDeviceGroupById($id, $deleteChildren)

delete device group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **deleteChildren** | **bool**|  | [optional] [default to false]

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDeviceGroupPropertyByName**
> RestNullObjectResponse DeleteDeviceGroupPropertyByName($gid, $name)

delete device group property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gid** | **int32**|  | 
 **name** | **string**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDevicePropertyByName**
> RestNullObjectResponse DeleteDevicePropertyByName($deviceId, $name)

delete device  property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **name** | **string**|  | 

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteRoleById**
> RestRoleResponse DeleteRoleById($id)

delete role




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestRoleResponse**](RestRoleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteServiceGroupById**
> RestNullObjectResponse DeleteServiceGroupById($id, $deleteChildren)

delete service group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **deleteChildren** | **int32**|  | [optional] [default to 0]

### Return type

[**RestNullObjectResponse**](RestNullObjectResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAdminById**
> RestAdminResponse GetAdminById($id, $fields)

get admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestAdminResponse**](RestAdminResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAdminList**
> RestAdminPaginationResponse GetAdminList($fields, $size, $offset, $filter)

get admin list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestAdminPaginationResponse**](RestAdminPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertById**
> RestAlertResponse GetAlertById($id, $needMessage, $customColumns, $fields)

get alert




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 
 **needMessage** | **bool**|  | [optional] [default to false]
 **customColumns** | **string**|  | [optional] 
 **fields** | **string**|  | [optional] 

### Return type

[**RestAlertResponse**](RestAlertResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertList**
> RestAlertPaginationResponse GetAlertList($needMessage, $customColumns, $fields, $size, $offset, $searchId, $filter)

get alert list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **needMessage** | **bool**|  | [optional] [default to false]
 **customColumns** | **string**|  | [optional] 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **searchId** | **string**|  | [optional] 
 **filter** | **string**|  | [optional] 

### Return type

[**RestAlertPaginationResponse**](RestAlertPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertListByDeviceGroupId**
> RestAlertPaginationResponse GetAlertListByDeviceGroupId($id, $needMessage, $customColumns, $fields, $size, $offset, $searchId, $filter)

get device group alerts




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **needMessage** | **bool**|  | [optional] [default to false]
 **customColumns** | **string**|  | [optional] 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **searchId** | **string**|  | [optional] 
 **filter** | **string**|  | [optional] 

### Return type

[**RestAlertPaginationResponse**](RestAlertPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertListByDeviceId**
> RestAlertPaginationResponse GetAlertListByDeviceId($id, $needMessage, $customColumns, $fields, $size, $offset, $searchId, $filter)

get alerts




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **needMessage** | **bool**|  | [optional] [default to false]
 **customColumns** | **string**|  | [optional] 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **searchId** | **string**|  | [optional] 
 **filter** | **string**|  | [optional] 

### Return type

[**RestAlertPaginationResponse**](RestAlertPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertRuleById**
> RestAlertRuleResponse GetAlertRuleById($id, $fields)

get alert rule by id




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestAlertRuleResponse**](RestAlertRuleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertRuleList**
> RestAlertRulePaginationResponse GetAlertRuleList($fields, $size, $offset, $filter)

get alert rule list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestAlertRulePaginationResponse**](RestAlertRulePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetApiTokenListByAdminId**
> RestApiTokenPaginationResponse GetApiTokenListByAdminId($adminId, $fields, $size, $offset, $filter)

get apiToken by admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminId** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestApiTokenPaginationResponse**](RestApiTokenPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCollectorById**
> RestCollectorResponse GetCollectorById($id, $fields)

get collector




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestCollectorResponse**](RestCollectorResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCollectorGroupById**
> RestCollectorGroupResponse GetCollectorGroupById($id, $fields)

get collector group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestCollectorGroupResponse**](RestCollectorGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCollectorGroupList**
> RestCollectorGroupPaginationResponse GetCollectorGroupList($fields, $size, $offset, $filter)

get collector group list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestCollectorGroupPaginationResponse**](RestCollectorGroupPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCollectorList**
> RestCollectorPaginationResponse GetCollectorList($fields, $size, $offset, $filter)

get collector list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestCollectorPaginationResponse**](RestCollectorPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDashboardById**
> RestDashboardV1Response GetDashboardById($id, $fields)

get dashboard




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDashboardV1Response**](RestDashboardV1Response.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDashboardGroupById**
> RestDashboardGroupResponse GetDashboardGroupById($id, $fields)

get dashboard group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDashboardGroupResponse**](RestDashboardGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDashboardGroupList**
> RestDashboardGroupPaginationResponse GetDashboardGroupList($fields, $size, $offset, $filter)

get dashboard group list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDashboardGroupPaginationResponse**](RestDashboardGroupPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDashboardList**
> RestDashboardV1PaginationResponse GetDashboardList($fields, $size, $offset, $filter)

get dashboard list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDashboardV1PaginationResponse**](RestDashboardV1PaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceById**
> RestDeviceResponse GetDeviceById($id, $fields)

get device by id




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceResponse**](RestDeviceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceById**
> RestDeviceDatasourceResponse GetDeviceDatasourceById($deviceId, $id, $fields)

get device datasource 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceDatasourceResponse**](RestDeviceDatasourceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceDataById**
> RestDeviceDatasourceDataResponse GetDeviceDatasourceDataById($deviceId, $id, $period, $start, $end, $datapoints, $format)

get device datasource data 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **id** | **int32**|  | 
 **period** | **float64**|  | [optional] [default to 1]
 **start** | **int64**|  | [optional] [default to 0]
 **end** | **int64**|  | [optional] [default to 0]
 **datapoints** | **string**|  | [optional] 
 **format** | **string**|  | [optional] [default to json]

### Return type

[**RestDeviceDatasourceDataResponse**](RestDeviceDatasourceDataResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceAlertSettingById**
> RestDeviceDataSourceInstanceAlertSettingResponse GetDeviceDatasourceInstanceAlertSettingById($deviceId, $hdsId, $instanceId, $id, $fields)

get device instance alert setting




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **instanceId** | **int32**|  | 
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceDataSourceInstanceAlertSettingResponse**](RestDeviceDataSourceInstanceAlertSettingResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceById**
> RestDeviceDataSourceInstanceResponse GetDeviceDatasourceInstanceById($deviceId, $hdsId, $id, $fields)

get device instance 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceDataSourceInstanceResponse**](RestDeviceDataSourceInstanceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceData**
> RestDeviceDataSourceInstanceDataResponse GetDeviceDatasourceInstanceData($deviceId, $hdsId, $id, $period, $start, $end, $datapoints, $format)

get device instance data




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **id** | **int32**|  | 
 **period** | **float64**|  | [optional] [default to 1]
 **start** | **int64**|  | [optional] [default to 0]
 **end** | **int64**|  | [optional] [default to 0]
 **datapoints** | **string**|  | [optional] 
 **format** | **string**|  | [optional] [default to json]

### Return type

[**RestDeviceDataSourceInstanceDataResponse**](RestDeviceDataSourceInstanceDataResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceGraphData**
> RestGraphPlotResponse GetDeviceDatasourceInstanceGraphData($deviceId, $hdsId, $id, $graphId, $start, $end, $format)

get device instance graph data 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **id** | **int32**|  | 
 **graphId** | **int32**|  | 
 **start** | **int64**|  | [optional] 
 **end** | **int64**|  | [optional] 
 **format** | **string**|  | [optional] 

### Return type

[**RestGraphPlotResponse**](RestGraphPlotResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceGroupById**
> RestDeviceDataSourceInstanceGroupResponse GetDeviceDatasourceInstanceGroupById($deviceId, $deviceDsId, $id, $fields)

get device datasource instance group 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **deviceDsId** | **int32**|  | 
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceDataSourceInstanceGroupResponse**](RestDeviceDataSourceInstanceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceInstanceGroupOverviewGraphData**
> RestGraphPlotResponse GetDeviceDatasourceInstanceGroupOverviewGraphData($deviceId, $deviceDsId, $dsigId, $ographId, $start, $end, $format)

get device instance group overview graph data 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **deviceDsId** | **int32**|  | 
 **dsigId** | **int32**|  | 
 **ographId** | **int32**|  | 
 **start** | **int64**|  | [optional] 
 **end** | **int64**|  | [optional] 
 **format** | **string**|  | [optional] 

### Return type

[**RestGraphPlotResponse**](RestGraphPlotResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceDatasourceList**
> RestDeviceDatasourcePaginationResponse GetDeviceDatasourceList($deviceId, $fields, $size, $offset, $filter)

get device datasource list 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDeviceDatasourcePaginationResponse**](RestDeviceDatasourcePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupById**
> RestDeviceGroupResponse GetDeviceGroupById($id, $fields)

get device group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupResponse**](RestDeviceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupDatasourceAlertSetting**
> RestDeviceGroupDatasourceAlertConfigResponse GetDeviceGroupDatasourceAlertSetting($deviceGroupId, $dsId, $fields)

get device group datasource alert setting 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceGroupId** | **int32**|  | 
 **dsId** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupDatasourceAlertConfigResponse**](RestDeviceGroupDatasourceAlertConfigResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupDatasourceById**
> RestDeviceGroupDatasourceResponse GetDeviceGroupDatasourceById($deviceGroupId, $id, $fields)

get device group datasource




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceGroupId** | **int32**|  | 
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupDatasourceResponse**](RestDeviceGroupDatasourceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupDatasourceList**
> RestDeviceGroupDatasourceResponse GetDeviceGroupDatasourceList($deviceGroupId, $fields, $size, $offset, $filter)

get device group datasource list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceGroupId** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupDatasourceResponse**](RestDeviceGroupDatasourceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupList**
> RestDeviceGroupPaginationResponse GetDeviceGroupList($fields, $size, $offset, $filter)

get device group list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupPaginationResponse**](RestDeviceGroupPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupProperties**
> RestPropertyPaginationResponse GetDeviceGroupProperties($gid, $fields, $size, $offset, $filter)

get device group properties




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gid** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestPropertyPaginationResponse**](RestPropertyPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceGroupPropertyByName**
> RestPropertyResponse GetDeviceGroupPropertyByName($gid, $name, $fields)

get device group property by name




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gid** | **int32**|  | 
 **name** | **string**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceInstanceGraphDataOnlyByInstanceId**
> RestGraphPlotResponse GetDeviceInstanceGraphDataOnlyByInstanceId($instanceId, $graphId, $start, $end, $format)

get device instance graphData




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceId** | **int32**|  | 
 **graphId** | **int32**|  | 
 **start** | **int64**|  | [optional] 
 **end** | **int64**|  | [optional] 
 **format** | **string**|  | [optional] 

### Return type

[**RestGraphPlotResponse**](RestGraphPlotResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceList**
> RestDevicePaginationResponse GetDeviceList($fields, $size, $offset, $filter)

get device list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDevicePaginationResponse**](RestDevicePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDevicePropertiesList**
> RestPropertyPaginationResponse GetDevicePropertiesList($deviceId, $fields, $size, $offset, $filter)

get device properties




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestPropertyPaginationResponse**](RestPropertyPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDevicePropertyByName**
> RestPropertyResponse GetDevicePropertyByName($deviceId, $name, $fields)

get device property by name




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **name** | **string**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetImmediateDeviceListByDeviceGroupId**
> RestDevicePaginationResponse GetImmediateDeviceListByDeviceGroupId($id, $fields, $size, $offset, $filter)

get immediate devices under group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestDevicePaginationResponse**](RestDevicePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetRoleById**
> RestRoleResponse GetRoleById($id, $fields)

get role by id




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestRoleResponse**](RestRoleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetRoleList**
> RestRolePaginationResponse GetRoleList($fields, $size, $offset, $filter)

get role list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestRolePaginationResponse**](RestRolePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetServiceGraphData**
> RestGraphPlotResponse GetServiceGraphData($serviceId, $checkpointId, $graphName, $start, $end, $format)

get service graph data




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serviceId** | **int32**|  | 
 **checkpointId** | **int32**|  | 
 **graphName** | **string**|  | 
 **start** | **int64**|  | [optional] 
 **end** | **int64**|  | [optional] 
 **format** | **string**|  | [optional] 

### Return type

[**RestGraphPlotResponse**](RestGraphPlotResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetServiceGroupById**
> RestServiceGroupResponse GetServiceGroupById($id, $fields)

get service group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **fields** | **string**|  | [optional] 

### Return type

[**RestServiceGroupResponse**](RestServiceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetServiceGroupList**
> RestServiceGroupPaginationResponse GetServiceGroupList($fields, $size, $offset, $filter)

get service group list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestServiceGroupPaginationResponse**](RestServiceGroupPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSiteMonitorCheckPointList**
> RestSmCheckPointPaginationResponse GetSiteMonitorCheckPointList($fields, $size, $offset, $filter)

get site monitor checkpoint list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestSmCheckPointPaginationResponse**](RestSMCheckPointPaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUnmonitoredDeviceList**
> RestUnmonitoredDevicePaginationResponse GetUnmonitoredDeviceList($fields, $size, $offset, $filter)

get unmonitored device list




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fields** | **string**|  | [optional] 
 **size** | **int32**|  | [optional] [default to 50]
 **offset** | **int32**|  | [optional] [default to 0]
 **filter** | **string**|  | [optional] 

### Return type

[**RestUnmonitoredDevicePaginationResponse**](RestUnmonitoredDevicePaginationResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **InstallCollector**
> *os.File InstallCollector($collectorId, $osAndArch, $collectorVersion, $token, $monitorOthers, $collectorSize, $useEA)

install collector




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **collectorId** | **int32**|  | 
 **osAndArch** | **string**|  | 
 **collectorVersion** | **int32**|  | [optional] 
 **token** | **string**|  | [optional] 
 **monitorOthers** | **bool**|  | [optional] [default to true]
 **collectorSize** | **string**|  | [optional] [default to medium]
 **useEA** | **bool**|  | [optional] [default to false]

### Return type

[***os.File**](*os.File.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PatchDeviceById**
> RestDeviceResponse PatchDeviceById($body, $id, $opType, $patchFields)

patch a device




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDevice**](RestDevice.md)|  | 
 **id** | **int32**|  | 
 **opType** | **string**|  | [optional] [default to refresh]
 **patchFields** | **string**|  | [optional] 

### Return type

[**RestDeviceResponse**](RestDeviceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PatchDeviceGroupById**
> RestDeviceGroupResponse PatchDeviceGroupById($id, $body, $opType, $patchFields)

patch device group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestDeviceGroup**](RestDeviceGroup.md)|  | 
 **opType** | **string**|  | [optional] [default to refresh]
 **patchFields** | **string**|  | [optional] 

### Return type

[**RestDeviceGroupResponse**](RestDeviceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PatchServiceGroupById**
> RestServiceGroupResponse PatchServiceGroupById($id, $body, $opType, $patchFields)

patch service group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestServiceGroup**](RestServiceGroup.md)|  | 
 **opType** | **string**|  | [optional] [default to refresh]
 **patchFields** | **string**|  | [optional] 

### Return type

[**RestServiceGroupResponse**](RestServiceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ScheduleAutoDiscoveryByDeviceId**
> RestStringResponse ScheduleAutoDiscoveryByDeviceId($id)

schedule auto discovery for a host




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 

### Return type

[**RestStringResponse**](RestStringResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateAdminById**
> RestAdminResponse UpdateAdminById($id, $body, $changePassword)

update admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestAdmin**](RestAdmin.md)|  | 
 **changePassword** | **bool**|  | [optional] [default to false]

### Return type

[**RestAdminResponse**](RestAdminResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateAlertRuleById**
> RestAlertRuleResponse UpdateAlertRuleById($body, $id)

update alert rule




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestAlertRule**](RestAlertRule.md)|  | 
 **id** | **int32**|  | 

### Return type

[**RestAlertRuleResponse**](RestAlertRuleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateApiTokenByAdminId**
> RestApiTokenResponse UpdateApiTokenByAdminId($adminId, $apitokenId, $body)

update apiToken by admin




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminId** | **int32**|  | 
 **apitokenId** | **int32**|  | 
 **body** | [**RestApiToken**](RestApiToken.md)|  | 

### Return type

[**RestApiTokenResponse**](RestApiTokenResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateCollectorById**
> RestCollectorResponse UpdateCollectorById($id, $body)

update collector




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestCollector**](RestCollector.md)|  | 

### Return type

[**RestCollectorResponse**](RestCollectorResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateCollectorGroupById**
> RestCollectorGroupResponse UpdateCollectorGroupById($id, $body)

update collector group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestCollectorGroup**](RestCollectorGroup.md)|  | 

### Return type

[**RestCollectorGroupResponse**](RestCollectorGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDashboardById**
> RestDashboardV1Response UpdateDashboardById($id, $body, $overwriteGroupFields)

update dashboard




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestDashboardV1**](RestDashboardV1.md)|  | 
 **overwriteGroupFields** | **bool**|  | [optional] [default to false]

### Return type

[**RestDashboardV1Response**](RestDashboardV1Response.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDashboardGroupById**
> RestDashboardGroupResponse UpdateDashboardGroupById($id, $body)

update dashboard group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestDashboardGroup**](RestDashboardGroup.md)|  | 

### Return type

[**RestDashboardGroupResponse**](RestDashboardGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDevice**
> RestDeviceResponse UpdateDevice($body, $id, $opType)

update a device




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RestDevice**](RestDevice.md)|  | 
 **id** | **int32**|  | 
 **opType** | **string**|  | [optional] [default to refresh]

### Return type

[**RestDeviceResponse**](RestDeviceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceDatasourceInstanceAlertSettingById**
> RestDeviceDataSourceInstanceAlertSettingResponse UpdateDeviceDatasourceInstanceAlertSettingById($deviceId, $hdsId, $instanceId, $id, $body)

update device instance alert setting




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **instanceId** | **int32**|  | 
 **id** | **int32**|  | 
 **body** | [**RestDeviceDataSourceInstanceAlertSetting**](RestDeviceDataSourceInstanceAlertSetting.md)|  | 

### Return type

[**RestDeviceDataSourceInstanceAlertSettingResponse**](RestDeviceDataSourceInstanceAlertSettingResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceDatasourceInstanceById**
> RestDeviceDataSourceInstanceResponse UpdateDeviceDatasourceInstanceById($deviceId, $hdsId, $id, $body)

update device instance 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **hdsId** | **int32**|  | 
 **id** | **int32**|  | 
 **body** | [**RestDeviceDataSourceInstance**](RestDeviceDataSourceInstance.md)|  | 

### Return type

[**RestDeviceDataSourceInstanceResponse**](RestDeviceDataSourceInstanceResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceDatasourceInstanceGroupById**
> RestDeviceDataSourceInstanceGroupResponse UpdateDeviceDatasourceInstanceGroupById($deviceId, $deviceDsId, $id, $body)

update device datasource instance group 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **deviceDsId** | **int32**|  | 
 **id** | **int32**|  | 
 **body** | [**RestDeviceDataSourceInstanceGroup**](RestDeviceDataSourceInstanceGroup.md)|  | 

### Return type

[**RestDeviceDataSourceInstanceGroupResponse**](RestDeviceDataSourceInstanceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceGroupById**
> RestDeviceGroupResponse UpdateDeviceGroupById($id, $body)

update device group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestDeviceGroup**](RestDeviceGroup.md)|  | 

### Return type

[**RestDeviceGroupResponse**](RestDeviceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceGroupDatasourceAlertSetting**
> RestDeviceGroupDatasourceAlertConfigResponse UpdateDeviceGroupDatasourceAlertSetting($deviceGroupId, $dsId, $body)

update device group datasource alert setting 




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceGroupId** | **int32**|  | 
 **dsId** | **int32**|  | 
 **body** | [**RestDeviceGroupDataSourceAlertConfig**](RestDeviceGroupDataSourceAlertConfig.md)|  | 

### Return type

[**RestDeviceGroupDatasourceAlertConfigResponse**](RestDeviceGroupDatasourceAlertConfigResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDeviceGroupPropertyByName**
> RestPropertyResponse UpdateDeviceGroupPropertyByName($gid, $name, $body)

update device group property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gid** | **int32**|  | 
 **name** | **string**|  | 
 **body** | [**RestProperty**](RestProperty.md)|  | 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateDevicePropertyByName**
> RestPropertyResponse UpdateDevicePropertyByName($deviceId, $name, $body)

update device  property




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceId** | **int32**|  | 
 **name** | **string**|  | 
 **body** | [**RestProperty**](RestProperty.md)|  | 

### Return type

[**RestPropertyResponse**](RestPropertyResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateRoleById**
> RestRoleResponse UpdateRoleById($id, $body)

update role




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestRole**](RestRole.md)|  | 

### Return type

[**RestRoleResponse**](RestRoleResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateServiceGroupById**
> RestServiceGroupResponse UpdateServiceGroupById($id, $body)

update service group




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int32**|  | 
 **body** | [**RestServiceGroup**](RestServiceGroup.md)|  | 

### Return type

[**RestServiceGroupResponse**](RestServiceGroupResponse.md)

### Authorization

[LMv1](../README.md#LMv1)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

