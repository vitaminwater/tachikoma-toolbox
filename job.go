package tachikoma

type Job interface {
	GetName() string
	Run(interface{}) error
}
