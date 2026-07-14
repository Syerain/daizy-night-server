package constants

// the ways to register a user -- indicating params needed for db to create a user
type registerWayValue string

const (
	Legacy registerWayValue = "legacy"
	Github registerWayValue = "oauth-github"
)

// more attributes may be added in the future
type RegisterWay struct {
	Value registerWayValue
}

// var Legacy = RegisterWay{
// 	value: legacy,
// }

// var Github = RegisterWay{
// 	value: github,
// }

// currently meaningless
// func (r RegisterWay) Is(v *RegisterWay) bool {
// 	return r.value == v.value
// }

// func (r RegisterWay) GetType() registerWayValue {
// 	return r.value
// }
