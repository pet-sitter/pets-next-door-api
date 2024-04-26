package resourcemedia

type ResourceType string

const (
	SOSResourceType ResourceType = "sos_posts"
)

func (r ResourceType) String() string {
	return string(r)
}
