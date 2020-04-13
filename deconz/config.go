package deconz

// Config contains all field that contain configuration information of the
// Deconz instance.
type Config struct {
	UTC        string `json:"UTC"`
	Apiversion string `json:"apiversion"`
	Backup     struct {
		Errorcode int    `json:"errorcode"`
		Status    string `json:"status"`
	} `json:"backup"`
	Bridgeid         string `json:"bridgeid"`
	Datastoreversion string `json:"datastoreversion"`
	Devicename       string `json:"devicename"`
	Dhcp             bool   `json:"dhcp"`
	Factorynew       bool   `json:"factorynew"`
	Fwversion        string `json:"fwversion"`
	Gateway          string `json:"gateway"`
	Internetservices struct {
		Remoteaccess string `json:"remoteaccess"`
	} `json:"internetservices"`
	Ipaddress           string `json:"ipaddress"`
	Linkbutton          bool   `json:"linkbutton"`
	Localtime           string `json:"localtime"`
	Mac                 string `json:"mac"`
	Modelid             string `json:"modelid"`
	Name                string `json:"name"`
	Netmask             string `json:"netmask"`
	Networkopenduration int    `json:"networkopenduration"`
	Panid               int    `json:"panid"`
	Portalconnection    string `json:"portalconnection"`
	Portalservices      bool   `json:"portalservices"`
	Portalstate         struct {
		Communication string `json:"communication"`
		Incoming      bool   `json:"incoming"`
		Outgoing      bool   `json:"outgoing"`
		Signedon      bool   `json:"signedon"`
	} `json:"portalstate"`
	Proxyaddress     string      `json:"proxyaddress"`
	Proxyport        int         `json:"proxyport"`
	Replacesbridgeid interface{} `json:"replacesbridgeid"`
	Rfconnected      bool        `json:"rfconnected"`
	Starterkitid     string      `json:"starterkitid"`
	Swupdate         struct {
		Checkforupdate bool `json:"checkforupdate"`
		Devicetypes    struct {
			Bridge  bool          `json:"bridge"`
			Lights  []interface{} `json:"lights"`
			Sensors []interface{} `json:"sensors"`
		} `json:"devicetypes"`
		Notify      bool   `json:"notify"`
		Text        string `json:"text"`
		Updatestate int    `json:"updatestate"`
		URL         string `json:"url"`
	} `json:"swupdate"`
	Swupdate2 struct {
		Autoinstall struct {
			On         bool   `json:"on"`
			Updatetime string `json:"updatetime"`
		} `json:"autoinstall"`
		Bridge struct {
			Lastinstall string `json:"lastinstall"`
			State       string `json:"state"`
		} `json:"bridge"`
		Checkforupdate bool   `json:"checkforupdate"`
		Install        bool   `json:"install"`
		Lastchange     string `json:"lastchange"`
		Lastinstall    string `json:"lastinstall"`
		State          string `json:"state"`
	} `json:"swupdate2"`
	Swversion          string `json:"swversion"`
	Timeformat         string `json:"timeformat"`
	Timezone           string `json:"timezone"`
	UUID               string `json:"uuid"`
	Websocketnotifyall bool   `json:"websocketnotifyall"`
	Websocketport      int    `json:"websocketport"`
	Whitelist          map[string]struct {
		whitelist struct {
			CreateDate  string `json:"create date"`
			LastUseDate string `json:"last use date"`
			Name        string `json:"name"`
		}
	} `json:"whitelist"`
	Zigbeechannel int `json:"zigbeechannel"`
}
