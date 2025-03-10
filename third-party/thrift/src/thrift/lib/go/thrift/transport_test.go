/*
 * Copyright (c) Meta Platforms, Inc. and affiliates.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package thrift

import (
	"fmt"
	"io"
	"testing"
)

const TRANSPORT_BINARY_DATA_SIZE = 4096

var (
	transport_bdata  []byte // test data for writing; same as data
	transport_header map[string]string
)

func init() {
	transport_bdata = make([]byte, TRANSPORT_BINARY_DATA_SIZE)
	for i := 0; i < TRANSPORT_BINARY_DATA_SIZE; i++ {
		transport_bdata[i] = byte((i + 'a') % 255)
	}
	transport_header = map[string]string{"key": "User-Agent",
		"value": "Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1667.0 Safari/537.36"}
}

func TransportTest(t *testing.T, writer io.Writer, reader io.Reader) {
	buf := make([]byte, TRANSPORT_BINARY_DATA_SIZE)

	// Special case for header transport -- need to reset protocol on read
	var headerTrans *headerTransport
	if hdr, ok := reader.(*headerTransport); ok {
		headerTrans = hdr
	}

	_, err := writer.Write(transport_bdata)
	if err != nil {
		t.Fatalf("Transport %T cannot write binary data of length %d: %s", writer, len(transport_bdata), err)
	}
	err = flush(writer)
	if err != nil {
		t.Fatalf("Transport %T cannot flush write of binary data: %s", writer, err)
	}

	if headerTrans != nil {
		err = headerTrans.ResetProtocol()
		if err != nil {
			t.Errorf("Header Transport %T cannot read binary data frame", reader)
		}
	}
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Errorf("Transport %T cannot read binary data of length %d: %s", reader, TRANSPORT_BINARY_DATA_SIZE, err)
	}
	if n != TRANSPORT_BINARY_DATA_SIZE {
		t.Errorf("Transport %T read only %d instead of %d bytes of binary data", reader, n, TRANSPORT_BINARY_DATA_SIZE)
	}
	for k, v := range buf {
		if v != transport_bdata[k] {
			t.Fatalf("Transport %T read %d instead of %d for index %d of binary data 2", reader, v, transport_bdata[k], k)
		}
	}
	_, err = writer.Write(transport_bdata)
	if err != nil {
		t.Fatalf("Transport %T cannot write binary data 2 of length %d: %s", writer, len(transport_bdata), err)
	}
	err = flush(writer)
	if err != nil {
		t.Fatalf("Transport %T cannot flush write binary data 2: %s", writer, err)
	}

	if headerTrans != nil {
		err = headerTrans.ResetProtocol()
		if err != nil {
			t.Errorf("Header Transport %T cannot read binary data frame 2", reader)
		}
	}
	buf = make([]byte, TRANSPORT_BINARY_DATA_SIZE)
	read := 1
	for n = 0; n < TRANSPORT_BINARY_DATA_SIZE && read != 0; {
		read, err = reader.Read(buf[n:])
		if err != nil {
			t.Errorf("Transport %T cannot read binary data 2 of total length %d from offset %d: %s", reader, TRANSPORT_BINARY_DATA_SIZE, n, err)
		}
		n += read
	}
	if n != TRANSPORT_BINARY_DATA_SIZE {
		t.Errorf("Transport %T read only %d instead of %d bytes of binary data 2", reader, n, TRANSPORT_BINARY_DATA_SIZE)
	}
	for k, v := range buf {
		if v != transport_bdata[k] {
			t.Fatalf("Transport %T read %d instead of %d for index %d of binary data 2", reader, v, transport_bdata[k], k)
		}
	}
}

func transportHTTPClientTest(t *testing.T, writer io.Writer, reader io.Reader) {
	buf := make([]byte, TRANSPORT_BINARY_DATA_SIZE)

	// Need to assert type of Transport to HTTPClient to expose the Setter
	httpWPostTrans := writer.(*httpClient)
	httpWPostTrans.SetHeader(transport_header["key"], transport_header["value"])

	_, err := writer.Write(transport_bdata)
	if err != nil {
		t.Fatalf("Transport %T cannot write binary data of length %d: %s", writer, len(transport_bdata), err)
	}
	err = flush(writer)
	if err != nil {
		t.Fatalf("Transport %T cannot flush write of binary data: %s", writer, err)
	}
	// Need to assert type of Transport to HTTPClient to expose the Getter
	httpRPostTrans := reader.(*httpClient)
	readHeader := httpRPostTrans.GetHeader(transport_header["key"])
	if err != nil {
		t.Errorf("Transport %T cannot read HTTP Header Value", httpRPostTrans)
	}

	if transport_header["value"] != readHeader {
		t.Errorf("Expected HTTP Header Value %s, got %s", transport_header["value"], readHeader)
	}
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Errorf("Transport %T cannot read binary data of length %d: %s", reader, TRANSPORT_BINARY_DATA_SIZE, err)
	}
	if n != TRANSPORT_BINARY_DATA_SIZE {
		t.Errorf("Transport %T read only %d instead of %d bytes of binary data", reader, n, TRANSPORT_BINARY_DATA_SIZE)
	}
	for k, v := range buf {
		if v != transport_bdata[k] {
			t.Fatalf("Transport %T read %d instead of %d for index %d of binary data 2", reader, v, transport_bdata[k], k)
		}
	}
}

func TestIsEOF(t *testing.T) {
	if !isEOF(io.EOF) {
		t.Fatalf("expected true")
	}
	if !isEOF(fmt.Errorf("wrapped error: %w", io.EOF)) {
		t.Fatalf("expected true")
	}
	if !isEOF(NewTransportException(END_OF_FILE, "dummy")) {
		t.Fatalf("expected true")
	}
	if !isEOF(fmt.Errorf("wrapped trasport error: %w", NewTransportException(END_OF_FILE, "dummy"))) {
		t.Fatalf("expected true")
	}
}
