// SPDX-License-Identifier: Unlicense OR MIT

// +build darwin linux

package gl

import (
	"unsafe"
)

/*
#cgo linux LDFLAGS: -lGLESv2 -ldl
#cgo darwin,!ios LDFLAGS: -framework OpenGL

#include <stdlib.h>

#ifdef __APPLE__
#cgo CFLAGS: -DGL_SILENCE_DEPRECATION
	#include "TargetConditionals.h"
	#if TARGET_OS_IPHONE
	#include <OpenGLES/ES3/gl.h>
	#else
	#include <OpenGL/gl3.h>
	#endif
#else
#define __USE_GNU
#include <dlfcn.h>
#include <GLES2/gl2.h>
#include <GLES3/gl3.h>
#endif

static void (*_glInvalidateFramebuffer)(GLenum target, GLsizei numAttachments, const GLenum *attachments);

static void (*_glBeginQuery)(GLenum target, GLuint id);
static void (*_glDeleteQueries)(GLsizei n, const GLuint *ids);
static void (*_glEndQuery)(GLenum target);
static void (*_glGenQueries)(GLsizei n, GLuint *ids);
static void (*_glGetQueryObjectuiv)(GLuint id, GLenum pname, GLuint *params);

// The pointer-free version of glVertexAttribPointer, to avoid the Cgo pointer checks.
__attribute__ ((visibility ("hidden"))) void gio_glVertexAttribPointer(GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, uintptr_t offset) {
	glVertexAttribPointer(index, size, type, normalized, stride, (const GLvoid *)offset);
}

// The pointer-free version of glDrawElements, to avoid the Cgo pointer checks.
__attribute__ ((visibility ("hidden"))) void gio_glDrawElements(GLenum mode, GLsizei count, GLenum type, const uintptr_t offset) {
	glDrawElements(mode, count, type, (const GLvoid *)offset);
}

__attribute__ ((visibility ("hidden"))) void gio_glInvalidateFramebuffer(GLenum target, GLenum attachment) {
	// Framebuffer invalidation is just a hint and can safely be ignored.
	if (_glInvalidateFramebuffer != NULL) {
		_glInvalidateFramebuffer(target, 1, &attachment);
	}
}

__attribute__ ((visibility ("hidden"))) void gio_glBeginQuery(GLenum target, GLenum attachment) {
	_glBeginQuery(target, attachment);
}

__attribute__ ((visibility ("hidden"))) void gio_glDeleteQueries(GLsizei n, const GLuint *ids) {
	_glDeleteQueries(n, ids);
}

__attribute__ ((visibility ("hidden"))) void gio_glEndQuery(GLenum target) {
	_glEndQuery(target);
}

__attribute__ ((visibility ("hidden"))) void gio_glGenQueries(GLsizei n, GLuint *ids) {
	_glGenQueries(n, ids);
}

__attribute__ ((visibility ("hidden"))) void gio_glGetQueryObjectuiv(GLuint id, GLenum pname, GLuint *params) {
	_glGetQueryObjectuiv(id, pname, params);
}

__attribute__((constructor)) static void gio_loadGLFunctions() {
#ifdef __APPLE__
	#if TARGET_OS_IPHONE
	_glInvalidateFramebuffer = glInvalidateFramebuffer;
	_glBeginQuery = glBeginQuery;
	_glDeleteQueries = glDeleteQueries;
	_glEndQuery = glEndQuery;
	_glGenQueries = glGenQueries;
	_glGetQueryObjectuiv = glGetQueryObjectuiv;
	#endif
#else
	// Load libGLESv3 if available.
	dlopen("libGLESv3.so", RTLD_NOW | RTLD_GLOBAL);
	_glInvalidateFramebuffer = dlsym(RTLD_DEFAULT, "glInvalidateFramebuffer");
	// Fall back to EXT_invalidate_framebuffer if available.
	if (_glInvalidateFramebuffer == NULL) {
		_glInvalidateFramebuffer = dlsym(RTLD_DEFAULT, "glDiscardFramebufferEXT");
	}

	_glBeginQuery = dlsym(RTLD_DEFAULT, "glBeginQuery");
	if (_glBeginQuery == NULL)
		_glBeginQuery = dlsym(RTLD_DEFAULT, "glBeginQueryEXT");
	_glDeleteQueries = dlsym(RTLD_DEFAULT, "glDeleteQueries");
	if (_glDeleteQueries == NULL)
		_glDeleteQueries = dlsym(RTLD_DEFAULT, "glDeleteQueriesEXT");
	_glEndQuery = dlsym(RTLD_DEFAULT, "glEndQuery");
	if (_glEndQuery == NULL)
		_glEndQuery = dlsym(RTLD_DEFAULT, "glEndQueryEXT");
	_glGenQueries = dlsym(RTLD_DEFAULT, "glGenQueries");
	if (_glGenQueries == NULL)
		_glGenQueries = dlsym(RTLD_DEFAULT, "glGenQueriesEXT");
	_glGetQueryObjectuiv = dlsym(RTLD_DEFAULT, "glGetQueryObjectuiv");
	if (_glGetQueryObjectuiv == NULL)
		_glGetQueryObjectuiv = dlsym(RTLD_DEFAULT, "glGetQueryObjectuivEXT");
#endif
}
*/
import "C"

