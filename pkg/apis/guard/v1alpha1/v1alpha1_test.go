package v1alpha1

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Utility function to convert []interface{} to [][]interface{}
func toSlice(slices []interface{}) [][]interface{} {
	result := make([][]interface{}, len(slices))
	for sliceIndex := 0; sliceIndex < len(slices); sliceIndex++ {
		s := reflect.ValueOf(slices[sliceIndex])
		if s.Kind() != reflect.Slice {
			panic("InterfaceSlice() given a non-slice type")
		}

		// Keep the distinction between nil and empty slice input
		if s.IsNil() {
			return nil
		}

		converted := make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			converted[i] = s.Index(i).Interface()
		}
		result[sliceIndex] = converted
	}
	return result
}

func ValueTests_Test(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	args := toSlice(arguments)

	t.Run("Basics", func(t *testing.T) {

		// Initial Tests
		piles[9].Clear()
		piles[3].mergeI(piles[4])
		configs[3].fuseI(configs[4])
		configs[5].learnI(piles[5])
		configs[6].decideI(profiles[0])
		piles[6].Clear()

		// Test ProfileValue
		for i, v := range args {
			profiles[i].profileI(v...)
		}

		// Test PileValue
		for i, profile := range profiles {
			piles[i].addI(profile)
			piles[i].Clear()
			piles[i].addI(profile)
			piles[i].addI(profile)
			piles[0].mergeI(piles[i])
			piles[0].mergeI(piles[i])
		}

		// Test ConfigValue
		for i, pile := range piles {
			configs[i].learnI(pile)
			configs[0].fuseI(configs[i])
			configs[0].fuseI(configs[i])
			if str := configs[0].decideI(profiles[i]); str != "" {
				t.Errorf("config.Decide(profile) wrong decission: %s\nFor profile %s\nwhen using config %s\n", str, profiles[i], configs[0])
			}
		}
	})
}

func ValueTests_Test_WithMarshal(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	profile := profiles[0]
	pile := piles[0]
	config := configs[0]
	args := toSlice(arguments)

	t.Run("Basics", func(t *testing.T) {

		// Test ProfileValue
		profile.profileI(args[0]...)

		// Test PileValue
		pile.addI(profile)
		pile.Clear()
		pile.addI(profile)
		pile.addI(profile)
		pile.mergeI(pile)
		pile.mergeI(pile)
		var err error
		var bytes []byte
		if bytes, err = json.Marshal(pile); err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		if err = json.Unmarshal(bytes, &pile); err != nil {
			t.Errorf("json.Unmarshal Error %v", err.Error())
			t.Errorf("bytes: %s", string(bytes))
		}
		// Test ConfigValue
		config.learnI(pile)
		config.fuseI(config)
		config.fuseI(config)

		if str := config.decideI(profile); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}

		if bytes, err = json.Marshal(config); err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		if err = json.Unmarshal(bytes, &config); err != nil {
			t.Errorf("json.Unmarshal Error %v", err.Error())
			t.Errorf("bytes: %s", string(bytes))
		}
	})
}

