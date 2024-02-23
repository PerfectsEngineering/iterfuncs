package main

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"

	"cloud.google.com/go/bigquery"
	"github.com/perfectsengineering/iterfuncs"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	projectId = "bigquery-project"
)

func main() {
	ctx := context.Background()

	bqContainer, cleanupContainer, err := startBigQueryContainer(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	defer cleanupContainer()

	client, err := bigquery.NewClient(ctx, projectId, []option.ClientOption{
		option.WithEndpoint(bqContainer.URI),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		internaloption.SkipDialSettingsValidation(),
	}...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create BigQuery client")
	}
	defer client.Close()

	results, err := fetchDataFromBigQueryWithRange(ctx, client)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	fmt.Println(results)
}

type Row struct {
	Id   int
	Name string
}

func fetchDataFromBigQuery(ctx context.Context, client *bigquery.Client) ([]string, error) {

	// Construct a query to select data from a table
	query := client.Query("SELECT * FROM dataset1.table_a")

	// Run the query and get the results
	iter, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var results []string

	// Iterate over the results and append them to the array
	for {
		var row Row
		err := iter.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve row: %v", err)
		}

		result := fmt.Sprintf("Id: %d, Name: %s", row.Id, row.Name)
		results = append(results, result)
	}

	return results, nil
}

func fetchDataFromBigQueryWithRange(ctx context.Context, client *bigquery.Client) ([]string, error) {

	// Construct a query to select data from a table
	query := client.Query("SELECT * FROM dataset1.table_a")

	// Run the query and get the results
	iter, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var results []string

	// Iterate over the results and append them to the array
	for row, err := range iterfuncs.Range[Row](iter) {
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve row: %v", err)
		}

		result := fmt.Sprintf("Id: %d, Name: %s", row.Id, row.Name)
		results = append(results, result)
	}

	return results, nil
}

func startBigQueryContainer(ctx context.Context) (*gcloud.GCloudContainer, func(), error) {
	bigQueryContainer, err := RunBigQueryContainer(
		ctx,
		func() testcontainers.CustomizeRequestOption {
			return func(req *testcontainers.GenericContainerRequest) {
				// copy seed data into the container
				req.Files = append(req.Files, testcontainers.ContainerFile{
					// this assumes this program is run from the root of the repository
					// HostFilePath:      path.Join("/Users/perfect/Projects/iterfuncs", "testdata", "data.yaml"),
					HostFilePath:      path.Join(".", "testdata", "data.yaml"),
					ContainerFilePath: "/data.yaml",
				})

				req.Cmd = []string{"--project", projectId, "--data-from-yaml", "/data.yaml"}
			}
		}(),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to run container")
	}

	return bigQueryContainer, func() {
		if err := bigQueryContainer.Terminate(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to terminate container")
		}
	}, nil
}

func RunBigQueryContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*gcloud.GCloudContainer, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "ghcr.io/goccy/bigquery-emulator:0.4.3",
			ExposedPorts: []string{"9050/tcp", "9060/tcp"},
			WaitingFor:   wait.ForHTTP("/discovery/v1/apis/bigquery/v2/rest").WithPort("9050/tcp").WithStartupTimeout(time.Second * 5),
		},
		Started: true,
	}

	for _, opt := range opts {
		opt.Customize(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, err
	}

	spannerContainer, err := newGCloudContainer(ctx, 9050, container)
	if err != nil {
		return nil, err
	}

	// always prepend http:// to the URI
	spannerContainer.URI = "http://" + spannerContainer.URI

	return spannerContainer, nil
}

func newGCloudContainer(ctx context.Context, port int, c testcontainers.Container) (*gcloud.GCloudContainer, error) {
	mappedPort, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", port)))
	if err != nil {
		return nil, err
	}

	hostIP, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())

	gCloudContainer := &gcloud.GCloudContainer{
		Container: c,
		URI:       uri,
	}

	return gCloudContainer, nil
}
