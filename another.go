// Copyright The OpenTelemetry Authors, ashusvnath
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-logr/stdr"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	fooKey     = attribute.Key("ex.com/foo")
	barKey     = attribute.Key("ex.com/bar")
	lemonsKey  = attribute.Key("ex.com/lemons")
	anotherKey = attribute.Key("ex.com/another")
)

var tp *sdktrace.TracerProvider

// initTracer creates and registers trace provider instance.
func initTracer() error {
	exp, err :=
		otlptracehttp.New(context.Background(),
			otlptracehttp.WithInsecure())

	if err != nil {
		return fmt.Errorf("failed to initialize httptrace exporter: %w", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	// Set logging level to info to see SDK status messages
	stdr.SetVerbosity(5)

	// initialize trace provider.
	if err := initTracer(); err != nil {
		log.Panic(err)
	}

	// Create a named tracer with package path as its name.
	tracer := tp.Tracer("example/namedtracer/main")
	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()

	m0, _ := baggage.NewMemberRaw(string(fooKey), "foo1")
	m1, _ := baggage.NewMemberRaw(string(barKey), "bar1")
	b, _ := baggage.New(m0, m1)
	ctx = baggage.ContextWithBaggage(ctx, b)

	// var span trace.Span
	// ctx, span = tracer.Start(ctx, "operation")
	// defer span.End()
	for {
		// span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("time", time.Now().Nanosecond())))
		// span.SetAttributes(anotherKey.String("yes"))
		Operation(ctx, tracer)
		time.Sleep(time.Millisecond * 50)
	}
}

func Operation(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "Main operation")
	span.End()
	if err := SubOperation(ctx, tracer); err != nil {
		panic(err)
	}
	if err := PeerOperation(ctx, tracer); err != nil {
		panic(err)
	}
}

// SubOperation is an example to demonstrate the use of named tracer.
// It creates a named tracer with its package path.
func SubOperation(ctx context.Context, tr trace.Tracer) error {
	// Using global provider. Alternative is to have application provide a getter
	// for its component to get the instance of the provider.
	// tr := otel.Tracer("example/namedtracer/foo")

	var span trace.Span
	_, span = tr.Start(ctx, "Sub operation...")
	defer span.End()
	span.SetAttributes(lemonsKey.String("five"))
	span.AddEvent("Sub span event")
	time.Sleep(time.Millisecond * 15)
	SubChildOperation(ctx, tr)
	return nil
}

// SubOperation is an example to demonstrate the use of named tracer.
// It creates a named tracer with its package path.
func PeerOperation(ctx context.Context, tr trace.Tracer) error {
	// Using global provider. Alternative is to have application provide a getter
	// for its component to get the instance of the provider.
	// tr := otel.Tracer("example/namedtracer/foo")

	var span trace.Span
	_, span = tr.Start(ctx, "Peer operation...")
	defer span.End()
	span.SetAttributes(lemonsKey.String("five"))
	span.AddEvent("Peer span event")
	time.Sleep(time.Millisecond * 5)
	return nil
}

// SubOperation is an example to demonstrate the use of named tracer.
// It creates a named tracer with its package path.
func SubChildOperation(ctx context.Context, tr trace.Tracer) error {
	// Using global provider. Alternative is to have application provide a getter
	// for its component to get the instance of the provider.
	// tr := otel.Tracer("example/namedtracer/foo")

	var span trace.Span
	_, span = tr.Start(ctx, "Sub child operation...")
	defer span.End()
	span.SetAttributes(lemonsKey.String("five"))
	span.AddEvent("Sub child span event")
	time.Sleep(time.Millisecond * 10)
	return nil
}
