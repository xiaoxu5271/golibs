package golibs

import "time"

const (
	SUBJ_LOCAL_PREFIX = "iot.local."     // App and Drivers
	SUBJ_BRD_PREFIX   = "iot.broadcast." // Broadcast

)

// Standard Request and Reply

type STStdRsp struct {
	CmdId      int     `json:"cmd_id"`
	Cmd        int     `json:"cmd"`
	ResultCode EURCode `json:"result_code"`
	ErrCode    string  `json:"err_code"`
}

const (
	STD_RESULT_SUCCESS EURCode = iota
	STD_RESULT_FAILED
	STD_ERR_NULL = "NULL"
)

type EURCode int

func (s *STStdRsp) Reply(err error, cmdId int, cmd int, ntf *TopicFilt, head StdHead) {
	s.CmdId = cmdId
	s.Cmd = cmd

	if err != nil {
		s.ResultCode = STD_RESULT_FAILED
		s.ErrCode = err.Error()
	} else {
		s.ResultCode = STD_RESULT_SUCCESS
		s.ErrCode = STD_ERR_NULL
	}
	ntf.ReplyTo(head, s)
}

// Device and Software information

const (
	REQ_COMP_INF = "req.comp.info"
	RSP_COMP_INF = "rsp.comp.info"
)

type STInfReq struct {
	Cmd EUInfCmd `json:"cmd"`
}

const (
	GETINF_CMD_GET EUInfCmd = iota
)

type EUInfCmd int

type STInfRsp struct {
	Software  string `json:"software"`
	SfVersion string `json:"sfVersion"`
}

//nats req deploy
const (
	REQ_DEPLOY_FACTORY = "req.deploy.local.factory"
	RSP_DEPLOY_FACTORY = "rsp.deploy.local.factory"
)

// File Transfer

const (
	REQ_TRANS_FILE = "req.file.local.transfer"
	RSP_TRANS_FILE = "rsp.file.local.transfer"
)

type STTransReq struct {
	CmdId int           `json:"cmd_id"`
	Cmd   EUTransCmd    `json:"cmd"`
	Para  []STTransPara `json:"para"`
}

type STTransPara struct {
	Protocol  EUTransType    `json:"protocol"`
	SrcPath   string         `json:"src_path"`
	DstPath   string         `json:"dst_path"`
	Secret    EUTransSecret  `json:"secret"`
	User      string         `json:"user"`
	Passwd    string         `json:"passwd"`
	Cert      string         `json:"cert"`
	CheckMeth EUTransChkMeth `json:"check_meth"`
	CheckSum  string         `json:"check_sum"`
}

const (
	TRANSFILE_CMD_DWN EUTransCmd = iota
	TRANSFILE_CMD_UP
)

type EUTransCmd int

const (
	TRANSFILE_PROTO_HTTP EUTransType = iota
	TRANSFILE_PROTO_FTP
	TRANSFILE_PROTO_TFTP
	TRANSFILE_PROTO_MINIO
)

type EUTransType int

const (
	TRANSFILE_SECRET_NONE EUTransSecret = iota
	TRANSFILE_SECRET_USER
	TRANSFILE_SECRET_CERT
	TRANSFILE_SECRET_BOTH
)

type EUTransSecret int

const (
	TRANSFILE_CHECKMETH_NONE EUTransChkMeth = iota
	TRANSFILE_CHECKMETH_MD5
	TRANSFILE_CHECKMETH_SHA224
	TRANSFILE_CHECKMETH_SHA256
	TRANSFILE_CHECKMETH_SHA384
	TRANSFILE_CHECKMETH_SHA512
)

type EUTransChkMeth int

// Topology configuration

const (
	REQ_TOPO = "req.aVd.local"
	RSP_TOPO = "rsp.aVd.local"
)

type STTopoReq struct {
	CmdId int          `json:"cmd_id"`
	Cmd   EUTopoCmd    `json:"cmd"`
	Para  []STTopoPara `json:"para"`
}

