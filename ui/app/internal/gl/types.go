// +build !js

package gl

type (
	Buffer       struct{ V uint }
	Framebuffer  struct{ V uint }
	Program      struct{ V uint }
	Renderbuffer struct{ V uint }
	Shader       struct{ V uint }
	Texture      struct{ V uint }
	Query        struct{ V uint }
	Uniform      struct{ V int }
	Object       struct{ V uint }
)

func (u Uniform) Valid() bool {
	return u.V != -1
}
