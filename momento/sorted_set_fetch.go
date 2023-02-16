package momento

import (
	"context"
	"fmt"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

//////// Response

type SortedSetElement struct {
	Value []byte
	Score float64
}
type SortedSetFetchResponse interface {
	isSortedSetFetchResponse()
}

// SortedSetFetchMiss Miss Response to a cache SortedSetFetch api request.
type SortedSetFetchMiss struct{}

func (SortedSetFetchMiss) isSortedSetFetchResponse() {}

// SortedSetFetchHit Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchHit struct {
	Elements []*SortedSetElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

///// Request

type SortedSetOrder int

const (
	ASCENDING  SortedSetOrder = 0
	DESCENDING SortedSetOrder = 1
)

type SortedSetFetchNumResults interface {
	isSortedSetFetchNumResults()
}

type FetchAllElements struct{}

func (FetchAllElements) isSortedSetFetchNumResults() {}

type FetchLimitedElements struct {
	Limit uint32
}

func (FetchLimitedElements) isSortedSetFetchNumResults() {}

type SortedSetFetchRequest struct {
	CacheName       string
	SetName         string
	Order           SortedSetOrder
	NumberOfResults SortedSetFetchNumResults

	grpcRequest  *pb.XSortedSetFetchRequest
	grpcResponse *pb.XSortedSetFetchResponse
	response     SortedSetFetchResponse
}

func (r *SortedSetFetchRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	grpcReq := &pb.XSortedSetFetchRequest{
		SetName: []byte(r.SetName),
		Order:   pb.XSortedSetFetchRequest_Order(r.Order),
	}

	switch numResults := r.NumberOfResults.(type) {
	case *FetchAllElements:
		grpcReq.NumResults = &pb.XSortedSetFetchRequest_All{}
	case FetchAllElements:
		grpcReq.NumResults = &pb.XSortedSetFetchRequest_All{}
	case nil:
		grpcReq.NumResults = &pb.XSortedSetFetchRequest_All{}
	case *FetchLimitedElements:
		grpcReq.NumResults = &pb.XSortedSetFetchRequest_Limit{
			Limit: &pb.XSortedSetFetchRequest_XLimit{
				Limit: numResults.Limit,
			},
		}
	case FetchLimitedElements:
		grpcReq.NumResults = &pb.XSortedSetFetchRequest_Limit{
			Limit: &pb.XSortedSetFetchRequest_XLimit{
				Limit: numResults.Limit,
			},
		}
	default:
		return fmt.Errorf("unexpected fetch results type %T", r.NumberOfResults)
	}

	r.grpcRequest = grpcReq
	return nil
}

func (r *SortedSetFetchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetFetchRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		r.response = &SortedSetFetchHit{
			Elements: sortedSetGrpcElementToModel(grpcResp.Found.GetElements()),
		}
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func sortedSetGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []*SortedSetElement {
	var returnList []*SortedSetElement
	for _, element := range grpcSetElements {
		returnList = append(returnList, &SortedSetElement{
			Value: element.Name,
			Score: element.Score,
		})
	}
	return returnList
}
