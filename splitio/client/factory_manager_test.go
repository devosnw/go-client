package client

import (
	"github.com/splitio/go-client/splitio/conf"
	"github.com/splitio/go-toolkit/logging"

	"testing"
)

func TestFactoryManagerMultiple(t *testing.T) {
	sdkConf := conf.Default()
	sdkConf.Logger = logging.NewLogger(options)
	sdkConf.SplitFile = "../../testdata/splits.yaml"
	factory, _ := NewSplitFactory("localhost", sdkConf)
	client := factory.Client()

	if FactoryTrackerInstantiation["localhost"] != 1 {
		t.Error("It should be 1")
	}

	factory2, _ := NewSplitFactory("localhost", sdkConf)
	_ = factory2.Client()

	if FactoryTrackerInstantiation["localhost"] != 2 {
		t.Error("It should be 2")
	}
	expected := "Factory Instantiation: You already have 1 factory with this API Key. We recommend keeping only one " +
		"instance of the factory at all times (Singleton pattern) and reusing it throughout your application."
	if strMsg != expected {
		t.Error("Wrong logger message")
	}

	factory4, _ := NewSplitFactory("asdadd", sdkConf)
	client2 := factory4.Client()
	expected = "Factory Instantiation: You already have an instance of the Split factory. Make sure you definitely want " +
		"this additional instance. We recommend keeping only one instance of the factory at all times (Singleton pattern) and " +
		"reusing it throughout your application."
	if strMsg != expected {
		t.Error("Wrong logger message")
	}

	client.Destroy()

	if FactoryTrackerInstantiation["localhost"] != 1 {
		t.Error("It should be 1")
	}

	if FactoryTrackerInstantiation["asdadd"] != 1 {
		t.Error("It should be 1")
	}

	client.Destroy()

	if FactoryTrackerInstantiation["localhost"] != 1 {
		t.Error("It should be 1")
	}

	client2.Destroy()

	_, exist := FactoryTrackerInstantiation["asdadd"]
	if exist {
		t.Error("It should not exist")
	}

	factory3, _ := NewSplitFactory("localhost", sdkConf)
	_ = factory3.Client()
	expected = "Factory Instantiation: You already have 1 factory with this API Key. We recommend keeping only one " +
		"instance of the factory at all times (Singleton pattern) and reusing it throughout your application."
	if strMsg != expected {
		t.Error("Wrong logger message")
	}

	if FactoryTrackerInstantiation["localhost"] != 2 {
		t.Error("It should be 2")
	}

	factory5, _ := NewSplitFactory("localhost", sdkConf)
	_ = factory5.Client()
	expected = "Factory Instantiation: You already have 2 factories with this API Key. We recommend keeping only one " +
		"instance of the factory at all times (Singleton pattern) and reusing it throughout your application."
	if strMsg != expected {
		t.Error("Wrong logger message", strMsg)
	}
	if FactoryTrackerInstantiation["localhost"] != 3 {
		t.Error("It should be 3")
	}
}