type STTopoRsp struct {
	CmdId      int          `json:"cmd_id"`
	Cmd        EUTopoCmd    `json:"cmd"`
	ResultCode EURCode      `json:"result_code"`
	ErrCode    string       `json:"err_code"`
	Para       []STTopoPara `json:"para"`
}
type STTopoPara struct {
	State   EUTopoState   `json:"state"`
	Type    EUTopoExeType `json:"type"`
	AVd     string        `json:"aVd"`
	Driver  string        `json:"driver"`
	Version string        `json:"version"`
}

const (
	TOPO_STATE_EN EUTopoState = iota
	TOPO_STATE_DIS
)

type EUTopoState int

const (
	TOPO_CMD_ADD EUTopoCmd = iota
	TOPO_CMD_DEL
	TOPO_CMD_STOP
	TOPO_CMD_START
	TOPO_CMD_GET
	TOPO_CMD_SET
)

type EUTopoCmd int

const (
	TOPO_TYPE_APPLICATION EUTopoExeType = iota
	TOPO_TYPE_COMPUTE
	TOPO_TYPE_DRIVER
)

type EUTopoExeType int

// Scheduling Control

const (
	REQ_SCHEDULING_CTL   = "req.module.control"
	BRD_SCHEDULING_STATE = "scheduling.module.state"
	REQ_RUN_CTL          = "req.collect.local.link"
	RSP_RUN_CTL          = "rsp.collect.local.link"
)

type STSchedCtlReq struct {
	Cmd     EUScheCtlCmd `json:"cmd"`
	Id      string       `json:"id"`
	Program string       `json:"program"`
	Args    []string     `json:"args"`
}

type STSchedBrdState struct {
	Para []STSchedBrdStatePara `json:"para"`
}

type STSchedBrdStatePara struct {
	State   EUScheState `json:"state"`
	Id      string      `json:"id"`
	Program string      `json:"program"`
	Args    []string    `json:"args"`
}

// 第三方接口业务管理
type STRunCtlReq struct {
	CmdId int       `json:"cmd_id"`
	Cmd   RunCtlCmd `json:"cmd"`
}

const (
	SCHEDULING_CMD_STOP EUScheCtlCmd = iota
	SCHEDULING_CMD_START
	SCHEDULING_CMD_STOP_ALL
	SCHEDULING_CMD_START_ALL
	SCHEDULING_CMD_REBOOT
)

const (
	RUNCTL_CMD_CLOSE RunCtlCmd = RunCtlCmd(SCHEDULING_CMD_STOP)
)

type RunCtlCmd int

const (
	SCHEDULING_STS_STOP EUScheState = iota
	SCHEDULING_STS_START
)

type EUScheCtlCmd int
type EUScheState int

// Update

const (
	REQ_UPDATE = "req.comp.update"
	RSP_UPDATE = "rsp.comp.update"
)

type STUpdateReq struct {
	Cmd  EUUpdateCmd `json:"cmd"`
	Path string      `json:"path"`
	Key  string      `json:"key"`
}

const (
	UPDATE_CMD_MGT EUUpdateCmd = iota
	UPDATE_CMD_SCHEDULING
	UPDATE_CMD_WEB
	UPDATE_CMD_WEBBOX
	UPDATE_CMD_FACTORY
)

type EUUpdateCmd int

// Share business file

const (
	REQ_SHARE = "req.file.local.share"
	RSP_SHARE = "rsp.file.local.share"
)

type STShareReq struct {
	CmdId int           `json:"cmd_id"`
	Cmd   EUShareCmd    `json:"cmd"`
	Para  []STSharePara `json:"para"`
}

type STSharePara struct {
	Method EUShareMeth `json:"method"`
	Path   string      `json:"path"`
	User   string      `json:"user"`
	Key    string      `json:"key"`
}

const (
	SHARE_CMD_ENABLE EUShareCmd = iota
	SHARE_CMD_DISABLE
)

type EUShareCmd int

const (
	SHARE_USAGE_FTP EUShareMeth = iota
	SHARE_USAGE_SAMBA
)

type EUShareMeth int

// Gateway Factory Information

const (
	REQ_HW_INFO = "req.gateway.info"
	RSP_HW_INFO = "rsp.gateway.info"
)

type STHWInfo struct {
	HwSN   string `json:"hwSN"`
	HwMode string `json:"hwMode"`
	HwVer  string `json:"hwVer"`
	SfVer  string `json:"sfVer"`
	PD     string `json:"PD"`
}

