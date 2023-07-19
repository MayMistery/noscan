package scanlib

type FingerPrint struct {
	ProbeName        string
	MatchRegexString string

	Service         string
	ProductName     string
	Version         string
	Info            string
	Hostname        string
	OperatingSystem string
	DeviceType      string
}
