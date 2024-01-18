package sentinel

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/system"
)

func InitRule() error {
	_, err := flow.LoadRules([]*flow.Rule{ //流控规则
		{Resource: "flowDR", TokenCalculateStrategy: flow.Direct, ControlBehavior: flow.Reject, Threshold: 10},                                                                        //策略：直接模式，超QPS后，直接失败
		{Resource: "flowDT", TokenCalculateStrategy: flow.Direct, ControlBehavior: flow.Throttling, Threshold: 10, MaxQueueingTimeMs: 1000},                                           //策略：直接模式，匀速排队，排队间隔在1000/Threshold ms，超出最长排队等待时间后失败
		{Resource: "flowWR", TokenCalculateStrategy: flow.WarmUp, WarmUpPeriodSec: 10, WarmUpColdFactor: 3, ControlBehavior: flow.Reject, Threshold: 10},                              //策略：预热模式，预热10秒，预热因子默认为3(通过预热公式得出，因子需大于1)，超QPS后，直接失败
		{Resource: "flowWT", TokenCalculateStrategy: flow.WarmUp, WarmUpPeriodSec: 10, WarmUpColdFactor: 3, ControlBehavior: flow.Throttling, Threshold: 10, MaxQueueingTimeMs: 1000}, //策略：预热模式，预热10秒，预热因子默认为3(通过预热公式得出，因子需大于1)，匀速排队，排队间隔在1000/Threshold ms，超出最长排队等待时间后失败
	})
	if err != nil {
		return err
	}
	_, err = isolation.LoadRules([]*isolation.Rule{ //流量隔离规则
		{Resource: "rule2", MetricType: isolation.Concurrency, Threshold: 15},
	})
	if err != nil {
		return err
	}
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{ //熔断降级
		{Resource: "rule3", Strategy: circuitbreaker.SlowRequestRatio, RetryTimeoutMs: 500, MinRequestAmount: 20},
	})
	if err != nil {
		return err
	}
	_, err = system.LoadRules([]*system.Rule{ //系统自适应流控
		{MetricType: system.Load, TriggerCount: 4, Strategy: system.BBR},
		{MetricType: system.CpuUsage, TriggerCount: 0.1, Strategy: system.BBR},
	})
	if err != nil {
		return err
	}
	return nil
}