func ValueTests_SimpleTest(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	pile := piles[0]
	config := configs[0]
	args := toSlice(arguments)
	t.Run("Basics", func(t *testing.T) {
		// Test ProfileValue
		for i, v := range args {
			profiles[i].profileI(v...)
		}

		// Test PileValue
		pile.addI(profiles[0])

		// test ConfigValue
		config.learnI(pile)
		if str := config.decideI(profiles[0]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
		if str := config.decideI(profiles[1]); str == "" {
			t.Errorf("config.Decide(profile) expected a reject of %s after learning %s\n", args[1], args[0])
		}
		if str := config.decideI(profiles[2]); str == "" {
			t.Errorf("config.Decide(profile) expected a reject of %s after learning %s\n", args[2], args[0])
		}
	})
}

func ValueProfile_MarshalTest(t *testing.T, profiles []ValueProfile) {
	profileIn := profiles[0]
	profileOut := profiles[1]
	t.Run("Pile Marshal", func(t *testing.T) {

		var bytes []byte
		var err error
		if bytes, err = json.Marshal(profileIn); err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		if err = json.Unmarshal(bytes, &profileOut); err != nil {
			t.Errorf("json.Unmarshal Error %v", err.Error())
			t.Errorf("bytes: %s", string(bytes))
		}
	})
}

func ValuePile_MarshalTest(t *testing.T, piles []ValuePile) {
	pileIn := piles[0]
	pileOut := piles[1]
	t.Run("Pile Marshal", func(t *testing.T) {

		var bytes []byte
		var err error
		if bytes, err = json.Marshal(pileIn); err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		if err = json.Unmarshal(bytes, &pileOut); err != nil {
			t.Errorf("json.Unmarshal Error %v", err.Error())
			t.Errorf("bytes: %s", string(bytes))
		}
	})
}

func ValueConfig_MarshalTest(t *testing.T, configs []ValueConfig) {
	configIn := configs[0]
	configOut := configs[1]

	t.Run("Config Marshal", func(t *testing.T) {
		var bytes []byte
		var err error
		if bytes, err = json.Marshal(configIn); err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		if err = json.Unmarshal(bytes, &configOut); err != nil {
			t.Errorf("json.Unmarshal Error %v", err.Error())
			t.Errorf("bytes: %s", string(bytes))
		}
	})
}

func ValueTests_TestAdd(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	args := toSlice(arguments)
	t.Run("Basics", func(t *testing.T) {
		// Test ProfileValue
		profiles[0].profileI(args[0]...)
		profiles[1].profileI(args[1]...)

		// Test PileValue
		piles[0].addI(profiles[0])
		piles[0].addI(profiles[1])

		// test ConfigValue
		configs[0].learnI(piles[0])
		if str := configs[0].decideI(profiles[0]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
		if str := configs[0].decideI(profiles[1]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
	})
}

func ValueTests_TestMerge(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	args := toSlice(arguments)
	t.Run("Basics", func(t *testing.T) {

		// Test ProfileValue
		profiles[0].profileI(args[0]...)
		profiles[1].profileI(args[1]...)

		// Test PileValue
		piles[0].addI(profiles[0])
		piles[1].addI(profiles[1])
		piles[0].mergeI(piles[1])

		// test ConfigValue
		configs[0].learnI(piles[0])
		if str := configs[0].decideI(profiles[0]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
		if str := configs[0].decideI(profiles[1]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
	})
}

func ValueTests_TestFuse(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, arguments ...interface{}) {
	args := toSlice(arguments)
	t.Run("Basics", func(t *testing.T) {

		// Test ProfileValue
		profiles[0].profileI(args[0]...)
		profiles[1].profileI(args[1]...)

		// Test PileValue
		piles[0].addI(profiles[0])
		piles[1].addI(profiles[1])

		// test ConfigValue
		configs[0].learnI(piles[0])
		configs[1].learnI(piles[1])
		configs[0].fuseI(configs[1])
		if str := configs[0].decideI(profiles[0]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
		if str := configs[0].decideI(profiles[1]); str != "" {
			t.Errorf("config.Decide(profile) wrong decission: %s", str)
		}
	})
}

func ValueTests_All(t *testing.T, profiles []ValueProfile, piles []ValuePile, configs []ValueConfig, args ...interface{}) {
	ValueTests_SimpleTest(t, profiles, piles, configs, args...)
	ValueTests_Test(t, profiles, piles, configs, args...)
	ValueTests_TestAdd(t, profiles, piles, configs, args...)
	ValueTests_TestMerge(t, profiles, piles, configs, args...)
	ValueTests_TestFuse(t, profiles, piles, configs, args...)
	ValueConfig_MarshalTest(t, configs)
	ValuePile_MarshalTest(t, piles)
	ValueProfile_MarshalTest(t, profiles)
	ValueTests_Test_WithMarshal(t, profiles, piles, configs, args...)
}
