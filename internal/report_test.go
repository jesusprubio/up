package internal

import (
	"testing"
)

func TestReportString(t *testing.T) {
	r := Report{
		ProtocolID: "tcp",
		RHost:      "127.0.0.1:80",
		Time:       1,
		Extra:      "extra-0",
	}
	t.Run("returns a report using human format", func(t *testing.T) {
		got, err := r.String(HumanFormat)
		if err != nil {
			t.Fatal(err)
		}
		want := "✔ tcp             1ns            127.0.0.1:80 (extra-0)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("returns a report using JSON format", func(t *testing.T) {
		got, err := r.String(JSONFormat)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"protocol":"tcp","rhost":"127.0.0.1:80","time":1,"extra":"extra-0"}`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("returns a report using a grepable format", func(t *testing.T) {
		got, err := r.String(GrepFormat)
		if err != nil {
			t.Fatal(err)
		}
		want := "tcp\t1ns\t127.0.0.1:80\tok\textra-0"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestStringJSON(t *testing.T) {
	r := Report{
		ProtocolID: "tcp",
		RHost:      "127.0.0.1:80",
		Time:       1,
		Extra:      "extra-0",
	}
	t.Run("returns JSON format for successful probes", func(t *testing.T) {
		got, err := r.stringJSON()
		if err != nil {
			t.Fatal(err)
		}
		want := `{"protocol":"tcp","rhost":"127.0.0.1:80","time":1,"extra":"extra-0"}`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("returns JSON format for failed probes", func(t *testing.T) {
		rErr := r
		rErr.Extra = ""
		rErr.Error = "error-0"
		got, err := rErr.stringJSON()
		if err != nil {
			t.Fatal(err)
		}
		want := `{"protocol":"tcp","rhost":"127.0.0.1:80","time":1,"error":"error-0"}`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestStringHuman(t *testing.T) {
	r := Report{
		ProtocolID: "tcp",
		RHost:      "127.0.0.1:80",
		Time:       1,
		Extra:      "extra-0",
	}
	t.Run("returns human readable format for successful probes",
		func(t *testing.T) {
			got := r.stringHuman()
			want := "✔ tcp             1ns            127.0.0.1:80 (extra-0)"
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		},
	)
	t.Run("returns human readable format for failed probes",
		func(t *testing.T) {
			rErr := r
			rErr.Extra = ""
			rErr.Error = "error-0"
			got := rErr.stringHuman()
			want := "✘ tcp             1ns            127.0.0.1:80 (error-0)"
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		},
	)
}

func TestStringGrep(t *testing.T) {
	r := Report{
		ProtocolID: "tcp",
		RHost:      "127.0.0.1:80",
		Time:       1,
		Extra:      "extra-0",
	}
	t.Run("returns grep format for successful probes", func(t *testing.T) {
		got := r.stringGrep()
		want := "tcp\t1ns\t127.0.0.1:80\tok\textra-0"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("returns grep format for failed probes", func(t *testing.T) {
		rErr := r
		rErr.Extra = ""
		rErr.Error = "error-0"
		got := rErr.stringGrep()
		want := "tcp\t1ns\t127.0.0.1:80\terror\terror-0"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
