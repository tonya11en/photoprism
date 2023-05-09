package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_RawEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.NotEqual(t, c.DisableRaw(), c.RawEnabled())
}

func TestConfig_RawTherapeeBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.Contains(c.RawTherapeeBin(), "/bin/rawtherapee-cli"))
}

func TestConfig_RawTherapeeBlacklist(t *testing.T) {
	c := NewConfig(CliTestContext())

	c.options.RawTherapeeBlacklist = "foo,bar"
	assert.Equal(t, "foo,bar", c.RawTherapeeBlacklist())
	c.options.RawTherapeeBlacklist = ""
	assert.Equal(t, "", c.RawTherapeeBlacklist())
}

func TestConfig_RawTherapeeEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.RawTherapeeEnabled())

	c.options.DisableRawTherapee = true
	assert.False(t, c.RawTherapeeEnabled())
}

func TestConfig_DarktableBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.Contains(c.DarktableBin(), "/bin/darktable-cli"))
}

func TestConfig_DarktableBlacklist(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "raf,cr3", c.DarktableBlacklist())
}

func TestConfig_DarktablePresets(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.False(t, c.RawPresets())
}

func TestConfig_DarktableEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.DarktableEnabled())

	c.options.DisableDarktable = true
	assert.False(t, c.DarktableEnabled())
}

func TestConfig_SipsBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	bin := c.SipsBin()
	assert.Equal(t, "", bin)
}

func TestConfig_SipsEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.NotEqual(t, c.DisableSips(), c.SipsEnabled())
}

func TestConfig_HeifConvertBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	bin := c.HeifConvertBin()
	assert.Contains(t, bin, "/bin/heif-convert")
}

func TestConfig_HeifConvertEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.HeifConvertEnabled())

	c.options.DisableHeifConvert = true
	assert.False(t, c.HeifConvertEnabled())
}

func TestConfig_RsvgConvertBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	bin := c.RsvgConvertBin()
	assert.Contains(t, bin, "/bin/rsvg-convert")
}

func TestConfig_RsvgConvertEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.RsvgConvertEnabled())

	c.options.DisableVectors = true
	assert.False(t, c.RsvgConvertEnabled())
}
