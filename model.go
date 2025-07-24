package main

type SDKError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type Bundle struct {
	BundleDisplayLabel      string `json:"BundleDisplayLabel"`
	BundleID                string `json:"BundleId"`
	BundleSalesState        string `json:"BundleSalesState"`
	BundleType              string `json:"BundleType"`
	BundleTypeDescription   string `json:"BundleTypeDescription"`
	CPU                     int    `json:"CPU"`
	InternetChargeType      string `json:"InternetChargeType"`
	InternetMaxBandwidthOut int    `json:"InternetMaxBandwidthOut"`
	Memory                  int    `json:"Memory"`
	MonthlyTraffic          int    `json:"MonthlyTraffic"`
	Price                   struct {
		InstancePrice struct {
			Currency            string `json:"Currency"`
			Discount            int    `json:"Discount"`
			DiscountPrice       int    `json:"DiscountPrice"`
			OriginalBundlePrice int    `json:"OriginalBundlePrice"`
			OriginalPrice       int    `json:"OriginalPrice"`
		} `json:"InstancePrice"`
	} `json:"Price"`
	SupportLinuxUnixPlatform bool   `json:"SupportLinuxUnixPlatform"`
	SupportWindowsPlatform   bool   `json:"SupportWindowsPlatform"`
	SystemDiskSize           int    `json:"SystemDiskSize"`
	SystemDiskType           string `json:"SystemDiskType"`
	TrafficUnlimited         bool   `json:"TrafficUnlimited"`
}
type DescribeBundlesResp struct {
	Response struct {
		Error      *SDKError `json:"Error,omitempty"`
		BundleSet  []Bundle  `json:"BundleSet"`
		RequestID  string    `json:"RequestId"`
		TotalCount int       `json:"TotalCount"`
	} `json:"Response"`
}

type CreateInstanceResp struct {
	Response struct {
		Error         *SDKError `json:"Error,omitempty"`
		InstanceIDSet []string  `json:"InstanceIdSet,omitempty"`
		RequestID     string    `json:"RequestId"`
	} `json:"Response"`
}
