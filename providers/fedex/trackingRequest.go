package fedex

const (
	defaultTrackingCarrier   string = ""
	defaultTrackingQualifier string = "12024~621648671050~FDEG"
	defaultAppType           string = "WTRK"
	defaultDeviceType        string = "WTRK"
)

type trackingInfo struct {
	TrackingNumberInfo trackingNumberInfo `json:"trackNumberInfo"`
}

type trackingNumberInfo struct {
	Carrier           string `json:"trackingCarrier"`
	TrackingNumber    string `json:"trackingNumber"`
	TrackingQualifier string `json:"trackingQualifier"`
}

type trackingRequest struct {
	DeviceType             string         `json:"appDeviceType"`
	AppType                string         `json:"appType"`
	SupportCurrentLocation bool           `json:"supportCurrentLocation"`
	UniqueKey              string         `json:"uniqueKey"`
	TrackingInfo           []trackingInfo `json:"trackingInfo"`
}

func newTrackingRequest(trackingNumber string) *trackingRequest {

	return &trackingRequest{
		DeviceType:             defaultDeviceType,
		AppType:                defaultAppType,
		SupportCurrentLocation: true,
		UniqueKey:              "",
		TrackingInfo: []trackingInfo{
			trackingInfo{
				TrackingNumberInfo: trackingNumberInfo{
					TrackingNumber:    trackingNumber,
					TrackingQualifier: defaultTrackingQualifier,
					Carrier:           defaultTrackingCarrier,
				},
			},
		},
	}
}
