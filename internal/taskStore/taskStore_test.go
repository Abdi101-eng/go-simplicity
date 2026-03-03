package taskstore

import (
	"abdi/task-manager/internal/models"
	"testing"
)

func TestAdd(t *testing.T) {
	var storeTests = map[string]struct {
		title            string
		priority         models.Priority
		expectedTitle    string
		expectedPriority models.Priority
	}{
		"tests success case for adding a task": {title: "FirstTask",
			priority:         models.High,
			expectedTitle:    "FirstTask",
			expectedPriority: models.High,
		},
		"success low priority": {
			title: "Second Task", priority: models.Low,
			expectedTitle: "Second Task", expectedPriority: models.Low,
		},
	}

	for name, element := range storeTests {
		t.Run(name, func(t *testing.T) {
			store := Store{}
			task, _ := store.Add(element.title, element.priority)
			if got, want := task.Title, element.expectedTitle; got != want {
				t.Errorf("title: got %q, want %q", task.Title, element.expectedTitle)
			}
			if got, want := task.Priority, element.expectedPriority; got != want {
				t.Errorf("priority: got %q, want %q", task.Priority, element.expectedPriority)
			}
			if task.Id < 1 {
				t.Errorf("Id: got %d, want > 0", task.Id)
			}
			if task.Done != false {
				t.Errorf("Done: got %v, want false", task.Done)
			}
		})
	}
}

func TestComplete(t *testing.T) {
	var store = map[string]struct {
		id            int
		setupFunc     func() Store
		expectedError bool
		expectDone    bool
	}{
		"tests success case": {
			id: 1,
			setupFunc: func() Store {
				store := NewStore()
				store.Add("First Task", models.High)
				return *store
			},
			expectedError: false,
			expectDone:    true,
		},
		"test failure case": {
			id: 99,
			setupFunc: func() Store {
				store := NewStore()
				store.Add("First Task", models.High)
				return *store
			},
			expectedError: true,
			expectDone:    false,
		},
	}

	for name, tc := range store {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			store := tc.setupFunc()
			err := store.Complete(tc.id)
			if tc.expectedError && err == nil {
				t.Errorf("expected error, got %v", err)
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error, got %v", err)
			}
			if tc.expectDone {
				if got, want := store.list[0].Done, true; got != want {
					t.Errorf("got Done %v, want %v", got, want)
				}
			}

		})
	}

}

func TestDelete(t *testing.T) {
	var testStruct = map[string]struct {
		id            int
		setupFunc     func() Store
		expectedError bool
		expectDelete  bool
	}{
		"Success: Successfully deleted": {
			id: 1,
			setupFunc: func() Store {
				store := NewStore()
				store.Add("test", models.High)
				return *store
			},
			expectedError: false,
			expectDelete:  true,
		},
		"error: non existent id": {
			id: 99,
			setupFunc: func() Store {
				store := NewStore()
				store.Add("test", models.High)
				return *store
			},
			expectedError: true,
			expectDelete:  false,
		},
	}

	for name, tc := range testStruct {
		t.Run(name, func(t *testing.T) {
			store := tc.setupFunc()
			err := store.Delete(tc.id)

			if tc.expectedError && err == nil {
				t.Errorf("expected error, got %v", err)
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error, got %v", err)
			}

			if tc.expectDelete {
				for i := 0; i < len(store.list); i++ {
					if store.list[i].Id == tc.id {
						t.Errorf("deleted ID still present in the list")
					}
				}
			}
		})

	}

}
