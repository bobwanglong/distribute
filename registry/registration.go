package registry

type Registration struct {
	ServiceName ServiceName
	ServiceUrl  string
	RequiredServices []ServiceName // 依赖服务
	ServiceUpdateURL string // 提供通知更新的URL
	HeartbeatURL string
}

type ServiceName string

const (
	LogService = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	PortalService  = ServiceName("Portald")

)

// 每一条更新
type patchEntry struct{
	Name ServiceName
	URL string
}

// 
type patch struct{
	Added []patchEntry
	Removed []patchEntry
}