package zoho

import (
	"log"
	"testing"
)

// go test -run TestAuto
func TestAuto(t *testing.T) {
	err := AutoUpdateContact("477339000003181915")
	if err != nil {
		panic(err)
	}
}

// go test -run TestAuto477339000003181915
func TestAuto477339000003181915(t *testing.T) {
	err := AutoUpdateContact("477339000003181915")
	if err != nil {
		panic(err)
	}
}

// go test -run TestAutoAll
func TestAutoAll(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := AutoUpdateAllContacts()
	if err != nil {
		panic(err)
	}
}

func TestPhone1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("0630901698", "123")
	if p != "+33630901698" || o != "123" {
		panic("0630901698")
	}
}

func TestPhone2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("0730901698", "123")
	if p != "+33730901698" || o != "123" {
		panic("0730901698")
	}
}

func TestPhone3(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("0230901698", "123")
	if p != "" || o != "0230901698" {
		panic("0230901698")
	}
}

func TestPhone4(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("630901698", "123")
	if p != "+33630901698" || o != "123" {
		panic("630901698")
	}
}

func TestPhone5(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("730901698", "123")
	if p != "+33730901698" || o != "123" {
		panic("730901698")
	}
}

func TestPhone6(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("6730901698", "123")
	if p != "" || o != "6730901698" {
		panic("6730901698")
	}
}

func TestPhone7(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("+6730901698", "123")
	if p != "" || o != "+6730901698" {
		panic("+6730901698")
	}
}

func TestPhone8(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("+330630901698", "123")
	if p != "+33630901698" || o != "123" {
		panic("+330630901698")
	}
}

func TestPhone9(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("+330730901698", "123")
	if p != "+33730901698" || o != "123" {
		panic("+330730901698")
	}
}

func TestPhone10(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p, o := normalizePhone("", "+330730901698")
	if p != "+33730901698" || o != "" {
		panic("+330730901698")
	}
}
