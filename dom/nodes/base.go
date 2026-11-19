package nodes

type baseNode struct{}

func (b baseNode) String() string {
	return ""
}
func (b baseNode) Value() *string {
	return nil
}
