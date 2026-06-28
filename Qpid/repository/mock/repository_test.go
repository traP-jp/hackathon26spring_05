package mock

import "testing"

func TestListSuggestionsExcludesActionedUsers(t *testing.T) {
	repo := NewMockRepository()
	if err := repo.LikeUser("current-user", "alice"); err != nil {
		t.Fatalf("LikeUser returned error: %v", err)
	}
	if err := repo.NopeUser("current-user", "bob"); err != nil {
		t.Fatalf("NopeUser returned error: %v", err)
	}

	suggestions, err := repo.ListSuggestions("current-user", 20)
	if err != nil {
		t.Fatalf("ListSuggestions returned error: %v", err)
	}

	for _, suggestion := range suggestions {
		if suggestion.Username == "alice" || suggestion.Username == "bob" {
			t.Fatalf("got actioned user in suggestions: %v", suggestions)
		}
	}
}
