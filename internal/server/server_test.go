package server


import(
	"testing"
	"net/http"
	"net/http/httptest"
    "strings"
	 "github.com/Anshit-Gupta/kairo/internal/store"
)


func TestPutAndGet(t *testing.T){
	s,err := store.NewStore(t.TempDir() + "/test.log")
	 if err != nil {
        t.Fatal(err)
    }
    defer s.Close()
     
	srv := NewServer(s)
	req := httptest.NewRequest(
		http.MethodPut,
		"/name",
		strings.NewReader("Ansh"),
	)

	recorder := httptest.NewRecorder()

	srv.ServeHTTP(recorder,req)

	if recorder.Code!=http.StatusOK{
		t.Fatalf("expected 200 , got %d",recorder.Code)
	}


	req = httptest.NewRequest(
		http.MethodGet,
		"/name",
		nil,
	)
    recorder = httptest.NewRecorder()
	srv.ServeHTTP(recorder,req)

	if recorder.Code!=http.StatusOK{
		t.Fatalf("expected 200 , got %d",recorder.Code)
	}
    
	if recorder.Body.String() != "Ansh\n" {
    t.Fatalf("expected Ansh, got %q", recorder.Body.String())
    }
}

func TestDelete(t *testing.T){
	s,err := store.NewStore(t.TempDir() + "/test.log")
	 if err != nil {
        t.Fatal(err)
    }
    defer s.Close()
     
	srv := NewServer(s)
	req := httptest.NewRequest(
		http.MethodPut,
		"/name",
		strings.NewReader("Ansh"),
	)

	recorder := httptest.NewRecorder()

	srv.ServeHTTP(recorder,req)

	if recorder.Code!=http.StatusOK{
		t.Fatalf("expected 200 , got %d",recorder.Code)
	}

	req = httptest.NewRequest(
		http.MethodDelete,
		"/name",
		nil,
	)

	recorder = httptest.NewRecorder()

	srv.ServeHTTP(recorder,req)

	if recorder.Code !=http.StatusNoContent{
		t.Fatalf("expected 204 , got %d",recorder.Code)
	}

	

	req = httptest.NewRequest(
		http.MethodGet,
		"/name",
		nil,
	)

	recorder = httptest.NewRecorder()

	srv.ServeHTTP(recorder,req)

	if recorder.Code!=http.StatusNotFound{
		t.Fatalf("expected 404 , got %d",recorder.Code)
	}




}


func TestGetMissingKey(t *testing.T) {
    s, err := store.NewStore(t.TempDir() + "/test.log")
    if err != nil {
        t.Fatal(err)
    }
    defer s.Close()

    srv := NewServer(s)

    req := httptest.NewRequest(
        http.MethodGet,
        "/does-not-exist",
        nil,
    )

    recorder := httptest.NewRecorder()

    srv.ServeHTTP(recorder, req)

    if recorder.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", recorder.Code)
    }
}