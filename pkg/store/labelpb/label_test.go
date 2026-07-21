// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package labelpb

import (
	"fmt"
	"io"
	"testing"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/efficientgo/core/testutil"
)

func TestExtendLabels(t *testing.T) {
	testutil.Equals(t,
		labels.FromStrings("a", "1", "replica", "01", "xb", "2"),
		ExtendSortedLabels(
			labels.FromStrings("a", "1", "xb", "2"),
			labels.FromStrings("replica", "01"),
		))

	testutil.Equals(t,
		labels.FromStrings("replica", "01"),
		ExtendSortedLabels(
			labels.EmptyLabels(),
			labels.FromStrings("replica", "01"),
		))

	testutil.Equals(t, labels.FromStrings("a", "1", "replica", "01", "xb", "2"),
		ExtendSortedLabels(
			labels.FromStrings("a", "1", "replica", "NOT01", "xb", "2"),
			labels.FromStrings("replica", "01"),
		))

	testInjectExtLabels(testutil.NewTB(t))
}

func BenchmarkRealloc(b *testing.B) {
	lbls := ZLabelsFromPromLabels(labels.FromStrings(
		"__name__", "container_cpu_usage_seconds_total",
		"cluster", "prod-eu-west-1",
		"container", "thanos-receive",
		"cpu", "total",
		"endpoint", "https-metrics",
		"image", "quay.io/thanos/thanos:v0.35.0",
		"instance", "10.128.4.231:10250",
		"job", "kubelet",
		"metrics_path", "/metrics/cadvisor",
		"namespace", "monitoring",
		"node", "ip-10-128-4-231.eu-west-1.compute.internal",
		"pod", "thanos-receive-default-0",
		"prometheus", "monitoring/k8s",
		"prometheus_replica", "prometheus-k8s-0",
		"service", "kubelet",
	))
	work := make([]ZLabel, len(lbls))

	for _, intern := range []bool{false, true} {
		b.Run(fmt.Sprintf("intern=%v", intern), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				copy(work, lbls)
				ReAllocZLabelsStrings(&work, intern)
			}
		})
	}
}

func BenchmarkExtendLabels(b *testing.B) {
	testInjectExtLabels(testutil.NewTB(b))
}

var x labels.Labels

func testInjectExtLabels(tb testutil.TB) {
	in := labels.FromStrings(
		"__name__", "subscription_labels",
		"_id", "0dfsdfsdsfdsffd1e96-4432-9abe-e33436ea969a",
		"account", "1afsdfsddsfsdfsdfsdfsdfs",
		"ebs_account", "1asdasdad45",
		"email_domain", "asdasddgfkw.example.com",
		"endpoint", "metrics",
		"external_organization", "dfsdfsdf",
		"instance", "10.128.4.231:8080",
		"job", "sdd-acct-mngr-metrics",
		"managed", "false",
		"namespace", "production",
		"organization", "dasdadasdasasdasaaFGDSG",
		"pod", "sdd-acct-mngr-6669c947c8-xjx7f",
		"prometheus", "telemeter-production/telemeter",
		"prometheus_replica", "prometheus-telemeter-1",
		"risk", "5",
		"service", "sdd-acct-mngr-metrics",
		"support", "Self-Support", // Should be overwritten.
	)
	extLset := labels.FromMap(map[string]string{
		"replica": "1",
		"support": "Host-Support",
		"tenant":  "2342",
	})
	tb.ResetTimer()
	for i := 0; i < tb.N(); i++ {
		x = ExtendSortedLabels(in, extLset)

		if !tb.IsBenchmark() {
			testutil.Equals(tb, labels.FromStrings(
				"__name__", "subscription_labels",
				"_id", "0dfsdfsdsfdsffd1e96-4432-9abe-e33436ea969a",
				"account", "1afsdfsddsfsdfsdfsdfsdfs",
				"ebs_account", "1asdasdad45",
				"email_domain", "asdasddgfkw.example.com",
				"endpoint", "metrics",
				"external_organization", "dfsdfsdf",
				"instance", "10.128.4.231:8080",
				"job", "sdd-acct-mngr-metrics",
				"managed", "false",
				"namespace", "production",
				"organization", "dasdadasdasasdasaaFGDSG",
				"pod", "sdd-acct-mngr-6669c947c8-xjx7f",
				"prometheus", "telemeter-production/telemeter",
				"prometheus_replica", "prometheus-telemeter-1",
				"replica", "1",
				"risk", "5",
				"service", "sdd-acct-mngr-metrics",
				"support", "Host-Support",
				"tenant", "2342",
			), x)
		}
	}
	fmt.Fprint(io.Discard, x)
}
