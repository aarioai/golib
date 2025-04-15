package enumz

type LawyerCertType byte

const (
	LawyerCertA = 'A' // A 证
	LawyerCertB = 'B'
	LawyerCertC = 'C'
)

func NewLawyerCertType(t byte) (LawyerCertType, bool) {
	if t >= LawyerCertA && t <= LawyerCertC {
		return LawyerCertType(t), true
	}
	return LawyerCertType(0), false
}
func (t LawyerCertType) String() string {
	return string(t)
}
