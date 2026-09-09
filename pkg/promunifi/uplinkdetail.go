package promunifi

import (
	"github.com/unpoller/unifi/v6"
)

// exportUplinkDetail emits uplink topology that the controller sends and the
// unifi library parses, but that the exporter previously discarded: the parent
// device MAC, the parent device name, and the link type (wire/wireless).
//
// The gateway/switch uplink (unifi.Uplink) and the access point uplink are
// distinct struct shapes in the library -- notably the AP uplink carries
// uplink_mac, which is the parent device a mesh AP is backhauled to -- so the
// caller passes the fields explicitly rather than a struct.
//
// labels is the standard device label slice in descUSG order:
// {port, site_name, name, source, tag}.
func (u *promUnifi) exportUplinkDetail(
	r report,
	labels []string,
	uplinkMac, uplinkDevice, uplinkName, uplinkType, uplinkMedia string,
	up bool,
	txRate, rxRate, txBytesR, rxBytesR unifi.FlexInt,
) {
	if uplinkMac == "" && uplinkName == "" && uplinkType == "" {
		return // device carries no uplink object
	}

	linkUp := 0.0
	if up {
		linkUp = 1.0
	}

	ulLabels := []string{
		labels[1], labels[2], labels[3], labels[4],
		uplinkMac, uplinkDevice, uplinkName, uplinkType, uplinkMedia,
	}

	r.send([]*metric{
		// Info metrics are a constant 1; the labels are the payload. Link
		// state is a separate series so it can be alerted on.
		{u.USG.UplinkInfo, gauge, 1.0, ulLabels},
		{u.USG.UplinkUp, gauge, linkUp, ulLabels},
		{u.USG.UplinkTxRate, gauge, txRate, ulLabels},
		{u.USG.UplinkRxRate, gauge, rxRate, ulLabels},
		{u.USG.UplinkTxBytesRate, gauge, txBytesR, ulLabels},
		{u.USG.UplinkRxBytesRate, gauge, rxBytesR, ulLabels},
	})
}
