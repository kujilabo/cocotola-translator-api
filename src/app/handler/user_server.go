package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	pb "github.com/kujilabo/cocotola-translator-api/src/proto"
)

type userServer struct {
	pb.UnimplementedTranslatorUserServer
	userUsecase usecase.UserUsecase
}

func NewTranslatorUserServer(userUsecase usecase.UserUsecase) pb.TranslatorUserServer {
	return &userServer{
		userUsecase: userUsecase,
	}
}

func (s *userServer) DictionaryLookup(ctx context.Context, in *pb.DictionaryLookupParameter) (*pb.DictionaryLookupResponses, error) {

	fromLang, err := domain.NewLang2(in.FromLang2)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "bad request").Err()
	}

	toLang, err := domain.NewLang2(in.ToLang2)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "bad request").Err()
	}

	results, err := s.userUsecase.DictionaryLookup(ctx, fromLang, toLang, in.Text)
	if err != nil {
		return nil, err
	}

	dictionaryResponses := make([]*pb.DictionaryResponse, len(results))
	for i, r := range results {
		dictionaryResponses[i] = &pb.DictionaryResponse{
			Lang2:      r.GetLang2().String(),
			Text:       r.GetText(),
			Pos:        int32(r.GetPos()),
			Translated: r.GetTranslated(),
		}
	}

	return &pb.DictionaryLookupResponses{
		Results: dictionaryResponses,
	}, nil
}

func (s *userServer) DictionaryLookupWithPos(ctx context.Context, in *pb.DictionaryLookupWithPosParameter) (*pb.DictionaryLookupResponse, error) {

	fromLang, err := domain.NewLang2(in.FromLang2)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "bad request").Err()
	}

	toLang, err := domain.NewLang2(in.ToLang2)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "bad request").Err()
	}

	result, err := s.userUsecase.DictionaryLookupWithPos(ctx, fromLang, toLang, in.Text, domain.WordPos(in.Pos))
	if err != nil {
		return nil, err
	}

	return &pb.DictionaryLookupResponse{
		Result: &pb.DictionaryResponse{
			Lang2:      result.GetLang2().String(),
			Text:       result.GetText(),
			Pos:        int32(result.GetPos()),
			Translated: result.GetTranslated(),
		},
	}, nil
}
