package fontloader

import (
	"os"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

type FontLoaderSuite struct {
}

func Test(t *testing.T) {
	_ = Suite(&FontLoaderSuite{})
	TestingT(t)
}

func (s *FontLoaderSuite) TestLoadSans(c *C) {
	f, err := Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f, NotNil)
}

func (s *FontLoaderSuite) TestLoadSansDefaultCache(c *C) {
	f, err := LoadCache("sans")
	c.Assert(err, IsNil)
	c.Assert(f, NotNil)
}

func (s *FontLoaderSuite) TestLoadPath(c *C) {
	path, err := findFont("sans")
	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(path, "/"), Equals, true)

	_, err = os.Stat(path)
	c.Assert(err, IsNil)

	f, err := Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f, NotNil)
}

func (s *FontLoaderSuite) TestLoadBadFontFile(c *C) {
	f, err := Load("/proc/cpuinfo") // Definitely not a TTF file
	c.Assert(err, NotNil)
	c.Assert(f, IsNil)
}

func (s *FontLoaderSuite) TestNonExistentPath(c *C) {
	f, err := Load("/dev/this-path-should-not-exist")
	c.Assert(err, NotNil)
	c.Assert(f, IsNil)
}

func (s *FontLoaderSuite) TestCacheLoadSans(c *C) {
	f, err := Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f, NotNil)
}

func (s *FontLoaderSuite) TestCacheLoadPath(c *C) {
	cache := NewFontCache()

	path, err := findFont("sans")
	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(path, "/"), Equals, true)

	_, err = os.Stat(path)
	c.Assert(err, IsNil)

	f, err := cache.Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f, NotNil)
}

func (s *FontLoaderSuite) TestCacheLoadBadFontFile(c *C) {
	cache := NewFontCache()
	f, err := cache.Load("/proc/cpuinfo") // Definitely not a TTF file
	c.Assert(err, NotNil)
	c.Assert(f, IsNil)
}

func (s *FontLoaderSuite) TestCacheNonExistentPath(c *C) {
	cache := NewFontCache()
	f, err := cache.Load("/dev/this-path-should-not-exist")
	c.Assert(err, NotNil)
	c.Assert(f, IsNil)
}

func (s *FontLoaderSuite) TestNameCacheHit(c *C) {
	cache := NewFontCache()

	f1, err := cache.Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f1, NotNil)

	f2, err := cache.Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f2, NotNil)

	c.Assert(f1 == f2, Equals, true)
}

func (s *FontLoaderSuite) TestPathCacheHit(c *C) {
	cache := NewFontCache()

	path, err := findFont("sans")
	c.Assert(err, IsNil)

	f1, err := cache.Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f1, NotNil)

	f2, err := cache.Load(path)
	c.Assert(err, IsNil)
	c.Assert(f2, NotNil)

	c.Assert(f1 == f2, Equals, true)
}

func (s *FontLoaderSuite) TestCacheMiss(c *C) {
	cache := NewFontCache()

	f1, err := cache.Load("sans")
	c.Assert(err, IsNil)
	c.Assert(f1, NotNil)

	f2, err := cache.Load("serif")
	c.Assert(err, IsNil)
	c.Assert(f2, NotNil)

	c.Assert(f1 == f2, Equals, false)
}
