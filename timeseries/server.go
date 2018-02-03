package timeseries

func Start() {
	go startMetricsHandle()
	go startHealthHandle()
	select {}
}
