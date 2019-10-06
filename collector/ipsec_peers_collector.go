package collector

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type ipsecPeersCollector struct {
	props        []string
	descriptions map[string]*prometheus.Desc
}

func newIPSecPeersCollector() routerOSCollector {
	c := &ipsecPeersCollector{}
	c.init()
	return c
}

func (c *ipsecPeersCollector) init() {
	c.props = []string{"id", "remote-address", "state", "uptime", "rx-bytes", "tx-bytes", "rx-packets", "tx-packets"}

	const prefix = "ipsec_peers"
	labelNames := []string{"name", "address", "id", "remote_address", "state"}

	c.descriptions = make(map[string]*prometheus.Desc)

	for _, p := range c.props[3:] {
		c.descriptions[p] = descriptionForPropertyName(prefix, p, labelNames)
	}
}

func (c *ipsecPeersCollector) describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descriptions {
		ch <- d
	}
}

func (c *ipsecPeersCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectForStat(re, ctx)
	}

	return nil
}

func (c *ipsecPeersCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/ip/ipsec/active-peers/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching ipsec peers metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *ipsecPeersCollector) collectForStat(re *proto.Sentence, ctx *collectorContext) {
	remoteAddress := re.Map["remote-address"]
	id := re.Map["id"]
	state := re.Map["state"]

	for _, p := range c.props[3:] {
		c.collectMetricForProperty(p, id, remoteAddress, state, re, ctx)
	}
}

func (c *ipsecPeersCollector) collectMetricForProperty(property, id, remoteAddress, state string, re *proto.Sentence, ctx *collectorContext) {
	desc := c.descriptions[property]
	v, err := c.parseValueForProperty(property, re.Map[property])
	if err != nil {
		log.WithFields(log.Fields{
			"device":   ctx.device.Name,
			"id":       id,
			"property": property,
			"value":    re.Map[property],
			"error":    err,
		}).Error("error parsing ipsec peers metric value")
		return
	}

	valueType := prometheus.CounterValue
	switch property {
	case "state", "uptime":
		valueType = prometheus.GaugeValue
	}

	ctx.ch <- prometheus.MustNewConstMetric(desc, valueType, v, ctx.device.Name, ctx.device.Address, id, remoteAddress, state)
}

func (c *ipsecPeersCollector) parseValueForProperty(property, value string) (float64, error) {
	switch property {
	case "state":
		if value == "established" {
			return 1, nil
		}

		return 0, nil
	case "uptime":
		return parseDuration(value)
	default:
		if value == "" {
			return 0, nil
		}

		return strconv.ParseFloat(value, 64)
	}
}
