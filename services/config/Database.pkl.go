// Code generated from Pkl module `myorg.myteam.AppConfig`. DO NOT EDIT.
package config

type Database struct {
	// Username for the database
	Username string `pkl:"username"`

	// Password for the database
	Password string `pkl:"password"`

	// Hostname for the database
	Hostname string `pkl:"hostname"`

	// Port for the database
	Port uint16 `pkl:"port"`

	// Name of the database
	Name string `pkl:"name"`
}