// Route Management

const (
	REQ_ROUTE_MG = "req.route.local"
	RSP_ROUTE_MG = "rsp.route.local"
)

type STRouteMgReq struct {
	CmdId int             `json:"cmd_id"`
	Cmd   EURouteMgCmd    `json:"cmd"`
	Para  []STRouteMgPara `json:"para"`
}

type STRouteMgPara struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Objective string `json:"objective"`
}

const (
	ROUTEMG_CMD_ADD EURouteMgCmd = iota
	ROUTEMG_CMD_DEL
	ROUTEMG_CMD_GET
	ROUTEMG_CMD_SET
	ROUTEMG_CMD_CLS EURouteMgCmd = 100
)

type EURouteMgCmd int

type STRouteMgGetRsp struct {
	CmdId      int             `json:"cmd_id"`
	Cmd        EURouteMgCmd    `json:"cmd"`
	ResultCode EURCode         `json:"result_code"`
	ErrCode    string          `json:"err_code"`
	Para       []STRouteMgPara `json:"para"`
}

// Hardware Configuration and State

const (
	REQ_GATEWAY_SYS = "req.gateway.sys"
	RSP_GATEWAY_SYS = "rsp.gateway.sys"
)

type STHWSysReq struct {
	CmdId int         `json:"cmd_id"`
	Cmd   EUHWSysCmd  `json:"cmd"`
	Para  STHWSysPara `json:"para"`
}

type STHWSysRsp struct {
	CmdId int         `json:"cmd_id"`
	Cmd   EUHWSysCmd  `json:"cmd"`
	Para  STHWSysPara `json:"para"`
}

const (
	HWSYS_CMD_SET EUHWSysCmd = iota
	HWSYS_CMD_GET
	HWSYS_CMD_GETALL
)

type EUHWSysCmd int

type STHWSysPara struct {
	Version  uint             `json:"version"`
	WAN      *STHWSysWAN      `json:"wan"`
	Mobile   *STHWSysMobile   `json:"mobile"`
	LTE      *STHWSysLTE      `json:"LTE"`
	LAN      []STHWSysLAN     `json:"lan"`
	Route    []STHWSysRoute   `json:"route"`
	Bridge   []STHWSysBridge  `json:"bridge"`
	Position *STHWSysPosition `json:"position"`
	Cpu      *STHWSysCpu      `json:"cpu"`
	Mem      *STHWSysMem      `json:"mem"`
	Disk     *STHWSysDisk     `json:"disk"`
}

type STHWSysWAN struct {
	WAN string   `json:"wansel"`
	DNS []string `json:"dns"`
}

type STHWSysMobile struct {
	STHWSysMobileConfig `json:"config"`
	*STHWSysMobileState `json:"state"`
}

type STHWSysMobileConfig struct {
	Mode     EUHWSysMobileMode   `json:"mode"`
	Apn      string              `json:"apn"`
	User     string              `json:"user"`
	Password string              `json:"password"`
	Watchdog EUHWSysWatchdogMode `json:"watchdog"`
}

const (
	HWSYS_MOBILE_MODE_SYS  EUHWSysMobileMode = "SYS"
	HWSYS_MOBILE_MODE_USER EUHWSysMobileMode = "USER"
)

type EUHWSysMobileMode string

const (
	HWSYS_WATCHDOG_MODE_OFF EUHWSysWatchdogMode = iota
	HWSYS_WATCHDOG_MODE_ON
)

type EUHWSysWatchdogMode int

type STHWSysMobileState struct {
	Result    EUHWSysMobileResult `json:"result"`
	ErrorCode EUHWSysMobileECode  `json:"error_code"`
	ICCID     string              `json:"ICCID"`
	IMEI      string              `json:"IMEI"`
	IMSI      string              `json:"IMSI"`
	Model     string              `json:"model"`
	Soft      string              `json:"soft"`
}

const (
	HWSYS_MOBILE_RET_NOTEXIST EUHWSysMobileResult = iota
	HWSYS_MOBILE_RET_OK
	HWSYS_MOBILE_RET_FAULT
)

