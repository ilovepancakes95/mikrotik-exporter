package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type routesV6Collector struct {
	protocols         []string
	countDesc         *prometheus.Desc
	countProtocolDesc *prometheus.Desc
}

func newRoutesV6Collector() routerOSCollector {
	c := &routesV6Collector{}
	c.init()
	return c
}

func (c *routesV6Collector) init() {
	const prefix = "routes"
	labelNames := []string{"name", "address"}
	c.countDesc = description(prefix, "ipv6_total_count", "number of IPv6 routes in RIB", labelNames)
	c.countProtocolDesc = description(prefix, "ipv6_protocol_count", "number of IPv6 routes per protocol in RIB", append(labelNames, "protocol"))

	c.protocols = []string{"bgp", "static", "ospf", "dynamic", "connect"}
}

func (c *routesV6Collector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.countDesc
	ch <- c.countProtocolDesc
}

func (c *routesV6Collector) collect(ctx *collectorContext) error {
	err := c.collectForIPVersion("6", "ipv6", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *routesV6Collector) collectForIPVersion(ipVersion, topic string, ctx *collectorContext) error {
	err := c.collectCount(ipVersion, topic, ctx)
	if err != nil {
		return err
	}

	for _, p := range c.protocols {
		err := c.collectCountProtocol(ipVersion, topic, p, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *routesV6Collector) collectCount(ipVersion, topic string, ctx *collectorContext) error {
	reply, err := ctx.client.Run(fmt.Sprintf("/%s/route/print", topic), "?disabled=false", "=count-only=")
	if err != nil {
		log.WithFields(log.Fields{
			"ip_version": ipVersion,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error fetching routes metrics")
		return err
	}

	v, err := strconv.ParseFloat(reply.Done.Map["ret"], 32)
	if err != nil {
		log.WithFields(log.Fields{
			"ip_version": ipVersion,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error parsing routes metrics")
		return err
	}

	ctx.ch <- prometheus.MustNewConstMetric(c.countDesc, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address)
	return nil
}

func (c *routesV6Collector) collectCountProtocol(ipVersion, topic, protocol string, ctx *collectorContext) error {
	reply, err := ctx.client.Run(fmt.Sprintf("/%s/route/print", topic), "?disabled=false", fmt.Sprintf("?%s", protocol), "=count-only=")
	if err != nil {
		log.WithFields(log.Fields{
			"ip_version": ipVersion,
			"protocol":   protocol,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error fetching routes metrics")
		return err
	}

	v, err := strconv.ParseFloat(reply.Done.Map["ret"], 32)
	if err != nil {
		log.WithFields(log.Fields{
			"ip_version": ipVersion,
			"protocol":   protocol,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error parsing routes metrics")
		return err
	}

	ctx.ch <- prometheus.MustNewConstMetric(c.countProtocolDesc, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address, protocol)
	return nil
}
