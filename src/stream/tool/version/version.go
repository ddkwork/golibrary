package version

import (
	"github.com/ddkwork/golibrary/src/stream/tool/strconv"
	"strings"
)

type (
	Interface interface {
		String() string
		VerSion(version string) (ok bool)
		Major() uint64
		SetMajor(major uint64)
		Minor() uint64
		SetMinor(minor uint64)
		Patch() uint64
		SetPatch(patch uint64)
	}
	object struct {
		major uint64
		minor uint64
		patch uint64
	}
)

func New() Interface                    { return &object{} }
func (o *object) SetPatch(patch uint64) { o.patch = patch }
func (o *object) SetMinor(minor uint64) { o.minor = minor }
func (o *object) SetMajor(major uint64) { o.major = major }
func (o *object) Major() uint64         { return o.major }
func (o *object) Minor() uint64         { return o.minor }
func (o *object) Patch() uint64         { return o.patch }

func (o *object) String() string {
	s := strconv.Default
	var array []string
	array = append(
		array,
		s.FormatUint(o.Major(), 10),
		s.FormatUint(o.Minor(), 10),
		s.FormatUint(o.Patch(), 10),
	)
	return strings.Join(array, `.`)
}

func (o *object) VerSion(version string) (ok bool) {
	s := strconv.Default
	split := strings.Split(version, `.`)
	if !s.ParseUint(split[0], 10, 32) {
		return
	}
	o.SetMajor(s.Value())
	if !s.ParseUint(split[1], 10, 32) {
		return
	}
	o.SetMinor(s.Value())
	if !s.ParseUint(split[2], 10, 32) {
		return
	}
	o.SetPatch(s.Value())
	return true
}
