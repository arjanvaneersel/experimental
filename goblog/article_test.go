package goblog

import "testing"

func TestSetSlug(t *testing.T) {
	t.Log("While testing SetSlug")
	{
		a := Article{Title: "This is a test"}
		a.SetSlug()
		expected := "this-is-a-test"

		if a.Slug != expected {
			t.Fatalf("\tExpected slug to be \"%s\", but received \"%s\" instead.", expected, a.Slug)
		}
		t.Log("\tReceived expected slug")
	}
}

func TestSetSlugWithTranslation(t *testing.T) {
	t.Log("While testing SetSlug")
	{
		a := Article{
			Title: "This is a test",
			Translations: ArticleTranslations{
				"nl": &ArticleTranslation{
					Title: "Dit is een test",
				},
			},
		}
		a.SetSlug()
		expected := "dit-is-een-test"
		t.Log(a)

		if a.Translations["nl"].Slug != expected {
			t.Fatalf("\tExpected slug to be \"%s\", but received \"%s\" instead.", expected, a.Translations["nl"].Slug)
		}
		t.Log("\tReceived expected slug")
	}
}
