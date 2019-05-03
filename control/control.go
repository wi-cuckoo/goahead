package control

// Controller subprocess
type Controller interface {
	Start(*Unit) error // start an unit

	Stop(name string) error // stop an unit with name

	Status(name string) (*Status, error) // check status with name

	Reload(name string) error // reload an unit when changing its config
}