type EUHWSysMobileResult int

const (
	HWSYS_MOBILE_ECODE_NOSIM EUHWSysMobileECode = iota
	HWSYS_MOBILE_ECODE_ANT
	HWSYS_MOBILE_ECODE_DENY
)

type EUHWSysMobileECode int

type STHWSysLTE struct {
	RSSI int `json:"rssi"`
	SMR  int `json:"smr"`
}

type STHWSysLAN struct {
	STHWSysLANConfig `json:"config"`
	*STHWSysLANState `json:"state"`
}

type STHWSysLANConfig struct {
	Port    EUHWSysLanPort   `json:"port"`
	Speed   *EUHWSysLanSpeed `json:"speed"`
	DHCP    *EUHWSysLanDHCP  `json:"dhcp"`
	IPaddr  string           `json:"ipaddr"`
	Netmask string           `json:"netmask"`
	Gateway string           `json:"gateway"`
}

type STHWSysLANState struct {
	Carrier EUHWSysLanCarrier `json:"carrier"`
}

const (
	HWSYS_NET_PORT_LAN0   EUHWSysLanPort = "lan_0"
	HWSYS_NET_PORT_LAN1   EUHWSysLanPort = "lan_1"
	HWSYS_NET_PORT_MOBILE string         = "mobile"
)

type EUHWSysLanPort string

const (
	HWSYS_LAN_SPEED_AUTO EUHWSysLanSpeed = iota
	HWSYS_LAN_SPEED_10H
	HWSYS_LAN_SPEED_10F
	HWSYS_LAN_SPEED_100H
	HWSYS_LAN_SPEED_100F
)

type EUHWSysLanSpeed int

const (
	HWSYS_LAN_DHCP_STATIC EUHWSysLanDHCP = iota
	HWSYS_LAN_DHCP_DYNAMIC
)

type EUHWSysLanDHCP int

const (
	HWSYS_LAN_CARRIER_OFF EUHWSysLanCarrier = iota
	HWSYS_LAN_CARRIER_ON
	HWSYS_LAN_CARRIER_FAULT
)

type EUHWSysLanCarrier int

type STHWSysRoute struct {
	Port    string `json:"port"`
	Dest    string `json:"dest"`
	Gateway string `json:"gateway"`
	Mask    string `json:"mask"`
	Metric  int    `json:"metric"`
}

type STHWSysBridge struct {
	Port    string            `json:"port"`
	Mode    EUHWSysBridgeMode `json:"mode"`
	IPaddr  string            `json:"ipaddr"`
	Gateway string            `json:"gateway"`
	Netmask string            `json:"netmask"`
	Port0   string            `json:"port0"`
	Port1   string            `json:"port1"`
}

const (
	HWSYS_BRIDGE_MODE_OFF EUHWSysBridgeMode = iota
	HWSYS_BRIDGE_MODE_ON
)

type EUHWSysBridgeMode int

type STHWSysPosition struct {
	STHWSysPositionConfig `json:"config"`
	*STHWSysPositionState `json:"state"`
}

type STHWSysPositionConfig struct {
	Mode EUHWPositionMode `json:"mode"`
}

type STHWSysPositionState struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

const (
	HWSYS_POSITION_MODE_NONE EUHWPositionMode = iota
	HWSYS_POSITION_MODE_MOBILE
	HWSYS_POSITION_MODE_GNSS
)

type EUHWPositionMode int

type STHWSysCpu struct {
	Usage float64 `json:"usage"`
}

type STHWSysMem struct {
	Used  string `json:"used"`
	Total string `json:"total"`
}

type STHWSysDisk struct {
	Used  string `json:"used"`
	Total string `json:"total"`
}

// Business data pass-through

const (
	REQ_APPMES_UP   = "req.app.topic.up"
	REQ_APPMES_DOWN = "req.app.topic.down"
	RSP_APPMES_UP   = "rsp.app.topic.up"
	RSP_APPMES_DOWN = "rsp.app.topic.down"
	REQ_APPMES_SET  = "req.app.topic.set"
	RSP_APPMES_SET  = "rsp.app.topic.set"
)

