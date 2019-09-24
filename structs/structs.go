package structs

//AtlonaAudioWrapper .
type AtlonaAudioWrapper struct {
	Audio AtlonaAudioInWrapper `json:"audio"`
}

//AtlonaAudioInWrapper .
type AtlonaAudioInWrapper struct {
	AudIn AtlonaAudio `json:"audIn"`
}

//AtlonaAudio .
type AtlonaAudio struct {
	Input1 Audio `json:"digitalIn1"`
	Input2 Audio `json:"digitalIn2"`
	Input3 Audio `json:"digitalIn3"`
	Input4 Audio `json:"digitalIn4"`
	Input5 Audio `json:"digitalIn5"`
	Input6 Audio `json:"digitalIn6"`
	Aux1   Audio `json:"analogIn1"`
	Aux2   Audio `json:"analogIn2"`
	Mic1   Audio `json:"mic1"`
}

//Audio .
type Audio struct {
	Volume int  `json:"audioVol"`
	Mute   bool `json:"audioMute"`
}

//AtlonaVideoWrapper .
type AtlonaVideoWrapper struct {
	Video AtlonaVideoOutWrapper `json:"video"`
}

//AtlonaVideoOutWrapper .
type AtlonaVideoOutWrapper struct {
	VidOut AtlonaHdmiOutwrapper `json:"vidOut"`
}

//AtlonaHdmiOutwrapper .
type AtlonaHdmiOutwrapper struct {
	HdmiOut OutputWrapper `json:"hdmiOut"`
}

//OutputWrapper .
type OutputWrapper struct {
	Output1 Output `json:"hdmiOutA"`
	Output2 Output `json:"hdmiOutB"`
}

//Output .
type Output struct {
	Name string `json:"name"`
	Src  int    `json:"videoSrc"`
}

//NetworkWrapper .
type NetworkWrapper struct {
	Network EthernetWrapper `json:"network"`
}

//EthernetWrapper .
type EthernetWrapper struct {
	Ethernet NetworkInfo `json:"eth0"`
}

//NetworkInfo .
type NetworkInfo struct {
	MacAddress    string `json:"macAddr"`
	DomainName    string `json:"domainName"`
	DNSServer1    string `json:"dnsServer1"`
	DNSServer2    string `json:"dnsServer2"`
	IPSettings    IPInfo `json:"ipSettings"`
	LastIPAddress string `json:"lastIpaddr"`
}

//IPInfo .
type IPInfo struct {
	IPAddress string `json:"ipaddr"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
}
