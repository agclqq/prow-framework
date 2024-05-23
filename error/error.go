package pubErr

import "errors"

var ErrDuplicateName = errors.New("名称(name)重复")

var ErrInstanceStatusOutOfDate = errors.New("实例状态信息是过时的")

type ErrLabelType string

const (
	ErrLabelRedis  ErrLabelType = "[Redis Error]"
	ErrLabelK8S    ErrLabelType = "[K8S Error]"
	ErrLabelHarbor ErrLabelType = "[Harbor Error]"
	ErrLabelCP     ErrLabelType = "[CommonPipeline Error]"
)
