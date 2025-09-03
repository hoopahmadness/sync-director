package device

type deviceManager interface {
	GetDeviceById(id string) (*Device, bool)
}