type STAppSetReq struct {
	CmdId int         `json:"cmd_id"`
	Cmd   EUAppSetCmd `json:"cmd"`
	Topic []string    `json:"topic"`
}

type STAppSetRsp struct {
	CmdId      int         `json:"cmd_id"`
	Cmd        EUAppSetCmd `json:"cmd"`
	ResultCode EURCode     `json:"result_code"`
	ErrCode    string      `json:"err_code"`
	Para       []string    `json:"para"`
}

const (
	APPSET_CMD_SET EUAppSetCmd = iota
	APPSET_CMD_GET
)

type EUAppSetCmd int

type STAppmsgReq struct {
	CmdId int            `json:"cmd_id"`
	Cmd   EUAppmsgCmd    `json:"cmd"`
	Para  []STAppMsgPara `json:"para"`
}

type STAppMsgPara struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

const (
	APPMSG_CMD_SEND EUAppmsgCmd = iota
)

type EUAppmsgCmd int

//cmd req
type STCmdReq struct {
	GatewayId string                 `json:"gateway_id"`
	AppId     string                 `json:"app_id"`
	ObjDevId  string                 `json:"object_device_id"`
	CmdName   string                 `json:"command_name"`
	ServiceId string                 `json:"service_id"`
	Paras     map[string]interface{} `json:"paras"`
}

//cmd rsp
type STCmdRsp struct {
	GatewayId  string                 `json:"gateway_id"`
	AppId      string                 `json:"app_id"`
	ObjDevId   string                 `json:"object_device_id"`
	ResultCode int                    `json:"result_code"`
	ServiceId  string                 `json:"service_id"`
	RspName    string                 `json:"response_name"`
	Paras      map[string]interface{} `json:"paras"` //"result": "success"
}

//Properties report
type STPropRpt struct {
	GatewayId string          `json:"gateway_id"`
	AppId     string          `json:"app_id"`
	DeviceId  string          `json:"device_id"`
	Services  []STPorpService `json:"services"`
}

type STPorpService struct {
	ServiceId string                 `json:"service_id"`
	Data      map[string]interface{} `json:"data"`
	EventTime string                 `json:"event_time"`
}

//Properties set req
type STPropSetReq struct {
	GatewayId string                `json:"gateway_id"`
	AppId     string                `json:"app_id"`
	DeviceId  string                `json:"device_id"`
	Services  []STPorpSetReqService `json:"services"`
}

type STPorpSetReqService struct {
	ServiceId  string                 `json:"service_id"`
	Properties map[string]interface{} `json:"properties"`
}

//Properties set rsp
type STPropSetRsp struct {
	GatewayId  string `json:"gateway_id"`
	AppId      string `json:"app_id"`
	ObjDevId   string `json:"object_device_id"`
	ResultCode int    `json:"result_code"`
	ResultDesc string `json:"result_desc"` //"success"
}

//Properties get req
type STPropGetReq struct {
	GatewayId string `json:"gateway_id"`
	AppId     string `json:"app_id"`
	ObjDevId  string `json:"object_device_id"`
	ServiceId string `json:"service_id"`
}

//Properties get rsp
type STPropGetRsp struct {
	GatewayId string          `json:"gateway_id"`
	AppId     string          `json:"app_id"`
	ObjDevId  string          `json:"object_device_id"`
	Services  []STPorpService `json:"services"`
}

//Shadow  Get req
type STShadowGetReq struct {
	GatewayId string `json:"gateway_id"`
	AppId     string `json:"app_id"`
	ObjDevId  string `json:"object_device_id"`
	ServiceId string `json:"service_id"`
}

//Shadow  Get rsp
type STShadowGetRsp struct {
	GatewayId string     `json:"gateway_id"`
	AppId     string     `json:"app_id"`
	ObjDevId  string     `json:"object_device_id"`
	Shadow    []STShadow `json:"shadow"`
}

type STShadow struct {
	ServiceId string      `json:"service_id"`
	Desired   STShadowDes `json:"desired"`
	Reported  STShadowRpt `json:"reported"`
	Version   int         `json:"version"`
}

type STShadowDes struct {
	Properties map[string]interface{} `json:"properties"`
	EventTime  string                 `json:"event_time"`
}

