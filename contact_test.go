package zoho

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	xls "github.com/xuri/excelize/v2"
)

func TestCreateContact(t *testing.T) {
	id, err := CreateContact(Contact{
		FirstName:  "Test",
		LastName:   "TEST",
		LeadSource: "Test",
		OpenHouse:  &Date{time.Now()},
	})
	if err != nil {
		panic(err)
	}

	err = CancelOpenHouse(id)
	if err != nil {
		panic(err)
	}

	id, err = AddContactNote(id, "This is a note", "Bla bla bla")
	if err != nil {
		panic(err)
	}

}

func TestCreateDelete(t *testing.T) {
	_, err := CreateContact(Contact{
		Email:      "test132@test.com",
		FirstName:  "Test",
		LastName:   "TEST TEST",
		LeadSource: "Test",
		OpenHouse:  &Date{time.Now()},
	})
	if err != nil {
		panic(err)
	}

	id, err := FindContact("test132@test.com")
	if err != nil {
		panic(err)
	}
	if id == "" {
		panic("not found")
	}

	_, err = GetContact(id)
	if err != nil {
		panic(err)
	}
	err = DeleteContact(id)
	if err != nil {
		panic(err)
	}
}

func TestFindContact2(t *testing.T) {
	id, err := FindContact("eb232235@gmail.com")
	if err != nil {
		panic(err)
	}
	_, err = GetContact(id)
	if err != nil {
		panic(err)
	}
}

func TestContactEmails(t *testing.T) {
	_, err := GetContactEmails("477339000004039009")
	if err != nil {
		panic(err)
	}
}

func TestContactNotes(t *testing.T) {
	c, d, err := GetContactNotesCount("477339000004039009")
	if err != nil {
		panic(err)
	}
	t.Log(c, d)
}

const EMAIL = 9

// go test -run TestImportNomad
func TestImportNomad(t *testing.T) {
	var filePath = "C:\\Sync\\10-Prospection\\Génération de leads\\Nomad Education\\ALGOSUP - leads Nomad Education - 13 mars 2023 - 16 mars 2023.xlsx"

	f, err := xls.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	cn, err := xls.CoordinatesToCellName(9, 1)
	if err != nil {
		log.Fatal(err)
	}

	c1, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}
	if c1 != "Email" {
		panic("Column I should be 'Email'.")
	}

	var r = 2
	for {
		cn, err := xls.CoordinatesToCellName(EMAIL, r)
		if err != nil {
			log.Fatal(err)
		}

		email, err := f.GetCellValue(f.GetSheetList()[0], cn)
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(email) == "" {
			fmt.Println(r)
			break
		}
		fmt.Println(email)
		id, err := FindContact(email)
		if err != nil {
			panic(err)
		}
		if id != "" {
			fmt.Println(email, "already exists.")
			r++
			continue
		}

		var contact = Contact{
			Email:          email,
			Type:           "Prospect",
			LeadSource:     "Nomad Education",
			Salutation:     Salutation(f, r),
			FirstName:      First(f, r),
			LastName:       Last(f, r),
			Phone:          Phone(f, r),
			StudyLevel:     StudyLevel(f, r),
			MailingStreet:  Street(f, r),
			MailingCity:    City(f, r),
			MailingZip:     Zip(f, r),
			MailingCountry: Country(f, r),
			StudyingFor:    Studying(f, r),
			Pathway:        Pathway(f, r),
		}

		for _, o := range Options(f, r) {
			switch o {
			case "Numérique et sciences informatiques":
				contact.NSI = true
			case "Sciences de l'Ingénieur":
				contact.SciencesDeLIngNieur = true
			case "Mathématiques":
				contact.MathMatiques = true
			case "Physique-Chimie":
				contact.PhysiqueChimie = true
			case "Sciences de la Vie et de la Terre":
				contact.Biologie = true
			case "LLCE Anglais":
				contact.Anglais = true
			case "monde contemporain":
				contact.Anglais = true
			case "Humanités, littérature et philosophie":
				contact.LittRaturePhilosophie = true
			}
		}

		id, err = CreateContact(contact)
		if err != nil {
			panic(err)
		}

		pipeline := ""
		switch contact.StudyLevel {
		case "Première":
			pipeline = "2024-2025"
		case "Terminale", "BAC", "BAC+1", "BAC+2", "BAC+3":
			pipeline = "2023-2024"
		}

		if pipeline != "" {
			_, err = CreateDeal(Deal{
				DealName:  contact.FirstName + " " + contact.LastName,
				Stage:     "Prospect",
				ContactId: id,
				Pipeline:  pipeline,
			})
			if err != nil {
				panic(err)
			}
		}

		r++

	}
}

func First(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(7, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func Last(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(8, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}
	return strings.ToUpper(s)
}

func Phone(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(10, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

var study = map[string]string{
	"2ème année": "BAC+2",
	"3ème année": "BAC+3",
	"Terminale":  "Terminale",
}

func StudyLevel(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(17, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := study[s]
	if !ok {
		panic(fmt.Sprintf("'Niveau' value '%s' not found, line %d", s, row))
	}

	return v
}

var civ = map[string]string{
	"Monsieur":            "M.",
	"Madame":              "Mme",
	"Ne pas me prononcer": "",
	"":                    "",
}

func Salutation(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(6, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := civ[s]
	if !ok {
		panic(fmt.Sprintf("'Civilité' value '%s' not found, line %d", s, row))
	}

	return v
}

func Street(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(11, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func City(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(12, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func Zip(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(14, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func Country(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(15, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

var studying = map[string]string{
	"Bac général et technologique": "Bac",
	"Prépa":                        "Prépa",
	"Licence / L. AS":              "Licence",
}

func Studying(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(16, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := studying[s]
	if !ok {
		panic(fmt.Sprintf("'Diplôme en cours' value '%s' not found, line %d", s, row))
	}

	return v
}

var pathway = map[string]string{
	"Bac Général": "Général",
	"STI2D":       "STI2D",
	"STMG":        "STMG",
	"LEA":         "LEA",
	"MP":          "MP",
	"MPI":         "MPI",
	"PT":          "PT",
	"PSI":         "PSI",
	"PC":          "PC",
	"ATS":         "ATS",
	"BCPST":       "BCPST",
	"TSI":         "TSI",
	"TB":          "TB",
}

func Pathway(f *xls.File, row int) string {
	cn, err := xls.CoordinatesToCellName(19, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := pathway[s]
	if !ok {
		panic(fmt.Sprintf("'Filière' value '%s' not found, line %d", s, row))
	}

	return v
}

func Options(f *xls.File, row int) []string {
	cn, err := xls.CoordinatesToCellName(22, row)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.GetCellValue(f.GetSheetList()[0], cn)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(s, ", ")
}
