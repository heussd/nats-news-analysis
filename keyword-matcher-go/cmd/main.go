package main

import (
	"context"
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const name = "keyword-matcher-go"

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	//exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost("localhost")))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("environment", "production"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)
	tr := tp.Tracer("component-main")

	queue.WithArticleUrls(func(m *nats.Msg) {
		var startTime = time.Now()
		_, span := tr.Start(context.Background(), "retrieval")

		var url = string(m.Data)
		var fulltext = fulltextrss.RetrieveFullText(url)
		span.End()

		_, span = tr.Start(context.Background(), "analysis")

		var text = strings.Join([]string{
			fulltext.Title,
			fulltext.Excerpt,
			fulltext.Content,
		}, " ")

		var match, matchingText = keywords.Match(text)
		var elapsedTime = time.Since(startTime)
		
		span.End()
		if match {
			queue.PushToPocket(url, matchingText)
			fmt.Printf("✅ %s (analysis took %s)\n", url, elapsedTime)
		} else {
			fmt.Printf("❌ %s (analysis took %s)\n", url, elapsedTime)
		}
	})
	fmt.Println("\n🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
