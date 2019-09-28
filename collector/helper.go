package collector

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var durationRegex *regexp.Regexp
var durationParts [5]time.Duration
var wirelessRateRegex *regexp.Regexp

func init() {
	durationRegex = regexp.MustCompile(`(?:(\d*)w)?(?:(\d*)d)?(?:(\d*)h)?(?:(\d*)m)?(?:(\d*)s)?`)
	durationParts = [5]time.Duration{time.Hour * 168, time.Hour * 24, time.Hour, time.Minute, time.Second}

	wirelessRateRegex = regexp.MustCompile(`([\d.]+)Mbps.+`)
}

func metricStringCleanup(in string) string {
	return strings.Replace(in, "-", "_", -1)
}

func descriptionForPropertyName(prefix, property string, labelNames []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, prefix, metricStringCleanup(property)),
		property,
		labelNames,
		nil,
	)
}

func description(prefix, name, helpText string, labelNames []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, prefix, name),
		helpText,
		labelNames,
		nil,
	)
}

func splitStringToFloats(metric string) (float64, float64, error) {
	strs := strings.Split(metric, ",")

	m1, err := strconv.ParseFloat(strs[0], 64)
	if err != nil {
		return math.NaN(), math.NaN(), err
	}
	m2, err := strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return math.NaN(), math.NaN(), err
	}
	return m1, m2, nil
}

func parseDuration(duration string) (float64, error) {
	var u time.Duration

	reMatch := durationRegex.FindAllStringSubmatch(duration, -1)

	// should get one and only one match back on the regex
	if len(reMatch) != 1 {
		return 0, fmt.Errorf("invalid duration value sent to regex")
	} else {
		for i, match := range reMatch[0] {
			if match != "" && i != 0 {
				v, err := strconv.Atoi(match)
				if err != nil {
					log.WithFields(log.Fields{
						"duration": duration,
						"value":    match,
						"error":    err,
					}).Error("error parsing duration field value")
					return float64(0), err
				}
				u += time.Duration(v) * durationParts[i-1]
			}
		}
	}
	return u.Seconds(), nil
}

func parseDatetime(datetime string) (time.Time, error) {
	t, err := time.Parse("Jan/02/2006 15:04:05", datetime)
	if err != nil {
		log.WithFields(log.Fields{
			"datetime": datetime,
			"value":    t,
			"error":    err,
		}).Error("error parsing datetime field value")
		return time.Time{}, err
	}

	return t, nil
}

func parseWirelessRate(rate string) (float64, error) {
	reMatch := wirelessRateRegex.FindStringSubmatch(rate)

	// should get one and only one match back on the regex
	if len(reMatch) != 2 {
		return 0, fmt.Errorf("invalid wireless rate value sent to regex")
	} else {
		if reMatch[1] != "" {
			v, err := strconv.ParseFloat(reMatch[1], 64)
			if err != nil {
				log.WithFields(log.Fields{
					"wireless-rate": rate,
					"value":         reMatch[1],
					"error":         err,
				}).Error("error parsing wireless rate field value")
				return 0, err
			}
			return v, nil
		}
	}
	return 0, nil
}