type Functions struct{}

func (f *Functions) ActiveTexture(texture Enum) {
	C.glActiveTexture(C.GLenum(texture))
}

func (f *Functions) AttachShader(p Program, s Shader) {
	C.glAttachShader(C.GLuint(p), C.GLuint(s))
}

func (f *Functions) BeginQuery(target Enum, query Query) {
	C.gio_glBeginQuery(C.GLenum(target), C.GLenum(query))
}

func (f *Functions) BindAttribLocation(p Program, a Attrib, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.glBindAttribLocation(C.GLuint(p), C.GLuint(a), cname)
}

func (f *Functions) BindBuffer(target Enum, b Buffer) {
	C.glBindBuffer(C.GLenum(target), C.GLuint(b))
}

func (f *Functions) BindFramebuffer(target Enum, fb Framebuffer) {
	C.glBindFramebuffer(C.GLenum(target), C.GLuint(fb))
}

func (f *Functions) BindRenderbuffer(target Enum, fb Renderbuffer) {
	C.glBindRenderbuffer(C.GLenum(target), C.GLuint(fb))
}

func (f *Functions) BindTexture(target Enum, t Texture) {
	C.glBindTexture(C.GLenum(target), C.GLuint(t))
}

func (f *Functions) BlendEquation(mode Enum) {
	C.glBlendEquation(C.GLenum(mode))
}

func (f *Functions) BlendFunc(sfactor, dfactor Enum) {
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
}

func (f *Functions) BufferData(target Enum, src []byte, usage Enum) {
	var p unsafe.Pointer
	if len(src) > 0 {
		p = unsafe.Pointer(&src[0])
	}
	C.glBufferData(C.GLenum(target), C.GLsizeiptr(len(src)), p, C.GLenum(usage))
}

func (f *Functions) CheckFramebufferStatus(target Enum) Enum {
	return Enum(C.glCheckFramebufferStatus(C.GLenum(target)))
}

func (f *Functions) Clear(mask Enum) {
	C.glClear(C.GLbitfield(mask))
}

func (f *Functions) ClearColor(red float32, green float32, blue float32, alpha float32) {
	C.glClearColor(C.GLfloat(red), C.GLfloat(green), C.GLfloat(blue), C.GLfloat(alpha))
}

func (f *Functions) ClearDepthf(d float32) {
	C.glClearDepthf(C.GLfloat(d))
}

func (f *Functions) CompileShader(s Shader) {
	C.glCompileShader(C.GLuint(s))
}

func (f *Functions) CreateBuffer() Buffer {
	var handle C.GLuint
	C.glGenBuffers(1, &handle)
	return Buffer(handle)
}

func (f *Functions) CreateFramebuffer() Framebuffer {
	var handle C.GLuint
	C.glGenFramebuffers(1, &handle)
	return Framebuffer(handle)
}

func (f *Functions) CreateProgram() Program {
	return Program(C.glCreateProgram())
}

func (f *Functions) CreateQuery() Query {
	var handle C.GLuint
	C.gio_glGenQueries(1, &handle)
	return Query(handle)
}

func (f *Functions) CreateRenderbuffer() Renderbuffer {
	var handle C.GLuint
	C.glGenRenderbuffers(1, &handle)
	return Renderbuffer(handle)
}

func (f *Functions) CreateShader(ty Enum) Shader {
	return Shader(C.glCreateShader(C.GLenum(ty)))
}