type STShadowRpt struct {
	Properties map[string]interface{} `json:"properties"`
	EventTime  string                 `json:"event_time"`
}

//deploy lock
const (
	REQ_DEPLOY_LOCAL = "req.deploy.local"
	RSP_DEPLOY_LOCAL = "rsp.deploy.local"
)

const (
	DEPLOY_UNLOCK EUDeployState = iota
	DEPLOY_LOCK
)

type EUDeployState int

const (
	DEPLOY_CMD_LOCK_REQ EUDeployCmd = iota
	DEPLOY_CMD_UNLOCK_REQ
)

type EUDeployCmd int

type STDeployReq struct {
	CmdId    int           `json:"cmd_id"`
	Cmd      EUDeployCmd   `json:"cmd"`
	TimeOuts time.Duration `json:"timeouts"`
}

//configuration files create
const (
	REQ_USER_CFG_LOCAL = "req.configuration.local"
	RSP_USER_CFG_LOCAL = "rsp.configuration.local"
)

const (
	USER_CFG_CMD_SET EUUserCfgCmd = iota
	USER_CFG_CMD_GET
)

type EUUserCfgCmd int

const (
	USER_CFG_TYPE_APP EUUserCfgType = iota
	USER_CFG_TYPE_DRV
)

type EUUserCfgType int

type STUserCfgReq struct {
	CmdId int             `json:"cmd_id"`
	Cmd   EUUserCfgCmd    `json:"cmd"`
	Para  []STUserCfgPara `json:"para"`
}

type STUserCfgPara struct {
	Type        EUUserCfgType `json:"type"`
	Id          string        `json:"id"`
	FileName    string        `json:"file_name"`
	FileContent []byte        `json:"file_content"`
}

type STUserCfgRsp struct {
	CmdId      int             `json:"cmd_id"`
	Cmd        EUUserCfgCmd    `json:"cmd"`
	ResultCode EURCode         `json:"result_code"`
	ErrCode    string          `json:"err_code"`
	Para       []STUserCfgPara `json:"para"`
}

//driver and app manage
const (
	REQ_DRIVER_LOCAL = "req.driver.local"
	RSP_DRIVER_LOCAL = "rsp.driver.local"
)

const (
	DRV_CMD_ADD EUDrvCmd = iota
	DRV_CMD_DEL
	DRV_CMD_GET
	DRV_CMD_GET_ALL
)

const (
	DRV_TYPE_APPLICATION EUDrvTypeCmd = iota
	DRV_TYPE_COMPUTE
	DRV_TYPE_DRIVER
)

type EUDrvCmd int

type EUDrvTypeCmd int

type STDrvReq struct {
	CmdId int         `json:"cmd_id"`
	Cmd   EUDrvCmd    `json:"cmd"`
	Para  []STDrvPara `json:"para"`
}

type STDrvPara struct {
	Type   EUDrvTypeCmd `json:"type"`
	Path   string       `json:"path"`
	Driver string       `json:"driver"`
	Name   []string     `json:"name"`
}

type STDrvRsp struct {
	CmdId      int         `json:"cmd_id"`
	Cmd        EUDrvCmd    `json:"cmd"`
	ResultCode EURCode     `json:"result_code"`
	ErrCode    string      `json:"err_code"`
	Para       []STDrvPara `json:"para"`
}

/** Virtual Device Driver */

// Device Driver Collect Link Config

const (
	REQ_DRV_LINK_CFG = "req.collect.remote.set"
	RSP_DRV_LINK_CFG = "rsp.collect.remote.set"
)

type STDrvLinkCfgReq struct {
	CmdId int                `json:"cmd_id"`
	Cmd   EUDrvLinkCfgCmd    `json:"cmd"`
	Para  []STDrvLinkCfgPara `json:"para"`
}

type STDrvLinkCfgRsp struct {
	CmdId      int                `json:"cmd_id"`
	Cmd        EUDrvLinkCfgCmd    `json:"cmd"`
	ResultCode EURCode            `json:"result_code"`
	ErrCode    string             `json:"err_code"`
	Para       []STDrvLinkCfgPara `json:"para"`
}

