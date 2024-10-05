package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/infastin/gorack/validation"
)

type GeneralInfo struct {
	UserID     [16]byte `json:"user_id"`
	Device     Device   `json:"device"`
	AppVersion string   `json:"app_version"`
}

func (info *GeneralInfo) Validate() error {
	return validation.All(
		validation.Comparable(info.UserID, "user_id").Required(true),
		validation.Ptr(&info.Device, "device").With(validation.Custom),
		validation.String(info.AppVersion, "app_version").Required(true),
	)
}

type Device struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	BuildNumber  string `json:"build_number"`
	OS           string `json:"os"`
	OSVersion    string `json:"os_version"`
	ScreenWidth  uint32 `json:"screen_width"`
	ScreenHeight uint32 `json:"screen_height"`
}

func (d *Device) Validate() error {
	return validation.All(
		validation.String(d.Manufacturer, "manufacturer").Required(true),
		validation.String(d.Model, "model").Required(true),
		validation.String(d.BuildNumber, "build_number").Required(true),
	)
}

type Telemetry struct {
	Action    string         `json:"action"`
	Data      map[string]any `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}

func (t *Telemetry) Validate() error {
	return validation.All(
		validation.String(t.Action, "action").Required(true),
		validation.Map(t.Data, "data").NotNil(true),
		validation.Time(t.Timestamp, "timestamp").Required(true),
	)
}

type TrackRequest struct {
	Info GeneralInfo `json:"info"`
	Data []Telemetry `json:"data"`
}

func (tr *TrackRequest) Validate() error {
	return validation.All(
		validation.Ptr(&tr.Info, "info").With(validation.Custom),
		validation.Slice(tr.Data, "data").NotNil(true).ValuesPtrWith(validation.Custom),
	)
}

func main() {
	tr := TrackRequest{}
	err := tr.Validate()
	str, _ := json.Marshal(err)
	fmt.Printf("%s\n", str)
}