func (f *Functions) CreateTexture() Texture {
	var handle C.GLuint
	C.glGenTextures(1, &handle)
	return Texture(handle)
}

func (f *Functions) DeleteBuffer(v Buffer) {
	handle := C.GLuint(v)
	C.glDeleteBuffers(1, &handle)
}

func (f *Functions) DeleteFramebuffer(v Framebuffer) {
	handle := C.GLuint(v)
	C.glDeleteFramebuffers(1, &handle)
}

func (f *Functions) DeleteProgram(p Program) {
	C.glDeleteProgram(C.GLuint(p))
}

func (f *Functions) DeleteQuery(query Query) {
	handle := C.GLuint(query)
	C.gio_glDeleteQueries(1, &handle)
}

func (f *Functions) DeleteRenderbuffer(v Renderbuffer) {
	handle := C.GLuint(v)
	C.glDeleteRenderbuffers(1, &handle)
}

func (f *Functions) DeleteShader(s Shader) {
	C.glDeleteShader(C.GLuint(s))
}

func (f *Functions) DeleteTexture(v Texture) {
	handle := C.GLuint(v)
	C.glDeleteTextures(1, &handle)
}

func (f *Functions) DepthFunc(v Enum) {
	C.glDepthFunc(C.GLenum(v))
}

func (f *Functions) DepthMask(mask bool) {
	m := C.GLboolean(C.GL_FALSE)
	if mask {
		m = C.GLboolean(C.GL_TRUE)
	}
	C.glDepthMask(m)
}

func (f *Functions) DisableVertexAttribArray(a Attrib) {
	C.glDisableVertexAttribArray(C.GLuint(a))
}

func (f *Functions) Disable(cap Enum) {
	C.glDisable(C.GLenum(cap))
}

func (f *Functions) DrawArrays(mode Enum, first int, count int) {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}

func (f *Functions) DrawElements(mode Enum, count int, ty Enum, offset int) {
	C.gio_glDrawElements(C.GLenum(mode), C.GLsizei(count), C.GLenum(ty), C.uintptr_t(offset))
}

func (f *Functions) Enable(cap Enum) {
	C.glEnable(C.GLenum(cap))
}

func (f *Functions) EndQuery(target Enum) {
	C.gio_glEndQuery(C.GLenum(target))
}

func (f *Functions) EnableVertexAttribArray(a Attrib) {
	C.glEnableVertexAttribArray(C.GLuint(a))
}

func (f *Functions) Finish() {
	C.glFinish()
}

func (f *Functions) FramebufferRenderbuffer(target, attachment, renderbuffertarget Enum, renderbuffer Renderbuffer) {
	C.glFramebufferRenderbuffer(C.GLenum(target), C.GLenum(attachment), C.GLenum(renderbuffertarget), C.GLuint(renderbuffer))
}

func (f *Functions) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	C.glFramebufferTexture2D(C.GLenum(target), C.GLenum(attachment), C.GLenum(texTarget), C.GLuint(t), C.GLint(level))
}

func (f *Functions) GetError() Enum {
	return Enum(C.glGetError())
}

func (f *Functions) GetRenderbufferParameteri(target, pname Enum) int {
	// Hope this is enough room.
	var buf [100]C.GLint
	C.glGetRenderbufferParameteriv(C.GLenum(target), C.GLenum(pname), &buf[0])
	return int(buf[0])
}

func (f *Functions) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	// Hope this is enough room.
	var buf [100]C.GLint
	C.glGetFramebufferAttachmentParameteriv(C.GLenum(target), C.GLenum(attachment), C.GLenum(pname), &buf[0])
	return int(buf[0])
}

func (f *Functions) GetInteger(pname Enum) int {
	// Hope this is enough room.
	var buf [100]C.GLint
	C.glGetIntegerv(C.GLenum(pname), &buf[0])
	return int(buf[0])
}

func (f *Functions) GetProgrami(p Program, pname Enum) int {
	// Hope this is enough room.
	var buf [100]C.GLint
	C.glGetProgramiv(C.GLuint(p), C.GLenum(pname), &buf[0])
	return int(buf[0])
}

