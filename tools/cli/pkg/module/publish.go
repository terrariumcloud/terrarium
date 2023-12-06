package module

import (
	"context"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"io"
)

const DefaultChunkSize = 64 * 1024

type Metadata struct {
	Name        string
	Version     string
	Description string
	Source      string
	Maturity    module.Maturity
}

func Publish(client module.PublisherClient, source io.Reader, metadata Metadata) error {
	register := &module.RegisterModuleRequest{
		Name:        metadata.Name,
		Description: metadata.Description,
		Source:      metadata.Source,
		Maturity:    metadata.Maturity,
	}

	if _, err := client.Register(context.Background(), register); err != nil {
		return err
	}

	moduleVersion := &module.Module{
		Name:    metadata.Name,
		Version: metadata.Version,
	}

	begin := &module.BeginVersionRequest{
		Module: moduleVersion,
	}

	if _, err := client.BeginVersion(context.Background(), begin); err != nil {
		return err
	}

	end := &module.EndVersionRequest{
		Module: moduleVersion,
		Action: module.EndVersionRequest_PUBLISH,
	}

	if err := upload(client, source, moduleVersion); err != nil {
		end.Action = module.EndVersionRequest_DISCARD
		_, _ = client.EndVersion(context.Background(), end)
		return err
	}

	if _, err := client.EndVersion(context.Background(), end); err != nil {
		return err
	}

	return nil
}

func upload(client module.PublisherClient, source io.Reader, moduleVersion *module.Module) error {
	stream, err := client.UploadSourceZip(context.Background())

	if err != nil {
		return err
	}

	for {
		chunk := make([]byte, DefaultChunkSize)
		n, err := source.Read(chunk)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if n < len(chunk) {
			chunk = chunk[:n]
		}

		req := &module.UploadSourceZipRequest{
			Module:       moduleVersion,
			ZipDataChunk: chunk,
		}

		if err := stream.Send(req); err != nil {
			return err
		}
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		return err
	}
	return nil
}
