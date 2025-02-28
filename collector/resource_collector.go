package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
	"strconv"
	"strings"
)

type resourceCollector struct {
	props        []string
	descriptions map[string]*prometheus.Desc
}

func newResourceCollector() routerOSCollector {
	c := &resourceCollector{}
	c.init()
	return c
}

func (c *resourceCollector) init() {
	c.props = []string{"free-memory", "total-memory", "cpu-load", "free-hdd-space", "total-hdd-space", "uptime", "board-name", "version"}

	labelNames := []string{"name", "address", "boardname", "version"}
	c.descriptions = make(map[string]*prometheus.Desc)
	for _, p := range c.props[:6] {
		c.descriptions[p] = descriptionForPropertyName("system", p, labelNames)
	}
}

func (c *resourceCollector) describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descriptions {
		ch <- d
	}
}

func (c *resourceCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectForStat(re, ctx)
	}

	return nil
}

func (c *resourceCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/system/resource/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching system resource metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *resourceCollector) collectForStat(re *proto.Sentence, ctx *collectorContext) {
	for _, p := range c.props[:5] {
		c.collectMetricForProperty(p, prometheus.GaugeValue, re, ctx)
	}

	for _, p := range c.props[5:6] {
		c.collectMetricForProperty(p, prometheus.CounterValue, re, ctx)
	}
}

func (c *resourceCollector) collectMetricForProperty(property string, valueType prometheus.ValueType, re *proto.Sentence, ctx *collectorContext) {
	var v float64
	var err error

	boardname := re.Map["board-name"]
	version := re.Map["version"]

	if property == "uptime" {
		v, err = parseDuration(re.Map[property])
	} else {
		v, err = strconv.ParseFloat(re.Map[property], 64)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"device":   ctx.device.Name,
			"property": property,
			"value":    re.Map[property],
			"error":    err,
		}).Error("error parsing system resource metric value")
		return
	}

	desc := c.descriptions[property]
	ctx.ch <- prometheus.MustNewConstMetric(desc, valueType, v, ctx.device.Name, ctx.device.Address, boardname, version)
}
