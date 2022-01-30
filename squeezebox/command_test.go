package squeezebox

import "testing"

func TestGetTagValueFromResponseLineSimple(t *testing.T) {

	response := "color%3Ablack number%3A141 with_colon%3A%3A%3A%3A with_spaces%3AThis%20is%20with%20spaces surname%3AZ%C3%B6rner"

	tables := []struct {
		tag string
		value string
	}{
		{"color", "black"},
		{"number", "141"},
		{"with_colon", ":::"},
		{"with_spaces", "This is with spaces"},
		{"surname", "Zörner"},
		{"missing", ""},

	}

	for _, table := range tables {
		value, _ := getTagValueFromResponseLine(response, table.tag)
		if value != table.value {
			t.Errorf("Tag %s has wrong value: %s (expected: %s)", table.tag, value, table.value)
		}
	}
}


func TestGetTagValueFromResponseLineLMS(t *testing.T) {

    response := "00%3A04%3A20%3A07%3Af8%3Aeb status - 1 tags%3AK%2Cc player_name%3ASpitzboden player_connected%3A1 player_ip%3A192.168.178.97%3A28005 power%3A1 signalstrength%3A0 mode%3Aplay time%3A132.815542135239 rate%3A1 duration%3A447.76 can_seek%3A1 mixer%20volume%3A60 playlist%20repeat%3A2 playlist%20shuffle%3A0 playlist%20mode%3Aoff seq_no%3A0 playlist_cur_index%3A0 playlist_timestamp%3A1643452585.19849 playlist_tracks%3A1 playlist%20index%3A0 id%3A48388 title%3ADarkness%20(11%2F11) coverid%3Aa2453948"

	tables := []struct {
		tag string
		value string
	}{
		{"player_name", "Spitzboden"},
		{"coverid", "a2453948"},
		{"player_ip", "192.168.178.97:28005"},
	}

	for _, table := range tables {
		value, _ := getTagValueFromResponseLine(response, table.tag)
		if value != table.value {
			t.Errorf("Tag %s has wrong value: %s (expected: %s)", table.tag, value, table.value)
		}
	}
}

func TestGetTagValueFromResponseLineSpotify(t *testing.T) {

	response := "00%3A04%3A20%3A07%3Af8%3Aeb status - 1 tags%3AK%2Cc player_name%3ASpitzboden player_connected%3A1 player_ip%3A192.168.178.97%3A28031 power%3A1 signalstrength%3A0 mode%3Aplay remote%3A1 current_title%3A1.%20Camera's%20Rolling%20von%20Agnes%20Obel%20aus%20Myopia time%3A139.330393112183 rate%3A1 duration%3A283.773 mixer%20volume%3A30 playlist%20repeat%3A2 playlist%20shuffle%3A0 playlist%20mode%3Aoff seq_no%3A0 playlist_cur_index%3A0 playlist_timestamp%3A1643470359.94226 playlist_tracks%3A10 remoteMeta%3AHASH(0x275830c0) playlist%20index%3A0 id%3A-659974472 title%3ACamera's%20Rolling artwork_url%3Ahttps%3A%2F%2Fi.scdn.co%2Fimage%2Fab67616d0000b2737da3504b2ba9f1ecb43e0ac2 coverid%3A-659974472\n"

	tables := []struct {
		tag string
		value string
	}{
		{"current_title", "1. Camera's Rolling von Agnes Obel aus Myopia"},
		{"artwork_url", "https://i.scdn.co/image/ab67616d0000b2737da3504b2ba9f1ecb43e0ac2"},
	}

	for _, table := range tables {
		value, _ := getTagValueFromResponseLine(response, table.tag)
		if value != table.value {
			t.Errorf("Tag %s has wrong value: %s (expected: %s)", table.tag, value, table.value)
		}
	}
}
