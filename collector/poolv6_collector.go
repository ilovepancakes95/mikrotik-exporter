package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type poolv6Collector struct {
	usedCountDesc *prometheus.Desc
}

func (c *poolv6Collector) init() {
	const prefix = "ip_pool"

	labelNames := []string{"name", "address", "pool"}
	c.usedCountDesc = description(prefix, "ipv6_used_count", "number of used IP/prefixes in a IPv6 pool", labelNames)
}

func newPoolV6Collector() routerOSCollector {
	c := &poolv6Collector{}
	c.init()
	return c
}

func (c *poolv6Collector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.usedCountDesc
}

func (c *poolv6Collector) collect(ctx *collectorContext) error {
	err := c.collectForIPVersion("6", "ipv6", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *poolv6Collector) collectForIPVersion(ipVersion, topic string, ctx *collectorContext) error {
	names, err := c.fetchPoolNames(ipVersion, topic, ctx)
	if err != nil {
		return err
	}

	for _, n := range names {
		err := c.collectForPool(ipVersion, topic, n, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *poolv6Collector) fetchPoolNames(ipVersion, topic string, ctx *collectorContext) ([]string, error) {
	reply, err := ctx.client.Run(fmt.Sprintf("/%s/pool/print", topic), "=.proplist=name")
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching pool names")
		return nil, err
	}

	names := make([]string, len(reply.Re))
	for i, re := range reply.Re {
		names[i] = re.Map["name"]
	}

	return names, nil
}

func (c *poolv6Collector) collectForPool(ipVersion, topic, pool string, ctx *collectorContext) error {
	reply, err := ctx.client.Run(fmt.Sprintf("/%s/pool/used/print", topic), fmt.Sprintf("?pool=%s", pool), "=count-only=")
	if err != nil {
		log.WithFields(log.Fields{
			"pool":       pool,
			"ip_version": ipVersion,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error fetching pool counts")
		return err
	}

	v, err := strconv.ParseFloat(reply.Done.Map["ret"], 32)
	if err != nil {
		log.WithFields(log.Fields{
			"pool":       pool,
			"ip_version": ipVersion,
			"device":     ctx.device.Name,
			"error":      err,
		}).Error("error parsing pool counts")
		return err
	}

	ctx.ch <- prometheus.MustNewConstMetric(c.usedCountDesc, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address, pool)
	return nil
}
