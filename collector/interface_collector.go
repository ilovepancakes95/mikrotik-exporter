package collector

import (
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type interfaceCollector struct {
	props        []string
	descriptions map[string]*prometheus.Desc
}

func newInterfaceCollector() routerOSCollector {
	c := &interfaceCollector{}
	c.init()
	return c
}

func (c *interfaceCollector) init() {
	c.props = []string{"name", "comment", "mac-address", "type", "last-link-down-time", "last-link-up-time", "running", "actual-mtu", "link-downs", "rx-byte", "tx-byte", "rx-packet", "tx-packet", "rx-error", "tx-error", "rx-drop", "tx-drop"}

	labelNames := []string{"name", "address", "interface", "comment", "mac_address", "type"}
	c.descriptions = make(map[string]*prometheus.Desc)
	for _, p := range c.props[4:] {
		c.descriptions[p] = descriptionForPropertyName("interface", p, labelNames)
	}
}

func (c *interfaceCollector) describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descriptions {
		ch <- d
	}
}

func (c *interfaceCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	clock, err := c.fetchDatetime(ctx)
	if err != nil {
		return err
	}

	datetime, err := time.Parse("Jan/02/2006 15:04:05", clock[0].Map["date"]+" "+clock[0].Map["time"])
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectForStat(re, datetime, ctx)
	}

	return nil
}

func (c *interfaceCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/interface/print", "?disabled=false", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching interface metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *interfaceCollector) fetchDatetime(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/system/clock/print", "=.proplist=time,date")
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching clock metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *interfaceCollector) collectForStat(re *proto.Sentence, datetime time.Time, ctx *collectorContext) {
	name := re.Map["name"]
	comment := re.Map["comment"]
	macAddress := re.Map["mac-address"]
	interfaceType := re.Map["type"]

	for _, p := range c.props[4:8] {
		c.collectMetricForProperty(p, name, comment, macAddress, interfaceType, datetime, prometheus.GaugeValue, re, ctx)
	}

	for _, p := range c.props[8:] {
		c.collectMetricForProperty(p, name, comment, macAddress, interfaceType, datetime, prometheus.CounterValue, re, ctx)
	}
}

func (c *interfaceCollector) collectMetricForProperty(property, iface, comment, macAddress, interfaceType string, datetime time.Time, valueType prometheus.ValueType, re *proto.Sentence, ctx *collectorContext) {
	desc := c.descriptions[property]
	if value := re.Map[property]; value != "" {
		var v float64
		var err error
		if property == "last-link-down-time" || property == "last-link-up-time" {
			var t time.Time
			t, err = parseDatetime(value)
			if err != nil {
				log.WithFields(log.Fields{
					"device":    ctx.device.Name,
					"interface": iface,
					"property":  property,
					"value":     value,
					"error":     err,
				}).Error("error parsing interface duration metric value")
				return
			}

			v = datetime.Sub(t).Seconds()
		} else if property == "running" {
			if value == "true" {
				v = 1
			}
		} else {
			v, err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.WithFields(log.Fields{
					"device":    ctx.device.Name,
					"interface": iface,
					"property":  property,
					"value":     value,
					"error":     err,
				}).Error("error parsing interface metric value")
				return
			}
		}

		ctx.ch <- prometheus.MustNewConstMetric(desc, valueType, v, ctx.device.Name, ctx.device.Address, iface, comment, macAddress, interfaceType)
	}
}
