package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
	"strconv"
	"strings"
)

type dhcpLeaseCollector struct {
	props        []string
	descriptions *prometheus.Desc
}

func (c *dhcpLeaseCollector) init() {
	c.props = []string{"active-mac-address", "status", "expires-after", "active-address", "host-name"}

	labelNames := []string{"name", "address", "active_mac_address", "status", "expires_after", "active_address", "hostname"}
	c.descriptions = description("dhcp", "leases_info", "DHCP lease info", labelNames)

}

func newDHCPLCollector() routerOSCollector {
	c := &dhcpLeaseCollector{}
	c.init()
	return c
}

func (c *dhcpLeaseCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.descriptions
}

func (c *dhcpLeaseCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectMetric(ctx, re)
	}

	return nil
}

func (c *dhcpLeaseCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/ip/dhcp-server/lease/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching DHCP leases metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *dhcpLeaseCollector) collectMetric(ctx *collectorContext, re *proto.Sentence) {
	v := 1.0

	f, err := parseDuration(re.Map["expires-after"])
	if err != nil {
		log.WithFields(log.Fields{
			"device":   ctx.device.Name,
			"property": "expires-after",
			"value":    re.Map["expires-after"],
			"error":    err,
		}).Error("error parsing duration metric value")
		return
	}

	activemacaddress := re.Map["active-mac-address"]
	status := re.Map["status"]
	activeaddress := re.Map["active-address"]
	hostname := re.Map["host-name"]

	ctx.ch <- prometheus.MustNewConstMetric(c.descriptions, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address, activemacaddress, status, strconv.FormatFloat(f, 'f', 0, 64), activeaddress, hostname)
}