func (f *Functions) GetProgramInfoLog(p Program) string {
	var plen C.GLsizei
	C.glGetProgramInfoLog(C.GLuint(p), 0, &plen, nil)
	if plen == 0 {
		return ""
	}
	// Make room for the string and the null terminator.
	buf := make([]byte, plen+1)
	C.glGetProgramInfoLog(C.GLuint(p), C.GLsizei(len(buf)), &plen, (*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf[:len(buf)-1])
}

func (f *Functions) GetQueryObjectuiv(query Query, pname Enum) uint {
	// Hope this is enough room.
	var buf [100]C.GLuint
	C.gio_glGetQueryObjectuiv(C.GLuint(query), C.GLenum(pname), &buf[0])
	return uint(buf[0])
}

func (f *Functions) GetShaderi(s Shader, pname Enum) int {
	// Hope this is enough room.
	var buf [100]C.GLint
	C.glGetShaderiv(C.GLuint(s), C.GLenum(pname), &buf[0])
	return int(buf[0])
}

func (f *Functions) GetShaderInfoLog(s Shader) string {
	var plen C.GLsizei
	C.glGetShaderInfoLog(C.GLuint(s), 0, &plen, nil)
	if plen == 0 {
		return ""
	}
	// Make room for the string and the null terminator.
	buf := make([]byte, plen+1)
	C.glGetShaderInfoLog(C.GLuint(s), C.GLsizei(len(buf)), &plen, (*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf[:len(buf)-1])
}

func (f *Functions) GetString(pname Enum) string {
	str := C.glGetString(C.GLenum(pname))
	return C.GoString((*C.char)(unsafe.Pointer(str)))
}

func (f *Functions) GetUniformLocation(p Program, name string) Uniform {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Uniform(C.glGetUniformLocation(C.GLuint(p), cname))
}

func (f *Functions) InvalidateFramebuffer(target, attachment Enum) {
	C.gio_glInvalidateFramebuffer(C.GLenum(target), C.GLenum(attachment))
}

func (f *Functions) LinkProgram(p Program) {
	C.glLinkProgram(C.GLuint(p))
}

func (f *Functions) PixelStorei(pname Enum, param int32) {
	C.glPixelStorei(C.GLenum(pname), C.GLint(param))
}

func (f *Functions) Scissor(x, y, width, height int32) {
	C.glScissor(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func (f *Functions) RenderbufferStorage(target, internalformat Enum, width, height int) {
	C.glRenderbufferStorage(C.GLenum(target), C.GLenum(internalformat), C.GLsizei(width), C.GLsizei(height))
}

func (f *Functions) ShaderSource(s Shader, src string) {
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	strlen := C.GLint(len(src))
	C.glShaderSource(C.GLuint(s), 1, &csrc, &strlen)
}

func (f *Functions) TexImage2D(target Enum, level int, internalFormat int, width int, height int, format Enum, ty Enum, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glTexImage2D(C.GLenum(target), C.GLint(level), C.GLint(internalFormat), C.GLsizei(width), C.GLsizei(height), 0, C.GLenum(format), C.GLenum(ty), p)
}

func (f *Functions) TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glTexSubImage2D(C.GLenum(target), C.GLint(level), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLenum(format), C.GLenum(ty), p)
}

func (f *Functions) TexParameteri(target, pname Enum, param int) {
	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

func (f *Functions) Uniform1f(dst Uniform, v float32) {
	C.glUniform1f(C.GLint(dst), C.GLfloat(v))
}

func (f *Functions) Uniform1i(dst Uniform, v int) {
	C.glUniform1i(C.GLint(dst), C.GLint(v))
}

func (f *Functions) Uniform2f(dst Uniform, v0 float32, v1 float32) {
	C.glUniform2f(C.GLint(dst), C.GLfloat(v0), C.GLfloat(v1))
}

func (f *Functions) Uniform3f(dst Uniform, v0 float32, v1 float32, v2 float32) {
	C.glUniform3f(C.GLint(dst), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
}

func (f *Functions) Uniform4f(dst Uniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	C.glUniform4f(C.GLint(dst), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
}

func (f *Functions) UseProgram(p Program) {
	C.glUseProgram(C.GLuint(p))
}

func (f *Functions) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int) {
	var n C.GLboolean = C.GL_FALSE
	if normalized {
		n = C.GL_TRUE
	}
	C.gio_glVertexAttribPointer(C.GLuint(dst), C.GLint(size), C.GLenum(ty), n, C.GLsizei(stride), C.uintptr_t(offset))
}

func (f *Functions) Viewport(x int, y int, width int, height int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}