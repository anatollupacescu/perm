package example

import (
	"math/bits"
	"testing"
)

/*
There was a museum heist (or robbery). Exactly one of the three suspects — Alex, Blake, or Casey — is guilty.
Among the three, exactly one is telling the truth (truth-teller)
Their statements are:Alex: “Blake is guilty.”
Blake: “Casey is innocent.” (or sometimes “Casey did not do it.”)
Casey: “Alex is lying.” (which is equivalent to “Alex is guilty,” since if Alex is lying, and only one is guilty, it points back to guilt/innocence logic).
*/
func TestMuseumTheft(t *testing.T) {
	// exactly one guilty suspect
	// guilty suspects always lie
	// innocent suspects always tell the truth
	// who stole the painting?
	type suspect struct {
		is_guilty bool
		say       func(map[string]*suspect) bool
	}

	newWorld := func() map[string]*suspect {
		return map[string]*suspect{
			"Alex":  {say: func(m map[string]*suspect) bool { return m["Blake"].is_guilty }},
			"Blake": {say: func(m map[string]*suspect) bool { return !m["Casey"].is_guilty }},
			"Casey": {say: func(m map[string]*suspect) bool { return m["Alex"].is_guilty }}, // Equivalent to "Alex is lying"
		}
	}

	checkWorld := func(guiltyName string) (bool, map[string]*suspect) {
		world := newWorld()
		world[guiltyName].is_guilty = true

		valid := true
		for _, s := range world {
			statementIsTrue := s.say(world)

			// Rule: A liar must say something false; a truth-teller must say something true.
			// Therefore, is_guilty must NEVER equal statementIsTrue.
			if s.is_guilty == statementIsTrue {
				valid = false
				break
			}
		}

		return valid, world
	}

	var foundSolutions int
	var suspects = []string{"Alex", "Blake", "Casey"}

	for _, name := range suspects {
		valid, world := checkWorld(name)
		if valid {
			foundSolutions++
			t.Logf("Solution found for guilty set %v", name)
			for _, n := range suspects {
				t.Logf("- %s is guilty: %v", n, world[n].is_guilty)
			}
		}
	}

	if foundSolutions != 1 {
		t.Fatalf("want 1 solution, got %d", foundSolutions)
	}
}

const (
	AlexBit  = 1 << iota // 001
	BlakeBit             // 010
	CaseyBit             // 100
)

func TestMuseumTheftBitmask(t *testing.T) {
	type Suspect struct {
		Name      string
		Statement func(world int) bool
	}

	// Define suspects and their statements
	suspects := []Suspect{
		{Name: "Alex", Statement: func(w int) bool { return (w & BlakeBit) != 0 }},  // "Blake is guilty"
		{Name: "Blake", Statement: func(w int) bool { return (w & CaseyBit) == 0 }}, // "Casey is innocent"
		{Name: "Casey", Statement: func(w int) bool { return (w & AlexBit) != 0 }},  // "Alex is guilty"
	}

	numSuspects := len(suspects)
	totalWorlds := 1 << numSuspects
	var foundSolutions int

	for world := range totalWorlds {
		t.Logf("Checking world: %03b", world)

		// Print roles in this world
		for i, s := range suspects {
			role := "Innocent"
			if world&(1<<i) != 0 {
				role = "Guilty"
			}
			t.Logf(" - %s: %s", s.Name, role)
		}

		// Only consider worlds with exactly one guilty suspect
		guiltyCount := bits.OnesCount(uint(world))
		if guiltyCount != 1 {
			t.Logf("   Skipping: guilty count = %d (must be 1)", guiltyCount)
			continue
		}

		// Check statements relative to guilty/innocent
		matches := true
		for i, s := range suspects {
			isGuilty := (world & (1 << i)) != 0
			statementTrue := s.Statement(world)
			t.Logf("   Statement evaluation: %s says %d → %v", s.Name, i, statementTrue)
			// Guilty must lie, innocent must tell the truth
			if statementTrue == isGuilty {
				matches = false
				t.Logf("     Invalid: %s statement contradicts role", s.Name)
				break
			}
		}

		if matches {
			foundSolutions++
			t.Logf(" ✅ Found valid scenario! World: %03b", world)
			for i, s := range suspects {
				role := "Innocent"
				if world&(1<<i) != 0 {
					role = "Guilty"
				}
				t.Logf("   - %s is %s", s.Name, role)
			}
		}
	}

	if foundSolutions != 1 {
		t.Fatalf("Expected 1 solution, found %d", foundSolutions)
	}
}