type STDrvLinkCfgPara struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	DRV_LINK_CFG_CMD_SET EUDrvLinkCfgCmd = iota
	DRV_LINK_CFG_CMD_GET
)

type EUDrvLinkCfgCmd int

// Collection Link Control

const (
	REQ_DRV_LINK_CTL = "req.collect.remote.link"
	RSP_DRV_LINK_CTL = "rsp.collect.remote.link"
)

type STDrvLinkCtlReq struct {
	Cmd EUDrvLinkCtlCmd `json:"cmd"`
}

type STDrvLinkCtlRsp struct {
	Cmd        EUDrvLinkCtlCmd `json:"cmd"`
	ResultCode EURCode         `json:"result_code"`
	ErrCode    string          `json:"err_code"`
}

const (
	DRV_LINK_CTL_CMD_CNN EUDrvLinkCtlCmd = iota
	DRV_LINK_CTL_CMD_DISCNN
	DRV_LINK_CTL_CMD_RECNN
)

type EUDrvLinkCtlCmd int

// Collect Data Config

const (
	REQ_COLLECT_DATA_CFG = "req.collect.remote.config"
	RSP_COLLECT_DATA_CFG = "rsp.collect.remote.config"
)

type STCollectDataCfgReq struct {
	Cmd  EUCollectDataCfgCmd    `json:"cmd"`
	Para []STCollectDataCfgPara `json:"para"`
}

type STCollectDataCfgRsp = STCollectDataCfgReq

type STCollectDataCfgPara struct {
	Func    string           `json:"func"`
	Active  EUCollectDataAct `json:"active"`
	Refresh int              `json:"refresh"`
}

const (
	COLLECTION_CFG_CMD_SET EUCollectDataCfgCmd = iota
	COLLECTION_CFG_CMD_GET
	COLLECTION_CFG_CMD_GETALL
)

type EUCollectDataCfgCmd int

const (
	COLLECTION_CFG_DISABLE EUCollectDataAct = iota
	COLLECTION_CFG_ENABLE
)

type EUCollectDataAct int

// Collect User Data Config

const (
	REQ_COLLECT_UDATA_CFG = "req.collect.remote.userconfig"
	RSP_COLLECT_UDATA_CFG = "rsp.collect.remote.userconfig"
)

type STCollectionUsrCfgReq struct {
	CmdId int                  `json:"cmd_id"`
	Cmd   EUCollectUDataCfgCmd `json:"cmd"`
	Path  string               `json:"box_path"`
}

const (
	COLLECT_UDATA_CFG_CMD_SET EUCollectUDataCfgCmd = iota
	COLLECT_UDATA_CFG_CMD_GET
)

type EUCollectUDataCfgCmd int

// Device File Operation

const (
	REQ_DEV_FILE_OPERATION = "req.file.remote.operate"
	RSP_DEV_FILE_OPERATION = "rsp.file.remote.operate"
)

type STDevFileOpReq struct {
	CmdId int               `json:"cmd_id"`
	Cmd   EUDevFileOpCmd    `json:"cmd"`
	Para  []STDevFileOpPara `json:"para"`
}

type STDevFileOpRsp = STStdRsp

type STDevFileOpPara struct {
	Src string `json:"old_path"`
	Dst string `json:"new_path"`
}

const (
	DEV_FILE_OP_CMD_NEW EUDevFileOpCmd = iota
	DEV_FILE_OP_CMD_DEL
	DEV_FILE_OP_CMD_RENAME
	DEV_FILE_OP_CMD_MOVE
	DEV_FILE_OP_CMD_COPY
)

type EUDevFileOpCmd int

// Device File Transfer

const (
	REQ_DEV_FILE_TRANS = "req.file.remote.transfer"
	RSP_DEV_FILE_TRANS = "rsp.file.remote.transfer"
)

type STDevFileTransReq struct {
	CmdId int                  `json:"cmd_id"`
	Cmd   EUDevFileTransCmd    `json:"cmd"`
	Para  []STDevFileTransPara `json:"para"`
}

type STDevFileTransRsp = STStdRsp

type STDevFileTransPara struct {
	GwPath  string `json:"box_path"`
	DevPath string `json:"dev_path"`
}

