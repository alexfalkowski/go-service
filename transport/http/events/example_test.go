package events_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	transportevents "github.com/alexfalkowski/go-service/v2/transport/http/events"
)

func ExampleReceiver_Register() {
	mux := http.NewServeMux()
	receiver := transportevents.NewReceiver(mux, nil)
	receiver.Register(context.Background(), "/events", func(_ context.Context, event events.Event) events.Result {
		fmt.Println(event.Type())
		return nil
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	sender := transportevents.NewSender(nil)
	event := events.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")

	result := sender.Send(events.ContextWithTarget(context.Background(), server.URL+"/events"), event)
	fmt.Println(events.IsACK(result))
	// Output:
	// example.type
	// true
}
