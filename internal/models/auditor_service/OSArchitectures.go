package auditor_service

type OSArchitectures struct {
	DarwinAmd64 FactOSArchitecture `json:"darwin_amd64"`
	LinuxArm    FactOSArchitecture `json:"linux_arm"`
}