const (
	DEV_FILE_TRANS_CMD_2DEV EUDevFileTransCmd = 100 + iota
	DEV_FILE_TRANS_CMD_2GW
)

type EUDevFileTransCmd int

// Device File List

const (
	REQ_DEV_FILE_LS = "req.file.remote.list"
	RSP_DEV_FILE_LS = "rsp.file.remote.list"
)

type STDevFileLsReq struct {
	CmdId int                `json:"cmd_id"`
	Cmd   EUDevFileLsCmd     `json:"cmd"`
	Para  STDevFileLsReqPara `json:"para"`
}

type STDevFileLsReqPara struct {
	DevPath string `json:"dev_path"`
}

type STDevFileLsRsp struct {
	CmdId      int                `json:"cmd_id"`
	Cmd        EUDevFileLsCmd     `json:"cmd"`
	ResultCode EURCode            `json:"result_code"`
	ErrCode    string             `json:"err_code"`
	Para       STDevFileLsRspPara `json:"para"`
}

type STDevFileLsRspPara struct {
	Dir  string                   `json:"dir"`
	List []STDevFileLsRspParaInfo `json:"list"`
}

type STDevFileLsRspParaInfo struct {
	Name string          `json:"name"`
	Type EUDevFileLsType `json:"type"`
	Size int             `json:"size"`
	Time int             `json:"time"`
}

const (
	DEV_FILE_LS_CMD EUDevFileLsCmd = 10 + iota
)

type EUDevFileLsCmd int

const (
	DEV_FILE_LS_TYPE_DIR EUDevFileLsType = iota
	DEV_FILE_LS_TYPE_FILE
)

type EUDevFileLsType int

// Collect data broadcast

const (
	BRD_COLLECT_DATA = "collect.remote.report"
)

type STCollectDataBrd struct {
	Name   string              `json:"name"`
	Number int                 `json:"number"`
	Para   []STCollectDataPara `json:"para"`
}

type STCollectDataPara struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 虚拟业务（更新/查询）
const (
	REQ_AVD_CONFIG = "req.aVd.config"
	RSP_AVD_CONFIG = "rsp.aVd.config"
)

type STAvdCfgReq struct {
	CmdId int            `json:"cmd_id"`
	Cmd   EUAvdCfgCmd    `json:"cmd"`
	Para  []STAvdCfgPara `json:"para"`
}

type STAvdCfgRsp struct {
	CmdId      int            `json:"cmd_id"`
	Cmd        EUAvdCfgCmd    `json:"cmd"`
	ResultCode EURCode        `json:"result_code"`
	ErrCode    string         `json:"err_code"`
	Para       []STAvdCfgPara `json:"para"`
}
type STAvdCfgPara struct {
	Type    EUAvdCfgExeType `json:"type"`
	AVd     string          `json:"aVd"`
	Driver  string          `json:"driver"`
	Version string          `json:"version"`
}

const (
	AVDCONF_CMD_UP EUAvdCfgCmd = iota
	AVDCONF_CMD_GET
)

type EUAvdCfgCmd int

const (
	AVDCONF_TYPE_AVD EUAvdCfgExeType = iota
	AVDCONF_TYPE_APP_DRV
)

type EUAvdCfgExeType int

// 硬件通知到
type STGatwayReq struct {
	Para []STGatwayPara `json:"para"`
}

type STGatwayPara struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//factory.base.get

const (
	REQ_FACTORY_BASE_GET = "req.factory.base.get"
	RSP_FACTORY_BASE_GET = "rsp.factory.base.get"
)

type STFacBaseReq struct {
	CmdId int          `json:"cmd_id"`
	Cmd   EUFacBaseCmd `json:"cmd"`
}

type STFacBaseRsp struct {
	CmdId     int          `json:"cmd_id"`
	Cmd       EUFacBaseCmd `json:"cmd"`
	ProductId string       `json:"productid"`
	DeviceId  string       `json:"deviceid"`
}

const (
	FAC_BASE_GET EUFacBaseCmd = iota
)

type EUFacBaseCmd int

// Collect data write

const (
	REQ_COLLECT_WRITE = "req.collect.remote.write"
	RSP_COLLECT_WRITE = "rsp.collect.remote.write"
)
