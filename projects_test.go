package gerrit_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestProjectsService_ListProjects(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, testValues{
			"r": "(arch|benchmarks)",
			"n": "2",
		})

		fmt.Fprint(w, `)]}'`+"\n"+`{"arch":{"id":"arch","state":"ACTIVE"},"benchmarks":{"id":"benchmarks","state":"ACTIVE"}}`)
	})

	opt := &gerrit.ProjectOptions{
		Regex: "(arch|benchmarks)",
	}
	opt.Limit = 2
	project, _, err := testClient.Projects.ListProjects(context.Background(), opt)
	if err != nil {
		t.Errorf("Projects.ListProjects returned error: %v", err)
	}

	want := &map[string]gerrit.ProjectInfo{
		"arch": {
			ID:    "arch",
			State: "ACTIVE",
		},
		"benchmarks": {
			ID:    "benchmarks",
			State: "ACTIVE",
		},
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.ListProjects returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/go/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"go","name":"go","parent":"All-Projects","description":"The Go Programming Language","state":"ACTIVE"}`)
	})

	project, _, err := testClient.Projects.GetProject(context.Background(), "go")
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "go",
		Name:        "go",
		Parent:      "All-Projects",
		Description: "The Go Programming Language",
		State:       "ACTIVE",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject_WithSlash(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/plugins/delete-project", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/projects/plugins%2Fdelete-project")

		c := `)]}'` + "\n" + `{"id":"plugins%2Fdelete-project","name":"plugins/delete-project","parent":"Public-Plugins","description":"A plugin which allows projects to be deleted from Gerrit via an SSH command","state":"ACTIVE"}`
		fmt.Fprint(w, c)
	})
	project, _, err := testClient.Projects.GetProject(context.Background(), "plugins/delete-project")
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "plugins%2Fdelete-project",
		Name:        "plugins/delete-project",
		Parent:      "Public-Plugins",
		Description: "A plugin which allows projects to be deleted from Gerrit via an SSH command",
		State:       "ACTIVE",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

// +func (s *ProjectsService) CreateProject(name string, input *ProjectInput) (*ProjectInfo, *Response, error) {
func TestProjectsService_CreateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &gerrit.ProjectInput{
		Description: "The Go Programming Language",
	}

	testMux.HandleFunc("/projects/go/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		v := new(gerrit.ProjectInput)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"go","name":"go","parent":"All-Projects","description":"The Go Programming Language"}`)
	})

	project, _, err := testClient.Projects.CreateProject(context.Background(), "go", input)
	if err != nil {
		t.Errorf("Projects.CreateProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "go",
		Name:        "go",
		Parent:      "All-Projects",
		Description: "The Go Programming Language",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.CreateProject returned %+v, want %+v", project, want)
	}
}

// +func (s *ProjectsService) GetProjectDescription(name string) (*string, *Response, error) {
func TestProjectsService_GetProjectDescription(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/go/description/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `)]}'`+"\n"+`"The Go Programming Language"`)
	})

	description, _, err := testClient.Projects.GetProjectDescription(context.Background(), "go")
	if err != nil {
		t.Errorf("Projects.GetProjectDescription returned error: %v", err)
	}

	want := "The Go Programming Language"

	if !reflect.DeepEqual(description, want) {
		t.Errorf("Projects.GetProjectDescription returned %+v, want %+v", description, want)
	}
}

func ExampleProjectsService_ListProjects() {
	instance := "https://chromium-review.googlesource.com/"
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.ProjectOptions{
		Description: true,
		Prefix:      "infra/infra/infra_l",
	}
	projects, _, err := client.Projects.ListProjects(ctx, opt)
	if err != nil {
		panic(err)
	}

	for name, p := range *projects {
		fmt.Printf("%s - State: %s\n", name, p.State)
	}

	// Output:
	// infra/infra/infra_libs - State: ACTIVE
}

func TestProjectsService_GetBranch(t *testing.T) {
	setup()
	defer teardown()

	existBranches := map[string]*gerrit.BranchInfo{
		"branch": {
			Ref:       "123",
			Revision:  "abcd1234",
			CanDelete: true,
		},
		"branch/foo": {
			Ref:       "456",
			Revision:  "deadbeef",
			CanDelete: false,
		},
	}

	testMux.HandleFunc("/projects/go/branches/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		branchName := r.URL.Path[len("/projects/go/branches/"):]

		branchInfo, ok := existBranches[branchName]
		if !ok {
			http.Error(w, branchName, http.StatusBadRequest)
		}

		branchInfoRaw, err := json.Marshal(&branchInfo)
		if err != nil {
			http.Error(w, branchName, http.StatusBadRequest)
		}

		fmt.Fprint(w, `)]}'`+"\n"+string(branchInfoRaw))
	})

	var tests = []struct {
		name     string
		branch   string
		expected *gerrit.BranchInfo
	}{
		{
			name:   "branch without slash",
			branch: "branch",
			expected: &gerrit.BranchInfo{
				Ref:       "123",
				Revision:  "abcd1234",
				CanDelete: true,
			},
		},
		{
			name:   "branch with slash",
			branch: "branch/foo",
			expected: &gerrit.BranchInfo{
				Ref:       "456",
				Revision:  "deadbeef",
				CanDelete: false,
			},
		},
	}

	for _, tc := range tests {
		branchInfo, _, err := testClient.Projects.GetBranch(context.Background(), "go", tc.branch)
		if err != nil {
			t.Errorf("tc %s: Projects.GetProject returned error: %v", tc.name, err)
		}

		if !reflect.DeepEqual(branchInfo, tc.expected) {
			t.Errorf("tc %s: Projects.GetBranch returned %+v, want %+v", tc.name, branchInfo, tc.expected)
		}
	}
}

func TestProjectsService_SetConfig(t *testing.T) {
	setup()
	defer teardown()

	var configInput gerrit.ConfigInput

	// test server url
	instance := "https://gerrit-test.uniontech.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	// login username and password
	client.Authentication.SetBasicAuth("", "")

	configInput.Description = "test change description"

	_, _, err = client.Projects.SetConfig("my-test", &configInput)
	if err != nil {
		t.Errorf("Projects.SetConfig returned error: %v", err)
	}

}
