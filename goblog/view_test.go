package goblog

import "testing"

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TestAddPath(t *testing.T) {
	t.Log("While testing addPath")
	{
		DefaultTemplates = []string{"base.html", "nav.html"}
		TemplateDirectory = "templates/"

		expectations := []string{"templates/base.html", "templates/nav.html", "templates/test1.html", "templates/test2.html"}
		result := addPath("test1.html", "test2.html")

		for _, expected := range expectations {
			if !stringInSlice(expected, result) {
				t.Fatalf("\tResult doesn't contain %s. %v", expected, result)
			}
		}
		t.Log("\tReceived expected result")
	}
}

func TestNewView(t *testing.T) {
	t.Log("While testing NewView")
	{
		DefaultTemplates = []string{"base.html", "nav.html"}
		TemplateDirectory = "templates/"
	}

	_, err := NewView("home", "home.html")
	if err != nil {
		t.Fatalf("Expected to receive a View, but got error: %s", err)
	}
}
