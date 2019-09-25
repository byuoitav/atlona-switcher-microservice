package structs

//AtlonaVideo .
type AtlonaVideo struct {
	Video struct {
		VidOut struct {
			HdmiOut struct {
				HdmiOutA struct {
					VideoSrc int `json:"videoSrc"`
				} `json:"hdmiOutA"`
				HdmiOutB struct {
					VideoSrc int `json:"videoSrc"`
				} `json:"hdmiOutB"`
			} `json:"hdmiOut"`
		} `json:"vidOut"`
	} `json:"video"`
}

//AtlonaNetwork .
type AtlonaNetwork struct {
	Network struct {
		Eth0 struct {
			MacAddr    string `json:"macAddr"`
			DomainName string `json:"domainName"`
			DNSServer1 string `json:"dnsServer1"`
			DNSServer2 string `json:"dnsServer2"`
			IPSettings struct {
				TelnetPort int    `json:"telnetPort"`
				Ipaddr     string `json:"ipaddr"`
				Netmask    string `json:"netmask"`
				Gateway    string `json:"gateway"`
			} `json:"ipSettings"`
			LastIpaddr string `json:"lastIpaddr"`
			BootProto  string `json:"bootProto"`
		} `json:"eth0"`
	} `json:"network"`
}

//AtlonaAudio .
type AtlonaAudio struct {
	Audio struct {
		AudOut struct {
			ZoneOut1 struct {
				AnalogOut struct {
					AudioMute  bool `json:"audioMute"`
					AudioDelay int  `json:"audioDelay"`
				} `json:"analogOut"`
				AudioVol int `json:"audioVol"`
			} `json:"zoneOut1"`
			ZoneOut2 struct {
				AnalogOut struct {
					AudioMute  bool `json:"audioMute"`
					AudioDelay int  `json:"audioDelay"`
				} `json:"analogOut"`
				AudioVol int `json:"audioVol"`
			} `json:"zoneOut2"`
		} `json:"audOut"`
	} `json:"audio"`
}
