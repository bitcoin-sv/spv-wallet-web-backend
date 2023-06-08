package buxclient

// AccessKey is a struct that contains access key data.
type AccessKey struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

// XPub is a struct that contains xpub data.
type XPub struct {
	Id             string `json:"id"`
	XPub           string `json:"xpub"`
	CurrentBalance uint64 `json:"current_balance"`
}

// GetAccessKey returns access key.
func (a *AccessKey) GetAccessKey() string {
	return a.Key
}

// GetAccessKeyId returns access key id.
func (a *AccessKey) GetAccessKeyId() string {
	return a.Id
}

// GetId returns xpub id.
func (x *XPub) GetId() string {
	return x.Id
}

// GetXPub returns xpub.
func (x *XPub) GetXPub() string {
	return x.XPub
}

// GetCurrentBalance returns current balance.
func (x *XPub) GetCurrentBalance() uint64 {
	return x.CurrentBalance
}
